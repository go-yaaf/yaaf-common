package database

import (
	. "github.com/go-yaaf/yaaf-common/entity"
	"time"
)

// IQuery Database Query interface
type IQuery interface {

	// Apply adds a callback to apply on each result entity in the query
	Apply(cb func(in Entity) Entity) IQuery

	// Filter Add single field filter
	Filter(filter QueryFilter) IQuery

	// MatchAll Add list of filters, all of them should be satisfied (AND)
	MatchAll(filters ...QueryFilter) IQuery

	// MatchAny Add list of filters, any of them should be satisfied (OR)
	MatchAny(filters ...QueryFilter) IQuery

	// Sort Add sort order by field,  expects sort parameter in the following form: field_name (Ascending) or field_name- (Descending)
	Sort(sort string) IQuery

	// Page Set page number (for pagination)
	Page(page int) IQuery

	// Limit Set page size limit (for pagination)
	Limit(page int) IQuery

	// List Execute a query to get list of entities by IDs (the criteria is ignored)
	List(entityIDs []string, keys ...string) (out []Entity, err error)

	// Find Execute the query based on the criteria, order and pagination
	Find(keys ...string) (out []Entity, total int64, err error)

	// Select is similar to find but with ability to retrieve specific fields
	Select(fields ...string) ([]Json, error)

	// Count Execute the query based on the criteria, order and pagination and return only the count of matching rows
	Count(keys ...string) (total int64, err error)

	// Aggregation Execute the query based on the criteria, order and pagination and return the provided aggregation function on the field
	// supported functions: count : avg, sum, min, max
	Aggregation(field, function string, keys ...string) (value float64, err error)

	// GroupCount Execute the query based on the criteria, grouped by field and return count per group
	GroupCount(field string, keys ...string) (out map[any]int64, total int64, err error)

	// GroupAggregation Execute the query based on the criteria, order and pagination and return the aggregated value per group
	// the data point is a calculation of the provided function on the selected field, each data point includes the number of documents and the calculated value
	// the total is the sum of all calculated values in all the buckets
	// supported functions: count : avg, sum, min, max
	GroupAggregation(field, function string, keys ...string) (out map[any]Tuple[int64, float64], total float64, err error)

	// Histogram returns a time series data points based on the time field, supported intervals: Minute, Hour, Day, week, month
	// the data point is a calculation of the provided function on the selected field, each data point includes the number of documents and the calculated value
	// the total is the sum of all calculated values in all the buckets
	// supported functions: count : avg, sum, min, max
	Histogram(field, function, timeField string, interval time.Duration, keys ...string) (out map[Timestamp]Tuple[int64, float64], total float64, err error)

	// Histogram2D returns a two-dimensional time series data points based on the time field, supported intervals: Minute, Hour, Day, week, month
	// the data point is a calculation of the provided function on the selected field
	// supported functions: count : avg, sum, min, max
	Histogram2D(field, function, dim, timeField string, interval time.Duration, keys ...string) (out map[Timestamp]map[any]Tuple[int64, float64], total float64, err error)

	// FindSingle Execute query based on the where criteria to get a single (the first) result
	FindSingle(keys ...string) (entity Entity, err error)

	// GetMap Execute query based on the criteria, order and pagination and return the results as a map of id->Entity
	GetMap(keys ...string) (out map[string]Entity, err error)

	// GetIDs executes a query based on the where criteria, order and pagination and return the results as a list of Ids
	GetIDs(keys ...string) (out []string, err error)

	// Delete the entities satisfying the criteria
	Delete(keys ...string) (total int64, err error)

	// SetField Update single field of all the documents meeting the criteria in a single transaction
	SetField(field string, value any, keys ...string) (total int64, err error)

	// SetFields Update multiple fields of all the documents meeting the criteria in a single transaction
	SetFields(fields map[string]any, keys ...string) (total int64, err error)

	// ToString Get the string representation of the query
	ToString() string
}
