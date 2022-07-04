// Copyright 2022. Motty Cohen
//
// In-memory datastore implementation of IDatabase (used for testing)
//
package database

import (
	"fmt"
	. "github.com/mottyc/yaaf-common/entity"
)

// Database table interface
type ITable interface {

	// Get single entity by ID
	Get(entityID string) (entity Entity, err error)

	// Check if entity exists by ID
	Exists(entityID string) (result bool, err error)

	// Insert entity
	Insert(entity Entity) (added Entity, err error)

	// Update entity
	Update(entity Entity) (added Entity, err error)

	// Update entity or insert if not found
	Upsert(entity Entity) (added Entity, err error)

	// Delete entity
	Delete(entityID string) (err error)

	// Get access to the underlying data structure
	Table() (result map[string]Entity)
}

// In memory table implementation --------------------------------------------------------------------------------------

// Represent a table in the DB
type InMemoryTable struct {
	table map[string]Entity
}

func NewInMemTable() ITable {
	return &InMemoryTable{table: make(map[string]Entity)}
}

// Get single entity by ID
func (tbl *InMemoryTable) Get(entityID string) (entity Entity, err error) {
	if ent, ok := tbl.table[entityID]; ok {
		return ent, nil
	} else {
		return nil, fmt.Errorf("item not found")
	}
}

// Check if entity exists by ID
func (tbl *InMemoryTable) Exists(entityID string) (result bool, err error) {
	_, ok := tbl.table[entityID]
	return ok, nil
}

// Insert entity
func (tbl *InMemoryTable) Insert(entity Entity) (added Entity, err error) {
	entityID := fmt.Sprintf("%v", entity.ID())
	if _, ok := tbl.table[entityID]; ok {
		return nil, fmt.Errorf("item exists")
	} else {
		tbl.table[entityID] = entity
		return entity, nil
	}
}

// Update entity
func (tbl *InMemoryTable) Update(entity Entity) (added Entity, err error) {
	entityID := fmt.Sprintf("%v", entity.ID())
	if _, ok := tbl.table[entityID]; ok {
		tbl.table[entityID] = entity
		return entity, nil
	} else {
		return nil, fmt.Errorf("item not exists")
	}
}

// Update entity or insert if not found
func (tbl *InMemoryTable) Upsert(entity Entity) (added Entity, err error) {
	entityID := fmt.Sprintf("%v", entity.ID())
	tbl.table[entityID] = entity
	return entity, nil
}

// Delete entity
func (tbl *InMemoryTable) Delete(entityID string) (err error) {
	if _, ok := tbl.table[entityID]; ok {
		delete(tbl.table, entityID)
		return nil
	} else {
		return fmt.Errorf("item not found")
	}
}

// Get access to the underlying data structure
func (tbl *InMemoryTable) Table() (result map[string]Entity) {
	return tbl.table
}
