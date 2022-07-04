// Copyright 2022. Motty Cohen
//
// In-memory datastore implementation of IDatastore (used for testing)

package database

import (
	"fmt"
	"strings"
	"time"

	. "github.com/mottyc/yaaf-common/entity"
	"github.com/mottyc/yaaf-common/logger"
	"github.com/mottyc/yaaf-common/utils"
)

// region Database store definitions -----------------------------------------------------------------------------------

// Represent a db with tables
type InMemoryDatastore struct {
	db map[string]ITable
}

// Resolve index name from entity name
func indexName(ef EntityFactory, key string) (name string) {

	index := ef().TABLE()

	// Replace templates: {{accountId}}
	index = strings.Replace(index, "{{accountId}}", key, -1)

	// Replace templates: {{year}}
	index = strings.Replace(index, "{{year}}", time.Now().Format("2006"), -1)

	// Replace templates: {{month}}
	index = strings.Replace(index, "{{month}}", time.Now().Format("01"), -1)

	return index
}

// endregion

// region Factory and connectivity methods for Datastore ---------------------------------------------------------------

/**
 * Factory method for Datastore
 */
func NewInMemoryDatastore() (dbs IDatastore, err error) {
	return &InMemoryDatastore{db: make(map[string]ITable)}, nil
}

/**
 * Test connectivity for retries number of time with time interval (in seconds) between retries
 * @param retries - how many retries are required (max 10)
 * @param interval - time interval (in seconds) between retries (max 60)
 */
func (dbs *InMemoryDatastore) Ping(retries uint, interval uint) error {
	return nil
}

/**
 * Close Datastore and free resources
 */
func (dbs *InMemoryDatastore) Close() {
	logger.Debug("In memory datastore closed")
}

//endregion

// region Datastore Basic CRUD methods ----------------------------------------------------------------------------

/**
 * Get single entity by ID
 */
func (dbs *InMemoryDatastore) Get(factory EntityFactory, entityID string) (result Entity, err error) {

	entity := factory()
	table := tableName(entity.TABLE(), entity.KEY())
	if tbl, ok := dbs.db[table]; ok {
		return tbl.Get(entityID)
	} else {
		return nil, fmt.Errorf(TABLE_NOT_EXISTS)
	}
}

/**
 * Get list of entities by IDs
 */
func (dbs *InMemoryDatastore) List(factory EntityFactory, entityIDs []string) (list []Entity, err error) {

	entity := factory()
	table := tableName(entity.TABLE(), entity.KEY())

	list = make([]Entity, 0)

	if tbl, ok := dbs.db[table]; ok {
		for _, id := range entityIDs {
			if ent, err := tbl.Get(id); err == nil {
				list = append(list, ent)
			}
		}
	} else {
		return list, fmt.Errorf(TABLE_NOT_EXISTS)
	}
	return
}

/**
 * Check if entity exists by ID
 */
func (dbs *InMemoryDatastore) Exists(factory EntityFactory, entityID string) (result bool, err error) {

	entity := factory()
	table := tableName(entity.TABLE(), entity.KEY())

	if tbl, ok := dbs.db[table]; ok {
		return tbl.Exists(entityID)
	} else {
		return false, fmt.Errorf(TABLE_NOT_EXISTS)
	}
}

/**
 * Add new entity
 */
func (dbs *InMemoryDatastore) Insert(entity Entity) (added Entity, err error) {

	table := tableName(entity.TABLE(), entity.KEY())

	if _, ok := dbs.db[table]; !ok {
		dbs.db[table] = NewInMemTable()
	}

	return dbs.db[table].Insert(entity)
}

/**
 * Update existing entity in the data store
 */
func (dbs *InMemoryDatastore) Update(entity Entity) (updated Entity, err error) {
	table := tableName(entity.TABLE(), entity.KEY())

	if _, ok := dbs.db[table]; !ok {
		dbs.db[table] = NewInMemTable()
	}

	return dbs.db[table].Update(entity)
}

/**
 * Update existing entity in the data store or add it if it does not exist
 */
func (dbs *InMemoryDatastore) Upsert(entity Entity) (updated Entity, err error) {
	table := tableName(entity.TABLE(), entity.KEY())
	if tbl, ok := dbs.db[table]; ok {
		return tbl.Upsert(entity)
	} else {
		return nil, fmt.Errorf(TABLE_NOT_EXISTS)
	}
}

/**
 * Delete entity by id
 */
func (dbs *InMemoryDatastore) Delete(factory EntityFactory, entityID string) (err error) {

	entity := factory()

	table := tableName(entity.TABLE(), entity.KEY())
	if tbl, ok := dbs.db[table]; ok {
		return tbl.Delete(entityID)
	} else {
		return fmt.Errorf(TABLE_NOT_EXISTS)
	}
}

/**
 * Add multiple entities to data store (all must be of the same type)
 */
func (dbs *InMemoryDatastore) BulkInsert(entities []Entity) (affected int64, err error) {
	if len(entities) == 0 {
		return 0, nil
	}
	for _, ent := range entities {
		if _, err := dbs.Insert(ent); err == nil {
			affected += 1
		}
	}
	return affected, nil
}

/**
 * Update multiple entities in the data store (all must be of the same type)
 */
func (dbs *InMemoryDatastore) BulkUpdate(entities []Entity) (affected int64, err error) {
	if len(entities) == 0 {
		return 0, nil
	}
	for _, ent := range entities {
		if _, err := dbs.Update(ent); err == nil {
			affected += 1
		}
	}
	return affected, nil
}

/**
 * Update or insert multiple entities in the data store (all must be of the same type)
 */
func (dbs *InMemoryDatastore) BulkUpsert(entities []Entity) (affected int64, err error) {
	if len(entities) == 0 {
		return 0, nil
	}
	for _, ent := range entities {
		if _, err := dbs.Upsert(ent); err == nil {
			affected += 1
		}
	}
	return affected, nil
}

/**
 * Delete multiple entities by id list
 */
func (dbs *InMemoryDatastore) BulkDelete(factory EntityFactory, entityIDs []string) (affected int64, err error) {
	if len(entityIDs) == 0 {
		return 0, nil
	}
	for _, entityID := range entityIDs {
		if err := dbs.Delete(factory, entityID); err == nil {
			affected += 1
		}
	}
	return affected, nil
}

/**
 * Update single field of the document in a single transaction (eliminates the need to fetch - change - update)
 */
func (dbs *InMemoryDatastore) SetField(factory EntityFactory, entityID string, field string, value any) (err error) {
	fields := make(map[string]any)
	fields[field] = value
	return dbs.SetFields(factory, entityID, fields)
}

/**
 * Update some numeric fields of the document in a single transaction (eliminates the need to fetch - change - update)
 */
func (dbs *InMemoryDatastore) SetFields(factory EntityFactory, entityID string, fields map[string]any) (err error) {

	entity, fe := dbs.Get(factory, entityID)
	if fe != nil {
		return fe
	}

	// convert entity to Json
	js, fe := utils.JsonUtils().ToJson(entity)
	if fe != nil {
		return fe
	}

	// set fields
	for k, v := range fields {
		js[k] = v
	}

	toSet, fe := utils.JsonUtils().FromJson(factory, js)
	if fe != nil {
		return fe
	}

	_, fe = dbs.Update(toSet)
	return fe
}

/**
 * Helper method to construct query
 */
func (dbs *InMemoryDatastore) Query(factory EntityFactory) IQuery {
	return &inMemoryDatastoreQuery{
		db:         dbs,
		factory:    factory,
		allFilters: make([][]QueryFilter, 0),
		anyFilters: make([][]QueryFilter, 0),
		ascOrders:  make([]any, 0),
		descOrders: make([]any, 0),
		callbacks:  make([]func(in Entity) Entity, 0),
		limit:      100,
		page:       0,
	}
}

// endregion

// region Datastore Index methods --------------------------------------------------------------------------------------

/**
 * Test if index exists
 */
func (dbs *InMemoryDatastore) IndexExists(indexName string) (exists bool) {

	if _, ok := dbs.db[indexName]; ok {
		return true
	} else {
		return false
	}
}

/**
 * Create index by name
 */
func (dbs *InMemoryDatastore) CreateIndex(indexName string) (name string, err error) {
	// Create index
	if _, ok := dbs.db[indexName]; !ok {
		dbs.db[indexName] = &InMemoryTable{table: make(map[string]Entity)}
	}
	return indexName, nil
}

/**
 * Create index of entity and add entity field mapping
 */
func (dbs *InMemoryDatastore) CreateEntityIndex(ef EntityFactory, key string) (name string, err error) {

	idxName := indexName(ef, key)

	// Create index
	if _, ok := dbs.db[idxName]; !ok {
		dbs.db[idxName] = NewInMemTable()
	}
	return idxName, nil
}

/**
 * Drop index
 */
func (dbs *InMemoryDatastore) DropIndex(indexName string) (ack bool, err error) {

	if _, ok := dbs.db[indexName]; !ok {
		return false, nil
	} else {
		delete(dbs.db, indexName)
		return true, nil
	}
}

//endregion
