// Copyright 2022. Motty Cohen
//
// Entity interface and base entity for all persistent model entities
//
package entity

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jaevor/go-nanoid"
)

// region Json Document ------------------------------------------------------------------------------------------------

// Represent arbitrary JSON fields collection
type Json map[string]any

// Json document to store in Document object store (Postgres, ElasticSearch, Couchbase ...)
type JsonDoc struct {
	Id   string
	Data string
}

// endregion

// region Timestamp ----------------------------------------------------------------------------------------------------

// Epoch milliseconds Timestamp
type Timestamp int64

// Return current time as Epoch time milliseconds with delta in millis
func EpochNowMillis(delta int64) Timestamp {
	return Timestamp((time.Now().UnixNano() / 1000000) + delta)
}

// Return current time as Epoch time milliseconds with delta in millis
func Now() Timestamp {
	return EpochNowMillis(0)
}

// endregion

// region Entity Interface ---------------------------------------------------------------------------------------------
/**
 * Mark all concrete DB entities
 */
type Entity interface {
	// Return the entity unique Id
	ID() string

	// Return the entity table name (for sharded entities, table name include the suffix of the tenant id)
	TABLE() string

	// Get entity name
	NAME() string

	// Get entity shard key (tenant id)
	KEY() string
}

// Factory method signature for Entity
type EntityFactory func() Entity

// endregion

// region Base Entity --------------------------------------------------------------------------------------------------

/**
 * Base structure for Entity
 */
type BaseEntity struct {
	Id        string    `json:"id"`        // Unique object Id
	Key       string    `json:"key"`       // Shard (tenant) key
	CreatedOn Timestamp `json:"createdOn"` // When the object was created [Epoch milliseconds Timestamp]
	UpdatedOn Timestamp `json:"updatedOn"` // When the object was last updated [Epoch milliseconds Timestamp]
}

func (e BaseEntity) ID() string { return e.Id }

func (e BaseEntity) TABLE() string { return "" }

func (e BaseEntity) NAME() string { return fmt.Sprintf("%s %s", e.TABLE(), e.Id) }

func (e BaseEntity) KEY() string { return "" }

/**
 * Extract table or index name from entity.TABLE()
 */
func EntityIndex(entity Entity, tenantId string) string {

	table := entity.TABLE()

	// Replace templates: {{tenantId}}
	index := strings.Replace(table, "{{tenantId}}", tenantId, -1)

	// Replace templates: {{year}}
	index = strings.Replace(index, "{{year}}", time.Now().Format("2006"), -1)

	// Replace templates: {{month}}
	index = strings.Replace(index, "{{month}}", time.Now().Format("01"), -1)

	return index
}

// endregion

// region Entity Ids ---------------------------------------------------------------------------------------------------
/**
 * Generate new id based on nanoId (faster and smaller than GUID)
 */
func NewId() string {
	if generator, err := nanoid.Standard(21); err != nil {
		return strconv.FormatUint(uint64(time.Now().UnixNano()/1000000), 36)
	} else {
		return generator()
	}

}

/**
 * Generate new Global Unique Identifier
 */
func NewGuid() string {
	return uuid.New().String()
}

// endregion

// region Entity Actions -----------------------------------------------------------------------------------------------

type EntityAction int

const (
	AddEntity    EntityAction = 1
	UpdateEntity EntityAction = 2
	DeleteEntity EntityAction = 3
)

// endregion
