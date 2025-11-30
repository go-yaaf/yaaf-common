// Package entity Entity interface and base entity for all persistent model entities
package entity

import (
	"encoding/json"
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
// It is a map of string to any, used for flexible data structures.
// @Data
type Json map[string]any

// JsonDoc is a Json document to store in Document object store (Postgres, ElasticSearch, Couchbase ...)
// It serves as a generic container for JSON data with a unique identifier.
// @Data
type JsonDoc struct {
	Id   string // Id is the unique identifier of the document
	Data string // Data is the raw JSON string content
}

// endregion

// region Entity Interface ---------------------------------------------------------------------------------------------

// Entity is a marker interface for all serialized domain model entities with identity.
// Any struct that needs to be persisted or manipulated as a domain entity must implement this interface.
type Entity interface {
	// ID returns the entity unique Id.
	// This ID is used to uniquely identify the entity within its table or collection.
	ID() string

	// TABLE returns the entity table name.
	// For sharded entities, the table name may include the suffix of the tenant id.
	// This is used by the storage layer to determine where to persist the entity.
	TABLE() string

	// NAME returns the entity name.
	// This is typically used for logging, display, or logical identification of the entity type.
	NAME() string

	// KEY returns the entity sharding key (tenant/account id) based on one of the entity's attributes.
	// This key is crucial for distributed systems where data is partitioned by tenant or account.
	KEY() string
}

// EntityFactory is the factory method signature for creating new Entity instances.
// It is used by the framework to instantiate entities dynamically, for example when unmarshalling from a database.
type EntityFactory func() Entity

// endregion

// region Base Entity --------------------------------------------------------------------------------------------------

// BaseEntity is a base structure for any concrete Entity.
// It provides common fields like Id, CreatedOn, and UpdatedOn that are standard across most entities.
// Embed this struct in your domain entities to inherit these standard fields and basic behavior.
// @Data
type BaseEntity struct {
	Id        string    `json:"id"`        // Id is the unique object identifier
	CreatedOn Timestamp `json:"createdOn"` // CreatedOn is the timestamp when the object was created [Epoch milliseconds Timestamp]
	UpdatedOn Timestamp `json:"updatedOn"` // UpdatedOn is the timestamp when the object was last updated [Epoch milliseconds Timestamp]
}

func (e *BaseEntity) ID() string { return e.Id }

func (e *BaseEntity) TABLE() string { return "" }

func (e *BaseEntity) NAME() string { return fmt.Sprintf("%s %s", e.TABLE(), e.Id) }

func (e *BaseEntity) KEY() string { return "" }

// BaseAnalyticEntity is a base structure for analytical entities.
// Unlike BaseEntity, it does not enforce a standard ID or timestamp fields, making it suitable for
// data points, logs, or other analytical data that might have different identification or timing requirements.
type BaseAnalyticEntity struct {
}

func (e *BaseAnalyticEntity) ID() string { return "" }

func (e *BaseAnalyticEntity) TABLE() string { return "" }

func (e *BaseAnalyticEntity) NAME() string { return "" }

func (e *BaseAnalyticEntity) KEY() string { return "" }

// NewBaseEntity creates a new BaseEntity instance with the current time for CreatedOn and UpdatedOn.
// This is a helper function to initialize a BaseEntity with valid timestamps.
func NewBaseEntity() Entity {
	return &BaseEntity{CreatedOn: Now(), UpdatedOn: Now()}
}

// EntityIndex extracts the table or index name from entity.TABLE() and applies template replacements.
// It replaces placeholders like {{tenantId}}, {{accountId}}, {{year}}, and {{month}} with actual values.
// This is useful for dynamic table naming strategies, such as time-based or tenant-based sharding.
//
// Parameters:
//   - entity: The entity instance to derive the index from.
//   - tenantId: The tenant ID to replace {{tenantId}} or {{accountId}} placeholders.
//
// Returns:
//   - The resolved index or table name.
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

// region Base Entity Ex -----------------------------------------------------------------------------------------------

// BaseEntityEx is an extended base Entity with custom attributes.
// It includes all fields from BaseEntity and adds support for a status Flag and a generic properties map (Props).
// This is useful for entities that need soft-delete capabilities (via Flag) or extensible fields (via Props).
// @Data
type BaseEntityEx struct {
	Id        string    `json:"id"`        // Id is the unique object identifier
	CreatedOn Timestamp `json:"createdOn"` // CreatedOn is the timestamp when the object was created [Epoch milliseconds Timestamp]
	UpdatedOn Timestamp `json:"updatedOn"` // UpdatedOn is the timestamp when the object was last updated [Epoch milliseconds Timestamp]
	Flag      int64     `json:"flag"`      // Flag represents the entity status (e.g. -1 = deleted, 0 = active)
	Props     Json      `json:"props"`     // Props is a map of custom properties for extensibility
}

func (e *BaseEntityEx) ID() string { return e.Id }

func (e *BaseEntityEx) TABLE() string { return "" }

func (e *BaseEntityEx) NAME() string { return fmt.Sprintf("%s %s", e.TABLE(), e.Id) }

func (e *BaseEntityEx) KEY() string { return "" }

// NewBaseEntityEx creates a new BaseEntityEx instance with initialized timestamps and an empty Props map.
func NewBaseEntityEx() Entity {
	return &BaseEntityEx{CreatedOn: Now(), UpdatedOn: Now(), Props: Json{}}
}

// endregion

// region Simple Entity ------------------------------------------------------------------------------------------------

// SimpleEntity is a generic wrapper to express a primitive type as an Entity.
// It is useful when you need to treat a simple value (like a string or int) as a full-fledged Entity
// in the system, for example when passing it to functions that expect an Entity interface.
// @Data
type SimpleEntity[T any] struct {
	Value T `json:"value"` // Value is the wrapped entity value
}

func (e *SimpleEntity[T]) ID() string { return fmt.Sprintf("%v", e.Value) }

func (e *SimpleEntity[T]) TABLE() string { return "" }

func (e *SimpleEntity[T]) NAME() string { return fmt.Sprintf("%v", reflect.TypeOf(e.Value).Name()) }

func (e *SimpleEntity[T]) KEY() string { return fmt.Sprintf("%v", e.Value) }

// NewSimpleEntity creates a new SimpleEntity instance.
func NewSimpleEntity[T any]() Entity {
	return &SimpleEntity[T]{}
}

// endregion

// region Entities -----------------------------------------------------------------------------------------------------

// Entities is a generic wrapper to express a list of primitive or complex types as an Entity.
// It allows a collection of items to be treated as a single Entity unit.
// @Data
type Entities[T any] struct {
	Values []T `json:"values"` // Values is the list of wrapped entities or values
}

func (e *Entities[T]) ID() string { return "" }

func (e *Entities[T]) TABLE() string { return "" }

func (e *Entities[T]) NAME() string { return "" }

func (e *Entities[T]) KEY() string { return "" }

// NewEntities creates a new Entities instance with an empty list.
func NewEntities[T any]() Entity {
	return &Entities[T]{
		Values: make([]T, 0),
	}
}

// Add appends an item to the Entities list and returns the new count.
func (e *Entities[T]) Add(item T) int {
	e.Values = append(e.Values, item)
	return len(e.Values)
}

// endregion

// region Entity Ids ---------------------------------------------------------------------------------------------------

// ID returns a long string (alphanumeric) based on Epoch micro-seconds in base 36.
// It generates a unique identifier that is time-ordered and relatively compact.
func ID() string {
	return strconv.FormatUint(uint64(time.Now().UnixMicro()), 36)
}

// IDN returns a long string (digits only) based on Epoch micro-seconds.
// It is similar to ID() but uses only numeric characters.
func IDN() string {
	return fmt.Sprintf("%d", time.Now().UnixMicro())
}

// ShortID returns a short string (6 characters alphanumeric) based on epoch seconds in base 36.
// It accepts optional delta values to offset the timestamp, which can be useful for generating
// IDs in the past or future, or adding entropy.
func ShortID(delta ...int) string {
	value := uint64(time.Now().Unix())
	for _, d := range delta {
		value += uint64(d)
	}
	return strconv.FormatUint(value, 36)
}

// ShortIDN returns a short string (digits only) based on epoch seconds.
// It accepts optional delta values similar to ShortID.
func ShortIDN(delta ...int) string {
	value := uint64(time.Now().Unix())
	for _, d := range delta {
		value += uint64(d)
	}
	return fmt.Sprintf("%d", value)
}

// NanoID returns a long string (21 characters) based on the go-nanoid project.
// It is smaller and faster than UUID, and URL-friendly.
// If generation fails, it falls back to a time-based ID.
func NanoID() string {
	if generator, err := nanoid.Standard(21); err != nil {
		return strconv.FormatUint(uint64(time.Now().UnixMicro()), 36)
	} else {
		return generator()
	}
}

// GUID generates a new Global Unique Identifier (UUID v4).
// It uses the github.com/google/uuid library.
func GUID() string {
	return uuid.New().String()
}

// endregion

// region Entity Actions -----------------------------------------------------------------------------------------------

// EntityAction represents the type of action performed on an entity.
// @Data
type EntityAction int

const (
	// AddEntity represents an action to add a new entity.
	AddEntity EntityAction = 1
	// UpdateEntity represents an action to update an existing entity.
	UpdateEntity EntityAction = 2
	// DeleteEntity represents an action to delete an entity.
	DeleteEntity EntityAction = 3
)

// endregion

// region Clone Entity -------------------------------------------------------------------------------------------------

// CloneEntity performs a deep clone of an entity using JSON serialization.
// It creates a new instance using the provided EntityFactory and copies the data from the source entity.
//
// Parameters:
//   - ef: The EntityFactory to create the destination entity.
//   - src: The source entity to clone.
//
// Returns:
//   - A new Entity instance with the same data as src, or an error if cloning fails.
func CloneEntity(ef EntityFactory, src Entity) (Entity, error) {
	dst := ef()
	data, err := json.Marshal(src)
	if err != nil {
		return dst, err
	}
	err = json.Unmarshal(data, &dst)
	return dst, err
}

// endregion
