#!/bin/bash

set -eux -o pipefail

ccp_src/scripts/setup_ssh_to_cluster.sh

scp cluster_env_files/hostfile_all cdw:/tmp
tar -xzf gp_binary/gp.tgz
scp gpctl/gpctl cdw:/home/gpadmin/
scp gpservice/gpservice cdw:/home/gpadmin/

ssh -n cdw "
    set -eux -o pipefail

    export PATH=/usr/local/go/bin:\$PATH
    source /usr/local/greenplum-db-devel/greenplum_path.sh

    chmod +x gpctl
    chmod +x gpservice
    gpsync -f /tmp/hostfile_all gpctl =:/usr/local/greenplum-db-devel/bin/gpctl
    gpsync -f /tmp/hostfile_all gpservice =:/usr/local/greenplum-db-devel/bin/gpservice
    cd /home/gpadmin/gpdb_src/gpMgmt/bin/go-tools
    ./ci/scripts/generate_ssl_cert_multi_host.bash

    make integration FILE=/tmp/hostfile_all UTILITY=${utility}
"
