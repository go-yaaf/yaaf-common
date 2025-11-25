package database

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	. "github.com/go-yaaf/yaaf-common/entity"
	"github.com/go-yaaf/yaaf-common/logger"
	"github.com/go-yaaf/yaaf-common/utils"
)

const (
	NOT_IMPLEMENTED  = "not implemented"
	NOT_SUPPORTED    = "not supported"
	TABLE_NOT_EXISTS = "DbTable does not exist"
)

// region Database store definitions -----------------------------------------------------------------------------------

// InMemoryDatabase represents an in-memory database implementation using maps.
// It implements the IDatabase interface and is primarily used for testing and development.
type InMemoryDatabase struct {
	Db map[string]ITable `json:"db"` // Db is a map of table names to ITable instances.
}

// resolveTableName determines the table name for a given entity.
// It handles shard keys and template replacements (e.g., {{year}}, {{month}}).
func (d *InMemoryDatabase) resolveTableName(entity Entity, keys ...string) string {
	return tableName(entity.TABLE(), keys...)
}

// tableName resolves the table name from entity class name and shard keys.
// It is a package-level helper used by other in-memory components.
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

	return tblName
}

// endregion

// region Factory and connectivity methods for Database ----------------------------------------------------------------

// NewInMemoryDatabase creates a new instance of InMemoryDatabase.
func NewInMemoryDatabase() (IDatabase, error) {
	return &InMemoryDatabase{
		Db: make(map[string]ITable),
	}, nil
}

// Ping tests the database connectivity (always returns nil for in-memory database).
func (d *InMemoryDatabase) Ping(retries uint, intervalInSeconds uint) error {
	return nil
}

// Close closes the database connection (no-op for in-memory database).
func (d *InMemoryDatabase) Close() error {
	logger.Debug("In memory database closed")
	return nil
}

// CloneDatabase creates a copy of the current database instance.
// Note: This performs a shallow copy of the database structure.
func (d *InMemoryDatabase) CloneDatabase() (IDatabase, error) {
	return &InMemoryDatabase{
		Db: d.Db,
	}, nil
}

// endregion

// region Database store Basic CRUD methods ----------------------------------------------------------------------------

// Get retrieves a single entity by its ID.
func (d *InMemoryDatabase) Get(factory EntityFactory, entityID string, keys ...string) (result Entity, err error) {
	entity := factory()
	tableName := d.resolveTableName(entity, keys...)
	if table, ok := d.Db[tableName]; ok {
		return table.Get(entityID)
	} else {
		return nil, fmt.Errorf("table not found: %s", tableName)
	}
}

// List retrieves a list of entities by their IDs.
func (d *InMemoryDatabase) List(factory EntityFactory, entityIDs []string, keys ...string) (list []Entity, err error) {
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
		return list, fmt.Errorf(TABLE_NOT_EXISTS)
	}
	return
}

// Exists checks if an entity exists by its ID.
func (d *InMemoryDatabase) Exists(factory EntityFactory, entityID string, keys ...string) (result bool, err error) {
	entity := factory()
	tableName := d.resolveTableName(entity, keys...)
	if table, ok := d.Db[tableName]; ok {
		if _, err := table.Get(entityID); err == nil {
			return true, nil
		}
	}
	return false, nil
}

// Insert inserts a new entity into the database.
func (d *InMemoryDatabase) Insert(entity Entity) (added Entity, err error) {
	tableName := d.resolveTableName(entity, entity.KEY())
	if _, ok := d.Db[tableName]; !ok {
		d.Db[tableName] = NewInMemTable()
	}
	return d.Db[tableName].Insert(entity)
}

// Update updates an existing entity in the database.
func (d *InMemoryDatabase) Update(entity Entity) (updated Entity, err error) {
	tableName := d.resolveTableName(entity, entity.KEY())
	if table, ok := d.Db[tableName]; ok {
		return table.Update(entity)
	} else {
		return nil, fmt.Errorf("table not found: %s", tableName)
	}
}

// Upsert inserts or updates an entity in the database.
func (d *InMemoryDatabase) Upsert(entity Entity) (updated Entity, err error) {
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

// Delete deletes an entity from the database by its ID.
func (d *InMemoryDatabase) Delete(factory EntityFactory, entityID string, keys ...string) (err error) {
	entity := factory()
	tableName := d.resolveTableName(entity, keys...)
	if table, ok := d.Db[tableName]; ok {
		return table.Delete(entityID)
	} else {
		return fmt.Errorf("table not found: %s", tableName)
	}
}

// BulkInsert inserts multiple entities into the database.
func (d *InMemoryDatabase) BulkInsert(entities []Entity) (affected int64, err error) {
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

// BulkUpdate updates multiple entities in the database.
func (d *InMemoryDatabase) BulkUpdate(entities []Entity) (affected int64, err error) {
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

// BulkUpsert inserts or updates multiple entities in the database.
func (d *InMemoryDatabase) BulkUpsert(entities []Entity) (affected int64, err error) {
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

// BulkDelete deletes multiple entities from the database by their IDs.
func (d *InMemoryDatabase) BulkDelete(factory EntityFactory, entityIDs []string, keys ...string) (affected int64, err error) {
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
func (d *InMemoryDatabase) SetField(factory EntityFactory, entityID string, field string, value any, keys ...string) (err error) {
	fields := make(map[string]any)
	fields[field] = value
	return d.SetFields(factory, entityID, fields, keys...)
}

// SetFields updates multiple fields of an entity.
func (d *InMemoryDatabase) SetFields(factory EntityFactory, entityID string, fields map[string]any, keys ...string) (err error) {
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
func (d *InMemoryDatabase) BulkSetFields(factory EntityFactory, field string, values map[string]any, keys ...string) (affected int64, error error) {
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
func (d *InMemoryDatabase) Query(factory EntityFactory) IQuery {
	return &inMemoryDatabaseQuery{
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

// AdvancedQuery creates a new advanced query builder (not implemented).
func (d *InMemoryDatabase) AdvancedQuery(factory EntityFactory) IAdvancedQuery {
	panic("InMemoryDatabase: IAdvancedQuery is not implemented/supported")
}

// endregion

// region Database DDL and DML -----------------------------------------------------------------------------------------

// ExecuteDDL executes a Data Definition Language (DDL) query.
// The ddl parameter is a map of strings (table names) to array of strings (list of fields to index).
func (d *InMemoryDatabase) ExecuteDDL(ddl map[string][]string) (err error) {
	for table, fields := range ddl {
		logger.Debug("Creating DbTable: %s with fields indexes: %s", table, strings.Join(fields, ","))

		if _, ok := d.Db[table]; !ok {
			d.Db[table] = &InMemoryTable{DbTable: make(map[string]Entity)}
		}
	}
	return nil
}

// ExecuteSQL executes a raw SQL query (not supported).
func (d *InMemoryDatabase) ExecuteSQL(sql string, args ...any) (affected int64, err error) {
	return 0, fmt.Errorf(NOT_SUPPORTED)
}

// ExecuteQuery executes a native SQL query (not supported).
func (d *InMemoryDatabase) ExecuteQuery(source, sql string, args ...any) ([]Json, error) {
	return nil, fmt.Errorf(NOT_SUPPORTED)
}

// DropTable drops a table and its related indexes.
func (d *InMemoryDatabase) DropTable(table string) (err error) {
	delete(d.Db, table)
	return nil
}

// PurgeTable removes all data from a table (truncate).
func (d *InMemoryDatabase) PurgeTable(table string) (err error) {
	delete(d.Db, table)
	return nil
}

// endregion

// region Backup and Restore Database ----------------------------------------------------------------

// TODO: Restore is not working properly since table holds Entity interface pointers

// Backup creates a backup of the database to a file.
func (d *InMemoryDatabase) Backup(path string) error {
	if strings.HasSuffix(path, ".json") {
		return d.backupJson(path)
	} else {
		return d.backupBinary(path)
	}
}

// Restore restores the database from a backup file.
func (d *InMemoryDatabase) Restore(path string) error {
	if strings.HasSuffix(path, ".json") {
		return d.restoreJson(path)
	} else {
		return d.restoreBinary(path)
	}
}

// backupJson backs up the database to a file in JSON format.
func (d *InMemoryDatabase) backupJson(path string) error {
	// Ensure path
	folders := filepath.Dir(path)
	if err := os.MkdirAll(folders, 0755); err != nil {
		return err
	}

	// create file
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	err = enc.Encode(d.Db)
	_ = file.Close()
	return err
}

// backupBinary backs up the database to a file in binary format.
func (d *InMemoryDatabase) backupBinary(path string) error {
	// Ensure path
	folders := filepath.Dir(path)
	if err := os.MkdirAll(folders, 0755); err != nil {
		return err
	}

	// create file
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	enc := gob.NewEncoder(file)
	err = enc.Encode(d.Db)
	_ = file.Close()
	return err
}

// restoreJson restores the database from a JSON file.
func (d *InMemoryDatabase) restoreJson(path string) error {
	// open file
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	dec := json.NewDecoder(file)

	dst := make(map[string]*InMemoryTable)
	err = dec.Decode(&dst)
	_ = file.Close()
	return err
}

// restoreBinary restores the database from a binary file.
func (d *InMemoryDatabase) restoreBinary(path string) error {
	// open file
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	dec := gob.NewDecoder(file)
	err = dec.Decode(&d.Db)
	_ = file.Close()
	return err
}

// endregion
