libovndb
========

An OVNDB Library written in Go

## What is OVNDB?

OVSDB is Protocol for managing the configuration of OVN.
It's defined in [RFC 7047](http://tools.ietf.org/html/rfc7047)
OVNDB library developed based on [OVSDB Library](https://github.com/socketplane/libovsdb.git)

## Running example

On the node which has the OVN environment:

    go get github.com/golang/glog/
    go get github.com/socketplane/libovsdb
    go build -o ovndb-test ./examples/main.go
    ./ovndb-test

## e2e test case

    go get github.com/stretchr/testify/assert
    go test -c

##Change Log

v1.0:
Support lsp/lsw/addresset/acl creating/deleting
Cache supported will sync with ovndb automatically.


## Todo

Build dependency from makefile
Build example in makefile
Should have docker image with OVN to run the example
Run example supports in makefile
