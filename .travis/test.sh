#!/bin/sh -ex

rm -rf ovs
git clone https://github.com/openvswitch/ovs.git
cd ovs
./boot.sh && ./configure --enable-silent-rules
make -j4

srcdir=`pwd`
builddir=$srcdir
rm -rf sandbox
mkdir sandbox
sandbox=`cd sandbox && pwd`

# Below code is borrowed from OVS sandbox:
# https://github.com/openvswitch/ovs/blob/master/tutorial/ovs-sandbox

OVS_RUNDIR=$sandbox; export OVS_RUNDIR
OVS_LOGDIR=$sandbox; export OVS_LOGDIR
OVS_DBDIR=$sandbox; export OVS_DBDIR
OVS_SYSCONFDIR=$sandbox; export OVS_SYSCONFDIR
PATH=$builddir/ovsdb:$builddir/vswitchd:$builddir/utilities:$builddir/vtep:$PATH
PATH=$builddir/ovn/controller:$builddir/ovn/controller-vtep:$builddir/ovn/northd:$builddir/ovn/utilities:$PATH
export PATH

run() {
    echo "$@"
    (cd "$sandbox" && "$@") || exit 1
}
ovn_start_db() {
    local db=$1 model=$2 servers=$3 schema=$4
    local DB=$(echo $db | tr a-z A-Z)
    local schema_name=$(ovsdb-tool schema-name $schema)

    case $model in
        standalone | backup) ;;
        clustered)
            case $servers in
                [1-9] | [1-9][0-9]) ;;
                *) echo "${db}db servers must be between 1 and 99" >&2
                   exit 1
                   ;;
            esac
            ;;
        *)
            echo "unknown ${db}db model \"$model\"" >&2
            exit 1
            ;;
    esac

    ovn_start_ovsdb_server() {
        local i=$1; shift
        run ovsdb-server --detach --no-chdir \
               --pidfile=$db$i.pid -vconsole:off --log-file=$db$i.log \
               -vsyslog:off \
               --remote=db:$schema_name,${DB}_Global,connections \
               --private-key=db:$schema_name,SSL,private_key \
               --certificate=db:$schema_name,SSL,certificate \
               --ca-cert=db:$schema_name,SSL,ca_cert \
               --ssl-protocols=db:$schema_name,SSL,ssl_protocols \
               --ssl-ciphers=db:$schema_name,SSL,ssl_ciphers \
               --unixctl=${db}$i --remote=punix:$db$i.ovsdb ${db}$i.db "$@"
    }

    case $model in
        standalone)
            run ovsdb-tool create ${db}1.db "$schema"
            ovn_start_ovsdb_server 1
            remote=unix:${db}1.ovsdb
            ;;
        backup)
            for i in 1 2; do
                run ovsdb-tool create $db$i.db "$schema"
            done
            ovn_start_ovsdb_server 1
            ovn_start_ovsdb_server 2 --sync-from=unix:${db}1.ovsdb
            remote=unix:${db}1.ovsdb
            backup_note="$backup_note
The backup server of OVN $DB can be accessed by:
* ovn-${db}ctl --db=unix:`pwd`/sandbox/${db}2.ovsdb
* ovs-appctl -t `pwd`/sandbox/${db}2
The backup database file is sandbox/${db}2.db
"
            ;;
        clustered)
            for i in $(seq $servers); do
                if test $i = 1; then
                    run ovsdb-tool create-cluster ${db}1.db "$schema" unix:${db}1.raft;
                else
                    run ovsdb-tool join-cluster $db$i.db $schema_name unix:$db$i.raft unix:${db}1.raft
                fi
                ovn_start_ovsdb_server $i
            done
            remote=unix:${db}1.ovsdb
            for i in `seq 2 $servers`; do
                remote=$remote,unix:$db$i.ovsdb
            done
            for i in $(seq $servers); do
                run ovsdb-client wait unix:$db$i.ovsdb $schema_name connected
            done
            ;;
    esac
    eval OVN_${DB}_DB=\$remote
    eval export OVN_${DB}_DB
}

ovn_start_db nb standalone 1 $srcdir/ovn/ovn-nb.ovsschema

cd ../goovn; go test -v
pkill ovsdb-server
cd ../; rm -rf ovs
