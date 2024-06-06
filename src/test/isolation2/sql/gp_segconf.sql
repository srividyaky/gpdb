--
-- Test that having an incomplete set of mirrors do not cause dtx recovery dispatch
-- to panic/segfault besides just error/retry.
-- We could have this kind of state during gpinitsystem or gpexpand.
--

!\retcode gpconfig -c gp_dtx_recovery_interval -v 5;
!\retcode gpstop -u;

select pg_catalog.gp_add_segment_primary('localhost', 'localhost', 8001, '/non/exist/path1');
select pg_catalog.gp_add_segment_primary('localhost', 'localhost', 8002, '/non/exist/path2');
select pg_catalog.gp_add_segment_primary('localhost', 'localhost', 8003, '/non/exist/path3');
select pg_catalog.gp_add_segment_mirror(3::int2, 'localhost', 'localhost', 9001, '/non/exist/path4');

-- check that the dtx recovery process has done at least one dispatch w/o panic/segfault.
select gp_inject_fault('dtx_recovery_dispatch_caught_error', 'skip', dbid) from gp_segment_configuration where role = 'p' and content = -1;
select gp_wait_until_triggered_fault('dtx_recovery_dispatch_caught_error', 1, dbid) from gp_segment_configuration where role = 'p' and content = -1;
select gp_inject_fault('dtx_recovery_dispatch_caught_error', 'reset', dbid) from gp_segment_configuration where role = 'p' and content = -1;

-- restore
-1U: select pg_catalog.gp_remove_segment(dbid) from gp_segment_configuration where content = 3;
-1U: select pg_catalog.gp_remove_segment(dbid) from gp_segment_configuration where content = 4;
-1U: select pg_catalog.gp_remove_segment(dbid) from gp_segment_configuration where content = 5;

!\retcode gpconfig -r gp_dtx_recovery_interval;
!\retcode gpstop -u;
