#!/bin/sh -ex

SCRIPT=`realpath -s $0`
SCRIPTPATH=`dirname $SCRIPT`

. $SCRIPTPATH/test_common.sh

(
echo "building ovn in ${ovn_srcdir}"
cd ${ovn_srcdir}
rm -rf ovs
rm -rf ovn

git clone --depth 1 -b master https://github.com/openvswitch/ovs.git
git clone --depth 1 -b master https://github.com/ovn-org/ovn.git

cd ovs
./boot.sh && ./configure --enable-silent-rules
make -j4
)
