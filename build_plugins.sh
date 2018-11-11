#!/bin/bash

cur=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
BUILD_DIR=build

function build_plugin() {
    plugin_type=$1
    for dir in $(find plugins/$plugin_type -mindepth 1 -type d); do
        name=${dir##*"/"}
        src=$(ls $dir/*.go|grep -v .*_test.go)
        cmd="GO111MODULE=on go build -buildmode=plugin -o $BUILD_DIR/${name}_${plugin_type}.so $src"
        echo $cmd
        eval $cmd
    done
}

build_plugin parser $*
build_plugin output $*
