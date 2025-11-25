package database

import (
	"io"

	. "github.com/go-yaaf/yaaf-common/entity"
)

// IDatastore defines the interface for NoSQL Big Data (Document Store) wrapper implementations (e.g., Couchbase, ElasticSearch).
// It provides methods for CRUD operations, querying, and index management.
type IDatastore interface {

	// Closer includes method Close() to close the datastore connection.
	io.Closer

	// Ping tests the datastore connectivity.
	// It retries the connection 'retries' times with an 'intervalInSeconds' delay between attempts.
	Ping(retries uint, intervalInSeconds uint) error

	// CloneDatastore returns a clone (copy) of the datastore instance.
	// This is useful for creating a new instance with the same configuration but potentially different state or connection pool.
	CloneDatastore() (IDatastore, error)

	// Get retrieves a single entity by its ID.
	// The 'keys' argument is optional and can be used for sharding or other specific lookup mechanisms.
	Get(factory EntityFactory, entityID string, keys ...string) (result Entity, err error)

	// List retrieves multiple entities by their IDs.
	// The 'keys' argument is optional and can be used for sharding or other specific lookup mechanisms.
	List(factory EntityFactory, entityIDs []string, keys ...string) (list []Entity, err error)

	// Exists checks if an entity exists by its ID.
	// The 'keys' argument is optional and can be used for sharding or other specific lookup mechanisms.
	Exists(factory EntityFactory, entityID string, keys ...string) (result bool, err error)

	// Insert adds a new entity to the datastore.
	// It returns the added entity (potentially with updated fields like ID or timestamps) or an error.
	Insert(entity Entity) (added Entity, err error)

	// Update modifies an existing entity in the datastore.
	// It returns the updated entity or an error if the entity does not exist or the update fails.
	Update(entity Entity) (updated Entity, err error)

	// Upsert updates an entity if it exists, or creates it if it does not.
	// It returns the updated or inserted entity.
	Upsert(entity Entity) (updated Entity, err error)

	// Delete removes an entity by its ID.
	// The 'keys' argument is optional and can be used for sharding or other specific lookup mechanisms.
	Delete(factory EntityFactory, entityID string, keys ...string) (err error)

	// BulkInsert adds multiple entities to the datastore in a batch operation.
	// It returns the number of affected records and any error encountered.
	BulkInsert(entities []Entity) (affected int64, err error)

	// BulkUpdate modifies multiple entities in the datastore in a batch operation.
	// It returns the number of affected records and any error encountered.
	BulkUpdate(entities []Entity) (affected int64, err error)

	// BulkUpsert updates or inserts multiple entities in the datastore in a batch operation.
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

	// Query returns a new IQuery builder for the given entity factory.
	// The IQuery builder allows for constructing complex queries with filters, sorting, and pagination.
	Query(factory EntityFactory) IQuery

	// IndexExists checks if an index exists.
	IndexExists(indexName string) (exists bool)

	// CreateIndex creates a new index (without mapping).
	CreateIndex(indexName string) (name string, err error)

	// CreateEntityIndex creates an index for an entity and adds entity field mapping.
	CreateEntityIndex(factory EntityFactory, key string) (name string, err error)

	// ListIndices returns a list of all indices matching the pattern.
	ListIndices(pattern string) (map[string]int, error)

	// DropIndex drops an index.
	DropIndex(indexName string) (ack bool, err error)

	// ExecuteQuery executes a native query (e.g., KQL, N1QL) and returns the results as a list of Json maps.
	ExecuteQuery(source string, query string, args ...any) ([]Json, error)
}
