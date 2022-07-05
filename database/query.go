// Copyright 2022. Motty Cohen
//
// Database Query interface
//
package database

import (
	. "github.com/mottyc/yaaf-common/entity"
)

// Database Query interface
type IQuery interface {

	// Add callback to apply on each result entity in the query
	Apply(cb func(in Entity) Entity) IQuery

	// Add single field filter
	Filter(filter QueryFilter) IQuery

	// Add list of filters, all of them should be satisfied (AND)
	MatchAll(filters ...QueryFilter) IQuery

	// Add list of filters, any of them should be satisfied (OR)
	MatchAny(filters ...QueryFilter) IQuery

	// Add sort order by field,  expects sort parameter in the following form: field_name (Ascending) or field_name- (Descending)
	Sort(sort string) IQuery

	// Set page number (for pagination)
	Page(page int) IQuery

	// Set page size limit (for pagination)
	Limit(page int) IQuery

	// Execute a query to get list of entities by IDs (the criteria is ignored)
	List(entityIDs []string, keys ...string) (out []Entity, err error)

	// Execute the query based on the criteria, order and pagination
	Find(keys ...string) (out []Entity, total int64, err error)

	// Execute query based on the where criteria to get a single (the first) result
	FindSingle(keys ...string) (entity Entity, err error)

	// Execute query based on the criteria, order and pagination and return the results as a map of id->Entity
	GetMap(keys ...string) (out map[string]Entity, err error)

	// Execute query based on the where criteria, order and pagination and return the results as a list of Ids
	GetIds(keys ...string) (out []string, err error)

	// Delete the entities satisfying the criteria
	Delete(keys ...string) (total int64, err error)

	// Update single field of all the documents meeting the criteria in a single transaction
	SetField(field string, value any, keys ...string) (total int64, err error)

	// Update multiple fields of all the documents meeting the criteria in a single transaction
	SetFields(fields map[string]any, keys ...string) (total int64, err error)

	// Get the string representation of the query
	ToString() string
}
