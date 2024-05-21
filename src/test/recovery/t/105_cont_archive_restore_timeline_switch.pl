# Test for validating if continuous archive restore pause and resume
# can handle timeline switches that are archived while the standby is
# up and running.

use strict;
use warnings;
use PostgreSQL::Test::Cluster;
use PostgreSQL::Test::Utils;
use Test::More;

my $node_primary;
# Mirror standby
my $node_standby;
# Disaster recovery standby
my $node_standby_dr;

sub test_pause_in_recovery
{
	my ($restore_point, $test_lsn, $num_rows) = @_;

	# Wait until standby has replayed enough data
	my $caughtup_query = "SELECT pg_last_wal_replay_lsn() = '$test_lsn'::pg_lsn";
	$node_standby_dr->poll_query_until('postgres', $caughtup_query)
		or die "Timed out while waiting for standby to catch up";

	# Check data has been replayed
	my $result = $node_standby_dr->safe_psql('postgres', "SELECT count(*) FROM table_foo;");
	is($result, $num_rows, "check standby content for $restore_point");
	ok($node_standby_dr->safe_psql('postgres', 'SELECT pg_is_wal_replay_paused();') eq 't',
		"standby is paused in recovery on $restore_point");
}

# Initialize and start primary node
$node_primary = PostgreSQL::Test::Cluster->new('primary');
$node_primary->init(has_archiving => 1, allows_streaming => 1);
$node_primary->start;

# Create data before taking the backup
$node_primary->safe_psql('postgres', "CREATE TABLE table_foo AS SELECT generate_series(1,1000);");
# Take a backup from which all operations will be run
$node_primary->backup('my_backup');

# Create standby node based off my_backup, with streaming enabled
$node_standby = PostgreSQL::Test::Cluster->new("standby");
$node_standby->init_from_backup($node_primary, 'my_backup', has_streaming => 1);
$node_standby->start;

my $lsn0 = $node_primary->safe_psql('postgres', "SELECT pg_create_restore_point('rp0');");
# Switching WAL guarantees that the restore point is available to the standby
$node_primary->safe_psql('postgres', "SELECT pg_switch_wal();");

# Find the next WAL segment to be archived
my $walfile_to_be_archived = $node_primary->safe_psql('postgres',
	"SELECT pg_walfile_name(pg_current_wal_lsn());");

# Wait for the WAL to be written to disk on Standby and archived
my $archive_wait_query =
  "SELECT '$walfile_to_be_archived' <= last_archived_wal FROM pg_stat_archiver;";
$node_primary->poll_query_until('postgres', $archive_wait_query)
  or die "Timed out while waiting for WAL segment to be archived";
$node_primary->wait_for_catchup($node_standby);

# Create a new standby with restore enabled, based off the original backup
# Note: The recovery_target_timeline will be set to "latest" by default
$node_standby_dr = PostgreSQL::Test::Cluster->new("standby_dr");
$node_standby_dr->init_from_backup($node_primary, 'my_backup', has_restoring => 1);

# Enable `hot_standby`
$node_standby_dr->append_conf('postgresql.conf', qq(hot_standby = 'on'));

# Set rp0 as a restore point to pause on start up
$node_standby_dr->append_conf('postgresql.conf', qq(gp_pause_on_restore_point_replay = 'rp0'));
# Start the standby
$node_standby_dr->start;
test_pause_in_recovery('rp0', $lsn0, 1000);

# Make sure the promoted standby has archiving set to the same archive location before promotion
# so that the timeline history file gets archived.
ok($node_primary->safe_psql('postgres', 'SHOW archive_command;') eq $node_standby_dr->safe_psql('postgres', 'SHOW archive_command;'),
"check if archive_command is the same for both primary and standby");

# Promote the standby node while the primary is still up
# thus causing a timeline conflict
$node_standby->promote;

# Create some INSERT records and switch WAL to make this WAL segment
# file available to the standby
$node_primary->safe_psql('postgres', "INSERT INTO table_foo SELECT generate_series(1,10);");
$node_primary->safe_psql('postgres', "SELECT pg_switch_wal();");

# Create another simple table with small data on the promoted standby node
# Swtich WAL
$node_standby->safe_psql('postgres', "INSERT INTO table_foo VALUES (generate_series(1001,2000))");
$node_standby->safe_psql('postgres', "SELECT pg_switch_wal();");

# Create a second restore point. Switch WAL
my $lsn_standby = $node_standby->safe_psql('postgres', "SELECT pg_create_restore_point('rp_standby');");
$node_standby->safe_psql('postgres', "SELECT pg_switch_wal();");

# Set `gp_pause_on_restore_point_replay` GUC to second restore point
# send SIGHUP (pg_ctl reload) to second standby node and resume WAL
# Advance to rp_standby
$node_standby_dr->adjust_conf('postgresql.conf', 'gp_pause_on_restore_point_replay', "rp_standby");
$node_standby_dr->reload;
$node_standby_dr->safe_psql('postgres', "SELECT pg_wal_replay_resume();");

# Validate that the secondary Standby node reached the second restore point rp_standby and has not
# become blocked waiting for new WAL on the old timeline
test_pause_in_recovery('rp_standby', $lsn_standby, 2000);

$node_primary->teardown_node;
$node_standby->teardown_node;
$node_standby_dr->teardown_node;

done_testing();
