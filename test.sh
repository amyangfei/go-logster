#!/bin/bash -e
#
# Run all tests

cur=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

FMT="logster apps/logster plugins/output/stdout plugins/parser/sample"


echo "Running tests..."
go test -v $(go list ./... | grep -v /vendor/)


echo "Checking gofmt..."
cd $cur
fmtRes=$(gofmt -l $FMT)
if [ -n "${fmtRes}" ]; then
    echo -e "gofmt checking failed:\n${fmtRes}"
    exit 255
fi


echo "Build test..."
make clean && make


echo "Success"
