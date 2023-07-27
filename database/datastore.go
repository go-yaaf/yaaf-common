package database

import (
	. "github.com/go-yaaf/yaaf-common/entity"
	"io"
)

// IDatastore interface for NoSQL Big Data wrapper implementations
type IDatastore interface {

	// Closer includes method Close()
	io.Closer

	// Ping tests database connectivity for retries number of time with time interval (in seconds) between retries
	Ping(retries uint, intervalInSeconds uint) error

	// CloneDatastore Returns a clone (copy) of the instance
	CloneDatastore() (IDatastore, error)

	// Get a single entity by ID
	Get(factory EntityFactory, entityID string, keys ...string) (result Entity, err error)

	// List gets multiple entities by IDs
	List(factory EntityFactory, entityIDs []string, keys ...string) (list []Entity, err error)

	// Exists checks if entity exists by ID
	Exists(factory EntityFactory, entityID string, keys ...string) (result bool, err error)

	// Insert a new entity
	Insert(entity Entity) (added Entity, err error)

	// Update an existing entity
	Update(entity Entity) (updated Entity, err error)

	// Upsert update entity or create it if it does not exist
	Upsert(entity Entity) (updated Entity, err error)

	// Delete entity by id and shard (key)
	Delete(factory EntityFactory, entityID string, keys ...string) (err error)

	// BulkInsert inserts multiple entities
	BulkInsert(entities []Entity) (affected int64, err error)

	// BulkUpdate updates multiple entities
	BulkUpdate(entities []Entity) (affected int64, err error)

	// BulkUpsert update or insert multiple entities
	BulkUpsert(entities []Entity) (affected int64, err error)

	// BulkDelete delete multiple entities by IDs
	BulkDelete(factory EntityFactory, entityIDs []string, keys ...string) (affected int64, err error)

	// SetField update a single field of the document in a single transaction (eliminates the need to fetch - change - update)
	SetField(factory EntityFactory, entityID string, field string, value any, keys ...string) (err error)

	// SetFields update some fields of the document in a single transaction (eliminates the need to fetch - change - update)
	SetFields(factory EntityFactory, entityID string, fields map[string]any, keys ...string) (err error)

	// Query is a factory method for query builder Utility
	Query(factory EntityFactory) IQuery

	// IndexExists tests if index exists
	IndexExists(indexName string) (exists bool)

	// CreateIndex creates an index (without mapping)
	CreateIndex(indexName string) (name string, err error)

	// CreateEntityIndex creates an index of entity and add entity field mapping
	CreateEntityIndex(factory EntityFactory, key string) (name string, err error)

	// ListIndices returns a list of all indices matching the pattern
	ListIndices(pattern string) (map[string]int, error)

	// DropIndex drops an index
	DropIndex(indexName string) (ack bool, err error)
}
