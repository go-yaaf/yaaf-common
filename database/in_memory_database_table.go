package database

import (
	"fmt"

	. "github.com/go-yaaf/yaaf-common/entity"
)

// ITable defines the interface for a database table.
// It provides methods for basic CRUD operations on entities within the table.
type ITable interface {

	// Get retrieves a single entity by its ID.
	Get(entityID string) (entity Entity, err error)

	// Exists checks if an entity exists by its ID.
	Exists(entityID string) (result bool, err error)

	// Insert inserts a new entity into the table.
	Insert(entity Entity) (added Entity, err error)

	// Update updates an existing entity in the table.
	Update(entity Entity) (added Entity, err error)

	// Upsert inserts or updates an entity in the table.
	Upsert(entity Entity) (added Entity, err error)

	// Delete deletes an entity from the table by its ID.
	Delete(entityID string) (err error)

	// Table returns the underlying map storage of the table.
	Table() (result map[string]Entity)
}

// In memory DbTable implementation --------------------------------------------------------------------------------------

// InMemoryTable represents an in-memory implementation of the ITable interface.
// It uses a map to store entities.
type InMemoryTable struct {
	DbTable map[string]Entity `json:"dbTable"` // DbTable is the underlying map storage.
}

// NewInMemTable creates a new instance of InMemoryTable.
func NewInMemTable() ITable {
	return &InMemoryTable{DbTable: make(map[string]Entity)}
}

// Get retrieves a single entity by its ID.
func (tbl *InMemoryTable) Get(entityID string) (entity Entity, err error) {
	if ent, ok := tbl.DbTable[entityID]; ok {
		return ent, nil
	} else {
		return nil, fmt.Errorf("item not found")
	}
}

// Exists checks if an entity exists by its ID.
func (tbl *InMemoryTable) Exists(entityID string) (result bool, err error) {
	_, ok := tbl.DbTable[entityID]
	return ok, nil
}

// Insert inserts a new entity into the table.
func (tbl *InMemoryTable) Insert(entity Entity) (added Entity, err error) {
	entityID := fmt.Sprintf("%v", entity.ID())
	if _, ok := tbl.DbTable[entityID]; ok {
		return nil, fmt.Errorf("item exists")
	} else {
		tbl.DbTable[entityID] = entity
		return entity, nil
	}
}

// Update updates an existing entity in the table.
func (tbl *InMemoryTable) Update(entity Entity) (added Entity, err error) {
	entityID := fmt.Sprintf("%v", entity.ID())
	if _, ok := tbl.DbTable[entityID]; ok {
		tbl.DbTable[entityID] = entity
		return entity, nil
	} else {
		return nil, fmt.Errorf("item not exists")
	}
}

// Upsert inserts or updates an entity in the table.
func (tbl *InMemoryTable) Upsert(entity Entity) (added Entity, err error) {
	entityID := fmt.Sprintf("%v", entity.ID())
	tbl.DbTable[entityID] = entity
	return entity, nil
}

// Delete deletes an entity from the table by its ID.
func (tbl *InMemoryTable) Delete(entityID string) (err error) {
	if _, ok := tbl.DbTable[entityID]; ok {
		delete(tbl.DbTable, entityID)
		return nil
	} else {
		return fmt.Errorf("item not found")
	}
}

// Table returns the underlying map storage of the table.
func (tbl *InMemoryTable) Table() (result map[string]Entity) {
	return tbl.DbTable
}
