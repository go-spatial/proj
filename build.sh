#!/bin/sh

set -e

for i in . cmd/proj core merror mlog operations support   # gie
do
    echo "*** $i ***"
    pushd $i &> /dev/null
    go test -v -cover
    if [ "$?" -ne "0" ]
    then
        echo fail
        exit 1
    fi
    popd &> /dev/null
done

