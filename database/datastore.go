// Copyright 2022. Motty Cohen
//
// Database interface for NoSQL Big Data wrapper implementations
//
package database

import (
	. "github.com/mottyc/yaaf-common/entity"
)

// Datastore interface
type IDatastore interface {

	// Test database connectivity for retries number of time with time interval (in seconds) between retries
	Ping(retries uint, intervalInSeconds uint) error

	// Get single entity by ID
	Get(factory EntityFactory, entityID string) (result Entity, err error)

	// Get multiple entities by IDs
	List(factory EntityFactory, entityIDs []string) (list []Entity, err error)

	// Check if entity exists by ID
	Exists(factory EntityFactory, entityID string) (result bool, err error)

	// Insert new entity
	Insert(entity Entity) (added Entity, err error)

	// Update existing entity
	Update(entity Entity) (updated Entity, err error)

	// Update entity or create it if it does not exist
	Upsert(entity Entity) (updated Entity, err error)

	// Delete entity by id and shard (key)
	Delete(factory EntityFactory, entityID string) (err error)

	// Insert multiple entities
	BulkInsert(entities []Entity) (affected int64, err error)

	// Update multiple entities
	BulkUpdate(entities []Entity) (affected int64, err error)

	// Update or insert multiple entities
	BulkUpsert(entities []Entity) (affected int64, err error)

	// Delete multiple entities by IDs
	BulkDelete(factory EntityFactory, entityIDs []string) (affected int64, err error)

	// Update single field of the document in a single transaction (eliminates the need to fetch - change - update)
	SetField(factory EntityFactory, entityID string, field string, value any) (err error)

	// Update some fields of the document in a single transaction (eliminates the need to fetch - change - update)
	SetFields(factory EntityFactory, entityID string, fields map[string]any) (err error)

	// Utility struct method to build a query
	Query(factory EntityFactory) IQuery

	// Close DB and free resources
	Close()

	// Index Actions ---------------------------------------------------------------------------------------------------
	// Test if index exists
	IndexExists(indexName string) (exists bool)

	// Create index by name (without mapping)
	CreateIndex(indexName string) (name string, err error)

	// Create index of entity and add entity field mapping
	CreateEntityIndex(ef EntityFactory, key string) (name string, err error)

	// Drop index
	DropIndex(indexName string) (ack bool, err error)
}
