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
	NOT_IMPLEMENTED  = "not implemented"
	NOT_SUPPORTED    = "not supported"
	TABLE_NOT_EXISTS = "table does not exist"
)

// region Database store definitions -----------------------------------------------------------------------------------

// InMemoryDatabase represents in memory database with tables
type InMemoryDatabase struct {
	db map[string]ITable
}

// Resolve table name from entity class name and shard keys
func tableName(table string, keys ...string) (tblName string) {

	tblName = table

	if len(keys) == 0 {
		return tblName
	}

	// replace accountId placeholder with the first key
	tblName = strings.Replace(tblName, "{{accountId}}", "{{0}}", -1)

	for idx, key := range keys {
		placeHolder := fmt.Sprintf("{{%d}}", idx)
		tblName = strings.Replace(tblName, placeHolder, key, -1)
	}

	// Replace templates: {{year}}
	tblName = strings.Replace(tblName, "{{year}}", time.Now().Format("2006"), -1)

	// Replace templates: {{month}}
	tblName = strings.Replace(tblName, "{{month}}", time.Now().Format("01"), -1)

	// TODO: Replace templates: {{week}}

	return
}

// endregion

// region Factory and connectivity methods for Database ----------------------------------------------------------------

// NewInMemoryDatabase Factory method for database
func NewInMemoryDatabase() (dbs IDatabase, err error) {
	return &InMemoryDatabase{db: make(map[string]ITable)}, nil
}

// Ping Test database connectivity
// @param retries - how many retries are required (max 10)
// @param interval - time interval (in seconds) between retries (max 60)
func (dbs *InMemoryDatabase) Ping(retries uint, interval uint) error {
	return nil
}

// Close DB and free resources
func (dbs *InMemoryDatabase) Close() error {
	logger.Debug("In memory database closed")
	return nil
}

//endregion

// region Database store Basic CRUD methods ----------------------------------------------------------------------------

// Get single entity by ID
func (dbs *InMemoryDatabase) Get(factory EntityFactory, entityID string, keys ...string) (result Entity, err error) {

	entity := factory()
	table := tableName(entity.TABLE(), keys...)
	if tbl, ok := dbs.db[table]; ok {
		return tbl.Get(entityID)
	} else {
		return nil, fmt.Errorf(TABLE_NOT_EXISTS)
	}
}

// List get a list of entities by IDs
func (dbs *InMemoryDatabase) List(factory EntityFactory, entityIDs []string, keys ...string) (list []Entity, err error) {

	entity := factory()
	table := tableName(entity.TABLE(), keys...)

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

// Exists checks if entity exists by ID
func (dbs *InMemoryDatabase) Exists(factory EntityFactory, entityID string, keys ...string) (result bool, err error) {

	entity := factory()
	table := tableName(entity.TABLE(), keys...)

	if tbl, ok := dbs.db[table]; ok {
		return tbl.Exists(entityID)
	} else {
		return false, fmt.Errorf(TABLE_NOT_EXISTS)
	}
}

// Insert Add new entity
func (dbs *InMemoryDatabase) Insert(entity Entity) (added Entity, err error) {

	table := tableName(entity.TABLE(), entity.KEY())

	if _, ok := dbs.db[table]; !ok {
		dbs.db[table] = NewInMemTable()
	}
	return dbs.db[table].Insert(entity)
}

// Update existing entity in the data store
func (dbs *InMemoryDatabase) Update(entity Entity) (updated Entity, err error) {
	table := tableName(entity.TABLE(), entity.KEY())

	if _, ok := dbs.db[table]; !ok {
		dbs.db[table] = NewInMemTable()
	}
	return dbs.db[table].Update(entity)
}

// Upsert updates existing entity in the data store or add it if it does not exist
func (dbs *InMemoryDatabase) Upsert(entity Entity) (updated Entity, err error) {
	table := tableName(entity.TABLE(), entity.KEY())
	if tbl, ok := dbs.db[table]; ok {
		return tbl.Upsert(entity)
	} else {
		return nil, fmt.Errorf(TABLE_NOT_EXISTS)
	}
}

// Delete entity by id
func (dbs *InMemoryDatabase) Delete(factory EntityFactory, entityID string, keys ...string) (err error) {

	entity := factory()

	table := tableName(entity.TABLE(), keys...)
	if tbl, ok := dbs.db[table]; ok {
		return tbl.Delete(entityID)
	} else {
		return fmt.Errorf(TABLE_NOT_EXISTS)
	}
}

// BulkInsert adds multiple entities to data store (all must be of the same type)
func (dbs *InMemoryDatabase) BulkInsert(entities []Entity) (affected int64, err error) {
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

// BulkUpdate updates multiple entities in the data store (all must be of the same type)
func (dbs *InMemoryDatabase) BulkUpdate(entities []Entity) (affected int64, err error) {
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

// BulkUpsert update or insert multiple entities in the data store (all must be of the same type)
func (dbs *InMemoryDatabase) BulkUpsert(entities []Entity) (affected int64, err error) {
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

// BulkDelete delete multiple entities by id list
func (dbs *InMemoryDatabase) BulkDelete(factory EntityFactory, entityIDs []string, keys ...string) (affected int64, err error) {
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

// SetField updates single field of the document in a single transaction (eliminates the need to fetch - change - update)
func (dbs *InMemoryDatabase) SetField(factory EntityFactory, entityID string, field string, value any, keys ...string) (err error) {
	fields := make(map[string]any)
	fields[field] = value
	return dbs.SetFields(factory, entityID, fields, keys...)
}

// SetFields Updates some numeric fields of the document in a single transaction (eliminates the need to fetch - change - update)
func (dbs *InMemoryDatabase) SetFields(factory EntityFactory, entityID string, fields map[string]any, keys ...string) (err error) {

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
func (dbs *InMemoryDatabase) BulkSetFields(factory EntityFactory, field string, values map[string]any, keys ...string) (affected int64, error error) {

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

// Query is a builder method to construct query
func (dbs *InMemoryDatabase) Query(factory EntityFactory) IQuery {
	return &inMemoryDatabaseQuery{
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

// region Database DDL and DML -----------------------------------------------------------------------------------------

// ExecuteDDL create table and indexes
// The ddl parameter is a map of strings (table names) to array of strings (list of fields to index)
func (dbs *InMemoryDatabase) ExecuteDDL(ddl map[string][]string) (err error) {

	for table, fields := range ddl {
		logger.Debug("Creating table: %s with fields indexes: %s", table, strings.Join(fields, ","))

		if _, ok := dbs.db[table]; !ok {
			dbs.db[table] = &InMemoryTable{table: make(map[string]Entity)}
		}
	}
	return nil
}

// ExecuteSQL execute raw SQL command
func (dbs *InMemoryDatabase) ExecuteSQL(sql string, args ...any) (affected int64, err error) {
	return 0, fmt.Errorf(NOT_SUPPORTED)
}

// DropTable drop a table and its related indexes
func (dbs *InMemoryDatabase) DropTable(table string) (err error) {
	delete(dbs.db, table)
	return nil
}

// PurgeTable fast delete table content (truncate)
func (dbs *InMemoryDatabase) PurgeTable(table string) (err error) {
	delete(dbs.db, table)
	return nil
}

// endregion
