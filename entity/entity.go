// Package entity Entity interface and base entity for all persistent model entities
package entity

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jaevor/go-nanoid"
)

// region Json Document ------------------------------------------------------------------------------------------------

// Json Represent arbitrary JSON fields collection
type Json map[string]any

// JsonDoc is a Json document to store in Document object store (Postgres, ElasticSearch, Couchbase ...)
type JsonDoc struct {
	Id   string
	Data string
}

// endregion

// region Entity Interface ---------------------------------------------------------------------------------------------

// Entity is a marker interface for all serialized domain model entities with identity
type Entity interface {
	// ID return the entity unique Id
	ID() string

	// TABLE return the entity table name (for sharded entities, table name include the suffix of the tenant id)
	TABLE() string

	// NAME return the entity name
	NAME() string

	// KEY return the entity sharding key (tenant/account id) based on one of the entity's attributes
	KEY() string
}

// EntityFactory is the factory method signature for Entity
type EntityFactory func() Entity

// endregion

// region Base Entity --------------------------------------------------------------------------------------------------

// BaseEntity is a base structure for any concrete Entity
type BaseEntity struct {
	Id        string    `json:"id"`        // Unique object Id
	CreatedOn Timestamp `json:"createdOn"` // When the object was created [Epoch milliseconds Timestamp]
	UpdatedOn Timestamp `json:"updatedOn"` // When the object was last updated [Epoch milliseconds Timestamp]
}

func (e BaseEntity) ID() string { return e.Id }

func (e BaseEntity) TABLE() string { return "" }

func (e BaseEntity) NAME() string { return fmt.Sprintf("%s %s", e.TABLE(), e.Id) }

func (e BaseEntity) KEY() string { return "" }

func NewBaseEntity() Entity {
	return &BaseEntity{CreatedOn: Now(), UpdatedOn: Now()}
}

// EntityIndex extract table or index name from entity.TABLE()
func EntityIndex(entity Entity, tenantId string) string {

	table := entity.TABLE()

	// Replace templates: {{tenantId}} or {{accountId}}
	index := strings.Replace(table, "{{tenantId}}", tenantId, -1)
	index = strings.Replace(table, "{{accountId}}", tenantId, -1)

	// Replace templates: {{year}}
	index = strings.Replace(index, "{{year}}", time.Now().Format("2006"), -1)

	// Replace templates: {{month}}
	index = strings.Replace(index, "{{month}}", time.Now().Format("01"), -1)

	return index
}

// endregion

// region Simple Entity ------------------------------------------------------------------------------------------------

// SimpleEntity is a primitive type expressed as an Entity
type SimpleEntity[T any] struct {
	Value T `json:"value"` // entity value
}

func (e SimpleEntity[T]) ID() string { return fmt.Sprintf("%v", e.Value) }

func (e SimpleEntity[T]) TABLE() string { return "" }

func (e SimpleEntity[T]) NAME() string { return fmt.Sprintf("%v", reflect.TypeOf(e.Value).Name()) }

func (e SimpleEntity[T]) KEY() string { return fmt.Sprintf("%v", e.Value) }

func NewSimpleEntity[T any]() Entity {
	return &SimpleEntity[T]{}
}

// endregion

// region Entity Ids ---------------------------------------------------------------------------------------------------
/**
 * Generate new id based on nanoId (faster and smaller than GUID)
 */

// ID return a long string (alphanumeric) based on Epoch micro-seconds in base 36
func ID() string {
	return strconv.FormatUint(uint64(time.Now().UnixMicro()), 36)
}

// IDN return a long string (digits only) based on Epoch micro-seconds
func IDN() string {
	return fmt.Sprintf("%d", time.Now().UnixMicro())
}

// ShortID return a short string (6 characters alphanumeric) based on epoch seconds in base 36
func ShortID(delta ...int) string {
	value := uint64(time.Now().Unix())
	for _, d := range delta {
		value += uint64(d)
	}
	return strconv.FormatUint(value, 36)
}

// ShortIDN return a short string (digits only) based on epoch seconds
func ShortIDN(delta ...int) string {
	value := uint64(time.Now().Unix())
	for _, d := range delta {
		value += uint64(d)
	}
	return fmt.Sprintf("%d", value)
}

// NanoID return a long string (6 characters) based on go-nanoid project (smaller and faster than GUID)
func NanoID() string {
	if generator, err := nanoid.Standard(21); err != nil {
		return strconv.FormatUint(uint64(time.Now().UnixMicro()), 36)
	} else {
		return generator()
	}
}

// GUID generate new Global Unique Identifier
func GUID() string {
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
