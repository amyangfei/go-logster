#!/bin/bash -e
#
# Run all tests


PACKAGES=$(go list ./... | grep -vE 'vendor|examples')
FILES=$(find . -name "*.go" | grep -vE "vendor|examples")

echo "Running tests..."
if [ "$GL_TRAVIS_CI" = "on" ] ; then
    cover_opts="-covermode=count -coverprofile=coverage.out"
else
    cover_opts="-cover"
fi
GO111MODULE=on go test -v ${cover_opts} ${PACKAGES}

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
