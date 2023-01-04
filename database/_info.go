// Package database
//
// The database package includes interfaces for some data storage middlewares including:
// * IDatabase  - interface for RDBMS wrapper implementations
// * IDataCache - interface for distributed cache wrapper implementations
// * IDatastore - interface for NoSQL Big Data (Document Store) wrapper implementations
//
// The package also includes in-memory implementations of all the above mainly for testing but can be used
// for some use cases when data persistent is not required
package database
