package database

import (
	"fmt"
	"strings"

	. "github.com/go-yaaf/yaaf-common/entity"
	"github.com/go-yaaf/yaaf-common/logger"
	"github.com/go-yaaf/yaaf-common/utils"
)

const (
	INDEX_NOT_EXISTS = "index does not exist"
)

// region Datastore definitions ----------------------------------------------------------------------------------------

// InMemoryDatastore represents an in-memory implementation of the IDatastore interface.
// It uses a map to store data (simulating tables/collections) and is primarily used for testing and development.
type InMemoryDatastore struct {
	Db map[string]ITable `json:"db"` // Db is a map of collection names to ITable instances.
}

// resolveTableName determines the table name for a given entity.
func (d *InMemoryDatastore) resolveTableName(entity Entity, keys ...string) string {
	return tableName(entity.TABLE(), keys...)
}

// indexName resolves index name from entity name (alias for tableName for backward compatibility).
func indexName(table string, keys ...string) (name string) {
	return tableName(table, keys...)
}

// endregion

// region Factory and connectivity methods for Datastore ---------------------------------------------------------------

// NewInMemoryDatastore creates a new instance of InMemoryDatastore.
func NewInMemoryDatastore() (IDatastore, error) {
	return &InMemoryDatastore{
		Db: make(map[string]ITable),
	}, nil
}

// Ping tests the datastore connectivity (always returns nil for in-memory datastore).
func (d *InMemoryDatastore) Ping(retries uint, intervalInSeconds uint) error {
	logger.Debug("Pinging %d times with %d interval", retries, intervalInSeconds)
	return nil
}

// Close closes the datastore connection (no-op for in-memory datastore).
func (d *InMemoryDatastore) Close() error {
	logger.Debug("In memory datastore closed")
	return nil
}

// CloneDatastore creates a copy of the current datastore instance.
func (d *InMemoryDatastore) CloneDatastore() (IDatastore, error) {
	return &InMemoryDatastore{
		Db: d.Db,
	}, nil
}

// endregion

// region Datastore Basic CRUD methods ---------------------------------------------------------------------------------

// Get retrieves a single entity by its ID.
func (d *InMemoryDatastore) Get(factory EntityFactory, entityID string, keys ...string) (result Entity, err error) {
	entity := factory()
	tableName := d.resolveTableName(entity, keys...)
	if table, ok := d.Db[tableName]; ok {
		return table.Get(entityID)
	} else {
		return nil, fmt.Errorf("collection not found: %s", tableName)
	}
}

// List retrieves a list of entities by their IDs.
func (d *InMemoryDatastore) List(factory EntityFactory, entityIDs []string, keys ...string) (list []Entity, err error) {
	entity := factory()
	tableName := d.resolveTableName(entity, keys...)

	list = make([]Entity, 0)

	if table, ok := d.Db[tableName]; ok {
		for _, id := range entityIDs {
			if ent, err := table.Get(id); err == nil {
				list = append(list, ent)
			}
		}
	} else {
		return list, fmt.Errorf("collection not found: %s", tableName)
	}
	return
}

// Exists checks if an entity exists by its ID.
func (d *InMemoryDatastore) Exists(factory EntityFactory, entityID string, keys ...string) (result bool, err error) {
	entity := factory()
	tableName := d.resolveTableName(entity, keys...)
	if table, ok := d.Db[tableName]; ok {
		if _, err := table.Get(entityID); err == nil {
			return true, nil
		}
	}
	return false, nil
}

// Insert inserts a new entity into the datastore.
func (d *InMemoryDatastore) Insert(entity Entity) (added Entity, err error) {
	tableName := d.resolveTableName(entity, entity.KEY())
	if _, ok := d.Db[tableName]; !ok {
		d.Db[tableName] = NewInMemTable()
	}
	return d.Db[tableName].Insert(entity)
}

// Update updates an existing entity in the datastore.
func (d *InMemoryDatastore) Update(entity Entity) (updated Entity, err error) {
	tableName := d.resolveTableName(entity, entity.KEY())
	if table, ok := d.Db[tableName]; ok {
		return table.Update(entity)
	} else {
		return nil, fmt.Errorf("collection not found: %s", tableName)
	}
}

// Upsert inserts or updates an entity in the datastore.
func (d *InMemoryDatastore) Upsert(entity Entity) (updated Entity, err error) {
	tableName := d.resolveTableName(entity, entity.KEY())
	if _, ok := d.Db[tableName]; !ok {
		d.Db[tableName] = NewInMemTable()
	}
	if _, err := d.Db[tableName].Get(entity.ID()); err == nil {
		return d.Db[tableName].Update(entity)
	} else {
		return d.Db[tableName].Insert(entity)
	}
}

// Delete deletes an entity from the datastore by its ID.
func (d *InMemoryDatastore) Delete(factory EntityFactory, entityID string, keys ...string) (err error) {
	entity := factory()
	tableName := d.resolveTableName(entity, keys...)
	if table, ok := d.Db[tableName]; ok {
		return table.Delete(entityID)
	} else {
		return fmt.Errorf("collection not found: %s", tableName)
	}
}

// BulkInsert inserts multiple entities into the datastore.
func (d *InMemoryDatastore) BulkInsert(entities []Entity) (affected int64, err error) {
	if len(entities) == 0 {
		return 0, nil
	}
	for _, ent := range entities {
		if _, err := d.Insert(ent); err == nil {
			affected += 1
		}
	}
	return affected, nil
}

// BulkUpdate updates multiple entities in the datastore.
func (d *InMemoryDatastore) BulkUpdate(entities []Entity) (affected int64, err error) {
	if len(entities) == 0 {
		return 0, nil
	}
	for _, ent := range entities {
		if _, err := d.Update(ent); err == nil {
			affected += 1
		}
	}
	return affected, nil
}

// BulkUpsert inserts or updates multiple entities in the datastore.
func (d *InMemoryDatastore) BulkUpsert(entities []Entity) (affected int64, err error) {
	if len(entities) == 0 {
		return 0, nil
	}
	for _, ent := range entities {
		if _, err := d.Upsert(ent); err == nil {
			affected += 1
		}
	}
	return affected, nil
}

// BulkDelete deletes multiple entities from the datastore by their IDs.
func (d *InMemoryDatastore) BulkDelete(factory EntityFactory, entityIDs []string, keys ...string) (affected int64, err error) {
	if len(entityIDs) == 0 {
		return 0, nil
	}
	for _, entityID := range entityIDs {
		if err := d.Delete(factory, entityID, keys...); err == nil {
			affected += 1
		}
	}
	return affected, nil
}

// SetField updates a single field of an entity.
func (d *InMemoryDatastore) SetField(factory EntityFactory, entityID string, field string, value any, keys ...string) (err error) {
	fields := make(map[string]any)
	fields[field] = value
	return d.SetFields(factory, entityID, fields, keys...)
}

// SetFields updates multiple fields of an entity.
func (d *InMemoryDatastore) SetFields(factory EntityFactory, entityID string, fields map[string]any, keys ...string) (err error) {
	entity, fe := d.Get(factory, entityID, keys...)
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

	_, fe = d.Update(toSet)
	return fe
}

// BulkSetFields updates a specific field of multiple entities in a single transaction.
func (d *InMemoryDatastore) BulkSetFields(factory EntityFactory, field string, values map[string]any, keys ...string) (affected int64, error error) {
	list, _, fe := d.Query(factory).Find(keys...)
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
				if _, er := d.Update(toSet); er == nil {
					count += 1
				}
			}
		}
	}
	return int64(count), nil
}

// Query creates a new query builder for the specified entity factory.
func (d *InMemoryDatastore) Query(factory EntityFactory) IQuery {
	return &inMemoryDatastoreQuery{
		db:         d,
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
func (d *InMemoryDatastore) IndexExists(indexName string) (exists bool) {
	if _, ok := d.Db[indexName]; ok {
		return true
	} else {
		return false
	}
}

// CreateIndex creates an index on the specified fields (no-op for in-memory datastore).
func (d *InMemoryDatastore) CreateIndex(indexName string) (name string, err error) {
	// Create index
	if _, ok := d.Db[indexName]; !ok {
		d.Db[indexName] = &InMemoryTable{DbTable: make(map[string]Entity)}
	}
	return indexName, nil
}

// CreateEntityIndex creates an index for the entity (no-op for in-memory datastore).
func (d *InMemoryDatastore) CreateEntityIndex(factory EntityFactory, key string) (name string, err error) {
	idxName := d.resolveTableName(factory(), key)
	if _, ok := d.Db[idxName]; !ok {
		d.Db[idxName] = NewInMemTable()
	}
	return idxName, nil
}

// ListIndices retrieves a list of existing indices (returns empty list for in-memory datastore).
func (d *InMemoryDatastore) ListIndices(pattern string) (map[string]int, error) {
	result := make(map[string]int)
	for name, table := range d.Db {
		if strings.Contains(name, pattern) {
			result[name] = len(table.Table())
		}
	}
	return result, nil
}

// DropIndex drops an index (no-op for in-memory datastore).
func (d *InMemoryDatastore) DropIndex(indexName string) (ack bool, err error) {
	if _, ok := d.Db[indexName]; !ok {
		return false, nil
	} else {
		delete(d.Db, indexName)
		return true, nil
	}
}

// PurgeIndex removes all data from an index (truncate).
func (d *InMemoryDatastore) PurgeIndex(indexName string) (ack bool, err error) {
	if _, ok := d.Db[indexName]; !ok {
		return false, nil
	} else {
		delete(d.Db, indexName)
		return true, nil
	}
}

// ExecuteQuery executes a native query (not implemented).
func (d *InMemoryDatastore) ExecuteQuery(source string, query string, args ...any) ([]Json, error) {
	return nil, fmt.Errorf("not yet implemented")
}

// endregion
