TODO
======

- Transaction support. Currently the lib supports a single transaction per API call.
  For complete transaction support like Python/C OVSDB clients, we need to maintain the
  local changes and support commit, rollback and auto retry, etc.

- L3 APIs. Currently the lib supports L2 objects operations and ACLs. APIs for L3 objects
  such as logical routers and ports, gateways, are to be added.
