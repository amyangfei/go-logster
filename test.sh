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
for dir in apps/*/; do
    dir=${dir%/}
    if grep -q '^package main$' $dir/*.go 2>/dev/null; then
        echo "building $dir"
        go build -o $dir/$(basename $dir) ./$dir
    else
        echo "(skipped $dir)"
    fi
done


echo "Success"
