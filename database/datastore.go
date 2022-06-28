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

	// Test connectivity for retries number of time with time interval (in seconds) between retries
	Ping(retries uint, intervalInSeconds uint) error

	// Get single entity by ID
	Get(accountID, entityID string, ef EntityFactory) (entity Entity, err error)

	// Check if entity exists by ID
	Exists(accountID, entityID string, ef EntityFactory) (result bool, err error)

	// Insert new entity
	Insert(accountID string, entity Entity) (added Entity, err error)

	// Update existing entity
	Update(accountID string, entity Entity) (updated Entity, err error)

	// Delete existing entity
	Delete(accountId, entityId string, ef EntityFactory) (err error)

	// Insert multiple entities
	BulkInsert(accountID string, entities []Entity) (affected int64, err error)

	// Delete multiple entities by IDs
	BulkDelete(entityIDs []string, key string, ef EntityFactory) (affected int, err error)

	// Update single string field of the document in a single transaction
	SetField(accountId string, entityId string, ef EntityFactory, field string, value any) (err error)

	// Update some fields of the document in a single transaction (eliminates the need to fetch - change - update)
	SetFields(entityID, key string, ef EntityFactory, fields map[string]any) (err error)

	// Utility struct to build a query
	Query(f EntityFactory) IQuery

	// Close connection and free resources
	Close()

	// Index Actions ---------------------------------------------------------------------------------------------------

	// Create index of entity and add entity field mapping
	CreateEntityIndex(ef EntityFactory, accountId string) (name string, err error)

	// Create index by name (without mapping)
	CreateIndex(indexName string) (name string, err error)

	// Drop index
	DropIndex(indexName string) (ack bool, err error)

	// Test if index exists
	IndexExists(indexName string) (exists bool)

	// Get indices (selective and sorted)
	GetIndices(indicesPattern, requestedColumns, sortByColumn string) (indices []string, err error)
}
