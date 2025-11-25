package database

import (
	"time"

	. "github.com/go-yaaf/yaaf-common/entity"
)

// AggFunc represents an aggregation function type (e.g., count, sum, avg).
type AggFunc string

const (
	COUNT AggFunc = "count" // Count aggregation
	SUM   AggFunc = "sum"   // Sum aggregation
	AVG   AggFunc = "avg"   // Average aggregation
	MIN   AggFunc = "min"   // Minimum aggregation
	MAX   AggFunc = "max"   // Maximum aggregation
)

// IQuery defines the interface for building and executing database queries.
// It supports filtering, sorting, pagination, and aggregation.
type IQuery interface {

	// Apply adds a callback function to be applied to each result entity.
	Apply(cb func(in Entity) Entity) IQuery

	// Filter adds a single field filter to the query.
	Filter(filter QueryFilter) IQuery

	// Range adds a time frame filter on a specific time field.
	Range(field string, from Timestamp, to Timestamp) IQuery

	// MatchAll adds a list of filters, all of which must be satisfied (AND logic).
	MatchAll(filters ...QueryFilter) IQuery

	// MatchAny adds a list of filters, any of which must be satisfied (OR logic).
	MatchAny(filters ...QueryFilter) IQuery

	// Sort adds a sort order by field.
	// The sort parameter should be in the format: "field_name" (Ascending) or "field_name-" (Descending).
	Sort(sort string) IQuery

	// Page sets the requested page number for pagination (0-based).
	Page(page int) IQuery

	// Limit sets the page size limit for pagination.
	Limit(page int) IQuery

	// List executes the query to retrieve a list of entities by their IDs, ignoring other criteria.
	// The 'keys' argument is optional and can be used for sharding or other specific lookup mechanisms.
	List(entityIDs []string, keys ...string) (out []Entity, err error)

	// Find executes the query based on criteria, order, and pagination.
	// The 'keys' argument is optional and can be used for sharding or other specific lookup mechanisms.
	Find(keys ...string) (out []Entity, total int64, err error)

	// Select executes the query and returns specific fields as a list of Json maps.
	Select(fields ...string) ([]Json, error)

	// Count executes the query and returns the number of matching entities.
	// The 'keys' argument is optional and can be used for sharding or other specific lookup mechanisms.
	Count(keys ...string) (total int64, err error)

	// Aggregation executes an aggregation function on a field for the matching entities.
	// Supported functions: count, sum, avg, min, max.
	Aggregation(field string, function AggFunc, keys ...string) (value float64, err error)

	// GroupCount executes the query and returns the count of entities per group (grouped by field).
	GroupCount(field string, keys ...string) (out map[any]int64, total int64, err error)

	// GroupAggregation executes the query and returns the aggregated value per group.
	// Each data point includes the count of documents and the calculated value.
	GroupAggregation(field string, function AggFunc, keys ...string) (out map[any]Tuple[int64, float64], total float64, err error)

	// Histogram returns time series data points based on a time field.
	// Supported intervals: Minute, Hour, Day, Week, Month.
	Histogram(field string, function AggFunc, timeField string, interval time.Duration, keys ...string) (out map[Timestamp]Tuple[int64, float64], total float64, err error)

	// Histogram2D returns two-dimensional time series data points based on a time field.
	Histogram2D(field string, function AggFunc, dim, timeField string, interval time.Duration, keys ...string) (out map[Timestamp]map[any]Tuple[int64, float64], total float64, err error)

	// FindSingle executes the query and returns the first matching entity.
	// The 'keys' argument is optional and can be used for sharding or other specific lookup mechanisms.
	FindSingle(keys ...string) (entity Entity, err error)

	// GetMap executes the query and returns the results as a map of ID -> Entity.
	// The 'keys' argument is optional and can be used for sharding or other specific lookup mechanisms.
	GetMap(keys ...string) (out map[string]Entity, err error)

	// GetIDs executes the query and returns a list of IDs of the matching entities.
	// The 'keys' argument is optional and can be used for sharding or other specific lookup mechanisms.
	GetIDs(keys ...string) (out []string, err error)

	// Delete removes the entities matching the query criteria.
	// The 'keys' argument is optional and can be used for sharding or other specific lookup mechanisms.
	Delete(keys ...string) (total int64, err error)

	// SetField updates a single field for all documents matching the criteria.
	// The 'keys' argument is optional and can be used for sharding or other specific lookup mechanisms.
	SetField(field string, value any, keys ...string) (total int64, err error)

	// SetFields updates multiple fields for all documents matching the criteria.
	// The 'keys' argument is optional and can be used for sharding or other specific lookup mechanisms.
	SetFields(fields map[string]any, keys ...string) (total int64, err error)

	// ToString returns a string representation of the query.
	ToString() string
}
