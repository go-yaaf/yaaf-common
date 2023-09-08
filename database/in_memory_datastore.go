// Copyright 2022. Motty Cohen
//
// In-memory datastore implementation of IDatastore (used for testing)

package database

import (
	"fmt"
	"strings"
	"time"

	. "github.com/go-yaaf/yaaf-common/entity"
	"github.com/go-yaaf/yaaf-common/logger"
	"github.com/go-yaaf/yaaf-common/utils"
)

const (
	INDEX_NOT_EXISTS = "index does not exist"
)

// region Database store definitions -----------------------------------------------------------------------------------

// InMemoryDatastore Represent a db with tables
type InMemoryDatastore struct {
	db map[string]ITable
}

// Resolve index name from entity name
func indexName(table string, keys ...string) (name string) {

	index := table

	if len(keys) == 0 {
		return index
	}

	// Replace accountId placeholder with the first key
	index = strings.Replace(index, "{{accountId}}", "{{0}}", -1)

	for idx, key := range keys {
		placeHolder := fmt.Sprintf("{{%d}}", idx)
		index = strings.Replace(index, placeHolder, key, -1)
	}

	// Replace templates: {{year}}
	index = strings.Replace(index, "{{year}}", time.Now().Format("2006"), -1)

	// Replace templates: {{month}}
	index = strings.Replace(index, "{{month}}", time.Now().Format("01"), -1)

	// TODO: Replace templates: {{week}}

	return index
}

// endregion

// region Factory and connectivity methods for Datastore ---------------------------------------------------------------

// NewInMemoryDatastore Factory method for Datastore
func NewInMemoryDatastore() (dbs IDatastore, err error) {
	return &InMemoryDatastore{db: make(map[string]ITable)}, nil
}

// Ping tests database connectivity for retries number of time with time interval (in seconds) between retries
func (dbs *InMemoryDatastore) Ping(retries uint, interval uint) error {
	logger.Debug("Pinging %d times with %d interval", retries, interval)
	return nil
}

// Close Datastore and free resources
func (dbs *InMemoryDatastore) Close() error {
	logger.Debug("In memory datastore closed")
	return nil
}

// CloneDatastore Returns a clone (copy) of the instance
func (dbs *InMemoryDatastore) CloneDatastore() (IDatastore, error) {
	return dbs, nil
}

//endregion

// region Datastore Basic CRUD methods ----------------------------------------------------------------------------

// Get a single entity by ID
func (dbs *InMemoryDatastore) Get(factory EntityFactory, entityID string, keys ...string) (result Entity, err error) {

	entity := factory()
	index := indexName(entity.TABLE(), keys...)
	if tbl, ok := dbs.db[index]; ok {
		return tbl.Get(entityID)
	} else {
		return nil, fmt.Errorf(INDEX_NOT_EXISTS)
	}
}

// List gets multiple entities by IDs
func (dbs *InMemoryDatastore) List(factory EntityFactory, entityIDs []string, keys ...string) (list []Entity, err error) {

	index := indexName(factory().TABLE(), keys...)

	list = make([]Entity, 0)

	if tbl, ok := dbs.db[index]; ok {
		for _, id := range entityIDs {
			if ent, err := tbl.Get(id); err == nil {
				list = append(list, ent)
			}
		}
	} else {
		return list, fmt.Errorf(INDEX_NOT_EXISTS)
	}
	return
}

// Exists checks if entity exists by ID
func (dbs *InMemoryDatastore) Exists(factory EntityFactory, entityID string, keys ...string) (result bool, err error) {

	index := indexName(factory().TABLE(), keys...)

	if tbl, ok := dbs.db[index]; ok {
		return tbl.Exists(entityID)
	} else {
		return false, fmt.Errorf(INDEX_NOT_EXISTS)
	}
}

// Insert a new entity
func (dbs *InMemoryDatastore) Insert(entity Entity) (added Entity, err error) {

	index := indexName(entity.TABLE(), entity.KEY())

	if _, ok := dbs.db[index]; !ok {
		dbs.db[index] = NewInMemTable()
	}

	return dbs.db[index].Insert(entity)
}

// Update an existing entity
func (dbs *InMemoryDatastore) Update(entity Entity) (updated Entity, err error) {
	index := indexName(entity.TABLE(), entity.KEY())

	if _, ok := dbs.db[index]; !ok {
		dbs.db[index] = NewInMemTable()
	}

	return dbs.db[index].Update(entity)
}

// Upsert update entity or create it if it does not exist
func (dbs *InMemoryDatastore) Upsert(entity Entity) (updated Entity, err error) {
	index := indexName(entity.TABLE(), entity.KEY())
	if tbl, ok := dbs.db[index]; ok {
		return tbl.Upsert(entity)
	} else {
		return nil, fmt.Errorf(INDEX_NOT_EXISTS)
	}
}

// Delete entity by id and shard (key)
func (dbs *InMemoryDatastore) Delete(factory EntityFactory, entityID string, keys ...string) (err error) {
	index := indexName(factory().TABLE(), keys...)
	if tbl, ok := dbs.db[index]; ok {
		return tbl.Delete(entityID)
	} else {
		return fmt.Errorf(INDEX_NOT_EXISTS)
	}
}

// BulkInsert inserts multiple entities
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

// BulkUpdate updates multiple entities
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

// BulkUpsert update or insert multiple entities
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

// BulkDelete delete multiple entities by IDs
func (dbs *InMemoryDatastore) BulkDelete(factory EntityFactory, entityIDs []string, keys ...string) (affected int64, err error) {
	if len(entityIDs) == 0 {
		return 0, nil
	}
	for _, entityID := range entityIDs {
		if err := dbs.Delete(factory, entityID, keys...); err == nil {
			affected += 1
		}
	}
	return affected, nil
}

// SetField update a single field of the document in a single transaction (eliminates the need to fetch - change - update)
func (dbs *InMemoryDatastore) SetField(factory EntityFactory, entityID string, field string, value any, keys ...string) (err error) {
	fields := make(map[string]any)
	fields[field] = value
	return dbs.SetFields(factory, entityID, fields, keys...)
}

// SetFields update some fields of the document in a single transaction (eliminates the need to fetch - change - update)
func (dbs *InMemoryDatastore) SetFields(factory EntityFactory, entityID string, fields map[string]any, keys ...string) (err error) {

	entity, fe := dbs.Get(factory, entityID, keys...)
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

// BulkSetFields Update specific field of multiple entities in a single transaction (eliminates the need to fetch - change - update)
// The field is the name of the field, values is a map of entityId -> field value
func (dbs *InMemoryDatastore) BulkSetFields(factory EntityFactory, field string, values map[string]any, keys ...string) (affected int64, error error) {

	list, _, fe := dbs.Query(factory).Find(keys...)
	if fe != nil {
		return 0, fe
	}

	// convert entity to Json
	count := 0
	for _, entity := range list {
		js, err := utils.JsonUtils().ToJson(entity)
		if err != nil {
			return 0, err
		}

		// set field
		if val, ok := values[entity.ID()]; ok {
			js[field] = val

			if toSet, _ := utils.JsonUtils().FromJson(factory, js); toSet != nil {
				if _, er := dbs.Update(toSet); er == nil {
					count += 1
				}

			}
		}
	}
	return int64(count), nil
}

// Query is a factory method for query builder Utility
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

// IndexExists tests if index exists
func (dbs *InMemoryDatastore) IndexExists(indexName string) (exists bool) {
	if _, ok := dbs.db[indexName]; ok {
		return true
	} else {
		return false
	}
}

// CreateIndex creates an index (without mapping)
func (dbs *InMemoryDatastore) CreateIndex(indexName string) (name string, err error) {
	// Create index
	if _, ok := dbs.db[indexName]; !ok {
		dbs.db[indexName] = &InMemoryTable{table: make(map[string]Entity)}
	}
	return indexName, nil
}

// CreateEntityIndex creates an index of entity and add entity field mapping
func (dbs *InMemoryDatastore) CreateEntityIndex(factory EntityFactory, key string) (name string, err error) {

	idxName := indexName(factory().TABLE(), key)

	// Create index
	if _, ok := dbs.db[idxName]; !ok {
		dbs.db[idxName] = NewInMemTable()
	}
	return idxName, nil
}

// ListIndices returns a list of all indices matching the pattern
func (dbs *InMemoryDatastore) ListIndices(pattern string) (map[string]int, error) {
	result := make(map[string]int)
	for name, table := range dbs.db {
		if strings.Contains(name, pattern) {
			result[name] = len(table.Table())
		}
	}
	return result, nil
}

// DropIndex drops an index
func (dbs *InMemoryDatastore) DropIndex(indexName string) (ack bool, err error) {

	if _, ok := dbs.db[indexName]; !ok {
		return false, nil
	} else {
		delete(dbs.db, indexName)
		return true, nil
	}
}

// ExecuteQuery Execute native KQL query
func (dbs *InMemoryDatastore) ExecuteQuery(query string, args ...any) ([]Json, error) {
	return nil, fmt.Errorf("not yet implemented")
}

//endregion
