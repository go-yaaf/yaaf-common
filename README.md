# yaaf-common

[![Build](https://github.com/mottyc/yaaf-common/actions/workflows/build.yml/badge.svg)](https://github.com/mottyc/yaaf-common/actions/workflows/build.yml)

Common interfaces to YAAF (Yet Another Application Framework) library


## About
This project is part of the complete Go application framework to provide wrapper and utilities for common middleware components required to create micro-services (servers and workers)

### BaseConfig
Base utility to provide service configuration via environment variable.

Configuration using environment variables is quite common in container-orchestration frameworks (e.g. Docker, Kubernetes etc) and this utility
provides a simple way to define and access application/service specific configuration parameters through accessors.

The basic implementation includes the fundamentals plus a pre-set of common configuration variables for web servers/worker.
For each application/service you should write your own configuration utility inherit from this structure and add specific configuration variables.

### IDatabase
Interface of object-relational database. It is used by concrete implementations of simple ORM over RDBMS ACID databases (e.g. PostgreSQL, MySQL...)
Under the database folder you should find in memory implementation of IDatabase interface, used mainly for testing.

### IDatastore
Interface of document database. It is used by concrete implementations of NoSQL Document oriented databases (e.g. Elasticsearch, Couchbase...).
Under the database folder you should find in memory implementation of IDatastore interface, used mainly for testing.

### Logger
Simple wrapper for zap logger used as system-wide logging framework

### Utils
Collection of utility helpers

#### Adding dependency

```bash
$ go get -v -t github.com/mottyc/yaaf-common ./...
```
