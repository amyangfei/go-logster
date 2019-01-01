#!/bin/bash -e
#
# Run all tests


PACKAGES=$(go list ./... | grep -vE 'vendor')
FILES=$(find . -name "*.go" | grep -vE "vendor")
CURDIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
GOFAIL_DIR=$(for p in $(go list ./...); do echo ${p#"github.com/amyangfei/go-logster/"}; done)
GOFAIL_ENABLE=$(echo $GOFAIL_DIR | xargs gofail enable)
GOFAIL_DISABLE=$(echo $GOFAIL_DIR | xargs gofail disable)

echo "Running tests..."
GO111MODULE=on go get github.com/etcd-io/gofail
echo $GOFAIL_DIR | xargs gofail enable
# fix main package in all plugins dir
sed -i "s/^package .*/package main/g" $(find plugins -name \*.fail.go)
if [ "$GL_TRAVIS_CI" = "on" ] ; then
    cover_opts="-covermode=count -coverprofile=coverage.out"
else
    cover_opts="-cover"
    export UT_PARSER_PLUGIN_PATH=${CURDIR}/build/sample_parser.so
    export UT_OUTPUT_PLUGIN_PATH=${CURDIR}/build/stdout_output.so
fi
GO111MODULE=on go test -v ${cover_opts} ${PACKAGES}
echo $GOFAIL_DIR | xargs gofail disable

echo "Checking gofmt..."
gofmt -s -l -w ${FILES} 2>&1 | awk '{print} END{if(NR>0) {exit 1}}'

echo "Checking govet..."
GO111MODULE=on go vet -all ${PACKAGES} 2>&1 | awk '{print} END{if(NR>0) {exit 1}}'

# GO111MODULE=off go get github.com/kisielk/errcheck
# echo "errcheck"
# errcheck -blank ${PACKAGES} | grep -v "_test\.go" | awk '{print} END{if(NR>0) {exit 1}}'

GO111MODULE=off go get golang.org/x/lint/golint
echo "Checking golint..."
golint -set_exit_status ${PACKAGES}

if [ "$GL_TRAVIS_CI" = "on" ] ; then
    GO111MODULE=on $GOPATH/bin/goveralls -coverprofile=coverage.out -service=travis-ci
fi

echo "Success"
