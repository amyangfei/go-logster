#!/bin/bash

# install dependcies for local developing

cur=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
src_dst=$GOPATH/src/github.com/amyangfei/go-logster

rm -rf $src_dst
mkdir -p $src_dst
cp -r $cur/../logster $src_dst

cd $src_dst/logster
echo "installing go-logster/logster"
go install

echo "done!"
