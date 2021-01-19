#!/bin/sh -ex

SCRIPT=`realpath -s $0`
SCRIPTPATH=`dirname $SCRIPT`

sh $SCRIPTPATH/test_prepare.sh || exit 1
sh $SCRIPTPATH/test_run.sh || exit 1
