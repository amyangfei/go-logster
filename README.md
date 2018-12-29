# Logster implemented by golang

[![Go Report Card](https://goreportcard.com/badge/github.com/amyangfei/go-logster)](https://goreportcard.com/report/github.com/amyangfei/go-logster)
[![Build Status](https://travis-ci.org/amyangfei/go-logster.svg?branch=master)](https://travis-ci.org/amyangfei/go-logster)
[![Coverage Status](https://coveralls.io/repos/github/amyangfei/go-logster/badge.svg?branch=master)](https://coveralls.io/github/amyangfei/go-logster?branch=master)

This is a golang implemention of [logster](https://github.com/etsy/logster)

## Dependency

* go-logster uses the `logtail2` utility for gathering data from a logfile. This tool can be installed via package manager, but in fact use the package `logcheck`.

## Usage

As golang is a static programming language, we can't use dynamic parser/output class like the python version `logster`. However we can still create parser/output class easily with the help of [golang plugin system](https://golang.org/pkg/plugin/)

Two parser plugin samples and two output plugin samples are provided in this project. A parser plugin must implement the `Parser` interface defined in `logster/helper.go` and export an XXParser variable named `Parser`. An output plugin must implement the `Output` interface defined in `logster/helper.go` and export an XXOutput variable named `Output`.

When you finish your parser/output plugin and put the code to right dir, just use `make` to build `logster` binary and plugin shared object library. The version of Go must be 1.11 or above.

You can test go-logster from the command line. The --dry-run option will allow you to see the metrics being generated on stdout rather than sending them to your configured output. Besides you should provide the path of so file of parser/output plugin. For example:

```bash
$ ./build/logster -p server.101 -s run -l log \
    -o build/graphite_output.so -O '{"host": "127.0.0.1:8125", "protocol":"udp"}' \
    build/sample_parser.so test.log
```

Additional usage details can be found with the -h option:

```bash
$ ./logster -h
Usage:
  logster [OPTIONS] ParserPlugin LogFile

Application Options:
  -p, --metric-prefix=  Add prefix to all published metrics. This is for people that may multiple instances of same service on same host.
  -x, --metric-suffix=  Add suffix to all published metrics. This is for people that may add suffix at the end of their metrics.
  -l, --log-dir=        Where to store the logster logfile. (default: /var/log/logster)
  -s, --state-dir=      Where to store the tailer state file. (default: /var/run)
  -o, --output=         Where to send metrics (can specify multiple times)
  -v, --version         print version
  -d, --dry-run         Parse the log file but send stats to standard output.
  -D, --debug           Provide more verbose logging for debugging.
  -P, --parser-options= specific parser options
  -O, --output-options= specific output options

Help Options:
  -h, --help            Show this help message
```
