libovndb
========

A Go library for OVN DB access using native OVSDB protocol.
It is based on the [OVSDB Library](https://github.com/socketplane/libovsdb.git)

## What is OVN?

OVN (Open Virtual Network), is a SDN solution built on top of OVS (Open vSwitch).
The interface of OVN is its north-bound DB which is an OVSDB database.

## What is OVSDB?

OVSDB is a protocol for managing the configuration of OVS.
It's defined in [RFC 7047](http://tools.ietf.org/html/rfc7047)

## Why native OVSDB protocol?

There are projects accessing OVN based on the ovn-nbctl CLI. There are two major
issues with those approaches, which can be addressed by this library.

- Performance problem. Every command would trigger a separate OVSDB connection setup/teardown,
  initial OVSDB client cache population, etc., which would impact performance significantly. This library uses OVSDB
  protocol directly so that those overhead happens only once for all OVSDB operations.

- Caching problem. When there is a change in desired state, which requires updates in OVN, we need
  to figure out first what's the current state in OVN, which requires either maintaining a client cache or executing a "list" command everytime.
  This library maintains an internal cache and ensures it is always up to date with the remote DB with the help of native OVSDB support.

- String parsing problem. CLI based implementation needs extra convertion from the string output to Go internal data types, while it is not necessary
  with this library since OVSDB JSON RPC takes care of it.

## TODO

- Support transaction for multiple operations.
