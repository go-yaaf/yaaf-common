package database

import (
	"io"

	. "github.com/go-yaaf/yaaf-common/entity"
)

// IDatabase defines the interface for a Relational Database Management System (RDBMS) wrapper.
// It provides a unified API for CRUD operations, querying, schema management, and transaction handling.
type IDatabase interface {

	// Closer includes method Close() to close the database connection.
	io.Closer

	// Ping tests the database connectivity.
	// It retries the connection 'retries' times with an 'intervalInSeconds' delay between attempts.
	Ping(retries uint, intervalInSeconds uint) error

	// CloneDatabase returns a clone (copy) of the database instance.
	// This is useful for creating a new instance with the same configuration but potentially different state or connection pool.
	CloneDatabase() (IDatabase, error)

	// Get retrieves a single entity by its ID.
	// The 'keys' argument is optional and can be used for sharding or other specific lookup mechanisms.
	Get(factory EntityFactory, entityID string, keys ...string) (result Entity, err error)

	// List retrieves multiple entities by their IDs.
	// The 'keys' argument is optional and can be used for sharding or other specific lookup mechanisms.
	List(factory EntityFactory, entityIDs []string, keys ...string) (list []Entity, err error)

	// Exists checks if an entity exists by its ID.
	// The 'keys' argument is optional and can be used for sharding or other specific lookup mechanisms.
	Exists(factory EntityFactory, entityID string, keys ...string) (result bool, err error)

	// Insert adds a new entity to the database.
	// It returns the added entity (potentially with updated fields like ID or timestamps) or an error.
	Insert(entity Entity) (added Entity, err error)

	// Update modifies an existing entity in the database.
	// It returns the updated entity or an error if the entity does not exist or the update fails.
	Update(entity Entity) (updated Entity, err error)

	// Upsert updates an entity if it exists, or creates it if it does not.
	// It returns the updated or inserted entity.
	Upsert(entity Entity) (updated Entity, err error)

	// Delete removes an entity by its ID.
	// The 'keys' argument is optional and can be used for sharding or other specific lookup mechanisms.
	Delete(factory EntityFactory, entityID string, keys ...string) (err error)

	// BulkInsert adds multiple entities to the database in a batch operation.
	// It returns the number of affected records and any error encountered.
	BulkInsert(entities []Entity) (affected int64, err error)

	// BulkUpdate modifies multiple entities in the database in a batch operation.
	// It returns the number of affected records and any error encountered.
	BulkUpdate(entities []Entity) (affected int64, err error)

	// BulkUpsert updates or inserts multiple entities in the database in a batch operation.
	// It returns the number of affected records and any error encountered.
	BulkUpsert(entities []Entity) (affected int64, err error)

	// BulkDelete removes multiple entities by their IDs in a batch operation.
	// The 'keys' argument is optional and can be used for sharding or other specific lookup mechanisms.
	BulkDelete(factory EntityFactory, entityIDs []string, keys ...string) (affected int64, err error)

	// SetField updates a single field of an entity in a single transaction.
	// This eliminates the need to fetch, change, and update the entire entity.
	SetField(factory EntityFactory, entityID string, field string, value any, keys ...string) (err error)

	// SetFields updates multiple fields of an entity in a single transaction.
	// The 'fields' argument is a map where the key is the field name and the value is the new value.
	SetFields(factory EntityFactory, entityID string, fields map[string]any, keys ...string) (err error)

	// BulkSetFields updates a specific field for multiple entities in a single transaction.
	// The 'values' argument is a map where the key is the entity ID and the value is the new value for the field.
	BulkSetFields(factory EntityFactory, field string, values map[string]any, keys ...string) (affected int64, err error)

	// Query returns a new IQuery builder for the given entity factory.
	// The IQuery builder allows for constructing complex queries with filters, sorting, and pagination.
	Query(factory EntityFactory) IQuery

	// AdvancedQuery returns an IAdvancedQuery builder, which supports analytic queries.
	// It is intended for working with sharded tables where the KEY() method is used to resolve the table name.
	AdvancedQuery(factory EntityFactory) IAdvancedQuery

	// DDL Actions -----------------------------------------------------------------------------------------------------

	// ExecuteDDL executes Data Definition Language (DDL) commands to create tables and indexes.
	// The 'ddl' argument is a map where the key is the table name and the value is a list of fields to index.
	ExecuteDDL(ddl map[string][]string) (err error)

	// ExecuteSQL executes a raw SQL command.
	// It returns the number of affected records and any error encountered.
	ExecuteSQL(sql string, args ...any) (affected int64, err error)

	// ExecuteQuery executes a raw SQL query and returns the results as a list of Json maps.
	ExecuteQuery(source, sql string, args ...any) ([]Json, error)

	// DropTable drops a table and its indexes.
	DropTable(table string) (err error)

	// PurgeTable removes all data from a table (truncate).
	PurgeTable(table string) (err error)
}
