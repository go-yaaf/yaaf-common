package database

import (
	. "github.com/go-yaaf/yaaf-common/entity"
	"io"
)

// IDatabase Database interface
type IDatabase interface {

	// Closer includes method Close()
	io.Closer

	// Ping Test database connectivity for retries number of time with time interval (in seconds) between retries
	Ping(retries uint, intervalInSeconds uint) error

	// Get a single entity by ID
	Get(factory EntityFactory, entityID string, keys ...string) (result Entity, err error)

	// List Get multiple entities by IDs
	List(factory EntityFactory, entityIDs []string, keys ...string) (list []Entity, err error)

	// Exists Check if entity exists by ID
	Exists(factory EntityFactory, entityID string, keys ...string) (result bool, err error)

	// Insert new entity
	Insert(entity Entity) (added Entity, err error)

	// Update existing entity
	Update(entity Entity) (updated Entity, err error)

	// Upsert Update entity or create it if it does not exist
	Upsert(entity Entity) (updated Entity, err error)

	// Delete entity by id and shard (key)
	Delete(factory EntityFactory, entityID string, keys ...string) (err error)

	// BulkInsert Insert multiple entities
	BulkInsert(entities []Entity) (affected int64, err error)

	// BulkUpdate Update multiple entities
	BulkUpdate(entities []Entity) (affected int64, err error)

	// BulkUpsert Update or insert multiple entities
	BulkUpsert(entities []Entity) (affected int64, err error)

	// BulkDelete Delete multiple entities by IDs
	BulkDelete(factory EntityFactory, entityIDs []string, keys ...string) (affected int64, err error)

	// SetField Update single field of the document in a single transaction (eliminates the need to fetch - change - update)
	SetField(factory EntityFactory, entityID string, field string, value any, keys ...string) (err error)

	// SetFields Update some fields of the document in a single transaction (eliminates the need to fetch - change - update)
	SetFields(factory EntityFactory, entityID string, fields map[string]any, keys ...string) (err error)

	// BulkSetFields Update specific field of multiple entities in a single transaction (eliminates the need to fetch - change - update)
	// The field is the name of the field, values is a map of entityId -> field value
	BulkSetFields(factory EntityFactory, field string, values map[string]any, keys ...string) (affected int64, err error)

	// Query Utility struct method to build a query
	Query(factory EntityFactory) IQuery

	// DDL Actions -----------------------------------------------------------------------------------------------------

	// ExecuteDDL Execute DDL - create table and indexes
	ExecuteDDL(ddl map[string][]string) (err error)

	// ExecuteSQL Execute SQL - execute SQL command
	ExecuteSQL(sql string, args ...any) (affected int64, err error)

	// DropTable Drop table and indexes
	DropTable(table string) (err error)

	// PurgeTable Fast delete table content (truncate)
	PurgeTable(table string) (err error)
}
