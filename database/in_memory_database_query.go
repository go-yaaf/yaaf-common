package database

import (
	"fmt"
	"strings"
	"time"

	. "github.com/go-yaaf/yaaf-common/entity"
	"github.com/go-yaaf/yaaf-common/utils"
)

// region queryBuilder internal structure ------------------------------------------------------------------------------

type inMemoryDatabaseQuery struct {
	db         *InMemoryDatabase        // A reference to the underlying IDatabase
	factory    EntityFactory            // The entity factory method
	allFilters [][]QueryFilter          // List of lists of AND filters
	anyFilters [][]QueryFilter          // List of lists of OR filters
	ascOrders  []any                    // List of fields for ASC order
	descOrders []any                    // List of fields for DESC order
	callbacks  []func(in Entity) Entity // List of entity transformation callback functions
	page       int                      // Page number (for pagination)
	limit      int                      // Page size: how many results in a page (for pagination)
	rangeField string                   // Field name for range filter (must be timestamp field)
	rangeFrom  Timestamp                // Start timestamp for range filter
	rangeTo    Timestamp                // End timestamp for range filter
}

// endregion

// region QueryBuilder Construction Methods ----------------------------------------------------------------------------

// Apply adds callback to apply on each result entity in the query
func (s *inMemoryDatabaseQuery) Apply(cb func(in Entity) Entity) IQuery {
	if cb != nil {
		s.callbacks = append(s.callbacks, cb)
	}
	return s
}

// Filter adds a single field filter
func (s *inMemoryDatabaseQuery) Filter(filter QueryFilter) IQuery {
	if filter.IsActive() {
		s.allFilters = append(s.allFilters, []QueryFilter{filter})
	}
	return s
}

// Range add time frame filter on specific time field
func (s *inMemoryDatabaseQuery) Range(field string, from Timestamp, to Timestamp) IQuery {
	s.rangeField = field
	s.rangeFrom = from
	s.rangeTo = to
	return s
}

// MatchAll adds a list of filters, all of them should be satisfied (AND operator equivalent)
func (s *inMemoryDatabaseQuery) MatchAll(filters ...QueryFilter) IQuery {
	list := make([]QueryFilter, 0)
	for _, filter := range filters {
		if filter.IsActive() {
			list = append(list, filter)
		}
	}
	s.allFilters = append(s.allFilters, list)
	return s
}

// MatchAny adds a list of filters, any of them should be satisfied (OR operator equivalent)
func (s *inMemoryDatabaseQuery) MatchAny(filters ...QueryFilter) IQuery {
	list := make([]QueryFilter, 0)
	for _, filter := range filters {
		if filter.IsActive() == true {
			list = append(list, filter)
		}
	}
	s.anyFilters = append(s.allFilters, list)
	return s
}

// Sort adds sort order by field
// The expects sort parameter should be in the following form: field_name (Ascending) or field_name- (Descending)
func (s *inMemoryDatabaseQuery) Sort(sort string) IQuery {
	if sort == "" {
		return s
	}

	// as a default, order will be ASC
	if strings.HasSuffix(sort, "-") {
		s.descOrders = append(s.descOrders, sort[0:len(sort)-1])
	} else if strings.HasSuffix(sort, "+") {
		s.ascOrders = append(s.ascOrders, sort[0:len(sort)-1])
	} else {
		s.ascOrders = append(s.ascOrders, sort)
	}
	return s
}

// Limit sets the page size limit (for pagination)
func (s *inMemoryDatabaseQuery) Limit(limit int) IQuery {
	s.limit = limit
	return s
}

// Page sets the requested page number (used for pagination)
func (s *inMemoryDatabaseQuery) Page(page int) IQuery {
	s.page = page
	return s
}

// endregion

// region QueryBuilder Execution Methods -------------------------------------------------------------------------------

// List executes a query to get a list of entities by IDs (the criteria is ignored)
func (s *inMemoryDatabaseQuery) List(entityIDs []string, keys ...string) (out []Entity, err error) {
	result, err := s.db.List(s.factory, entityIDs, keys...)
	if err != nil {
		return nil, err
	}

	// Apply filters
	for _, entity := range result {
		transformed := s.processCallbacks(entity)
		if transformed != nil {
			out = append(out, transformed)
		}
	}
	return
}

// Find executes a query based on the criteria, order and pagination
// On each record, after the marshaling the result shall be transformed via the query callback chain
func (s *inMemoryDatabaseQuery) Find(keys ...string) (out []Entity, total int64, err error) {
	ent := s.factory()
	table := tableName(ent.TABLE(), keys...)

	tbl, ok := s.db.db[table]
	if !ok {
		return nil, 0, fmt.Errorf(TABLE_NOT_EXISTS)
	}

	// If range is defined, add it to the filters
	if len(s.rangeField) > 0 {
		rangeFilter := []QueryFilter{F(s.rangeField).Between(s.rangeFrom, s.rangeTo)}
		s.allFilters = append(s.allFilters, rangeFilter)
	}

	// Apply filters
	for _, entity := range tbl.Table() {
		ent := s.filter(entity)
		if ent == nil {
			continue
		}

		// apply callbacks
		transformed := s.processCallbacks(entity)
		if transformed != nil {
			out = append(out, transformed)
		}
	}

	return out, int64(len(out)), nil
}

// Select is similar to find but with ability to retrieve specific fields
func (s *inMemoryDatabaseQuery) Select(fields ...string) ([]Json, error) {
	return nil, fmt.Errorf(NOT_IMPLEMENTED)
}

// Count executes a query based on the criteria, order and pagination
// Returns only the count of matching rows
func (s *inMemoryDatabaseQuery) Count(keys ...string) (total int64, err error) {
	ent := s.factory()
	table := tableName(ent.TABLE(), keys...)

	tbl, ok := s.db.db[table]
	if !ok {
		return 0, fmt.Errorf(TABLE_NOT_EXISTS)
	}

	// If range is defined, add it to the filters
	if len(s.rangeField) > 0 {
		rangeFilter := []QueryFilter{F(s.rangeField).Between(s.rangeFrom, s.rangeTo)}
		s.allFilters = append(s.allFilters, rangeFilter)
	}

	// Apply filters
	for _, entity := range tbl.Table() {
		ent := s.filter(entity)
		if ent == nil {
			continue
		}

		// apply callbacks
		transformed := s.processCallbacks(entity)
		if transformed != nil {
			total += 1
		}
	}

	return total, nil
}

// Aggregation Execute the query based on the criteria, order and pagination and return the provided aggregation function on the field
// supported functions: count : avg, sum, min, max
func (s *inMemoryDatabaseQuery) Aggregation(field string, function AggFunc, keys ...string) (value float64, err error) {
	return 0, fmt.Errorf("not yet implemented")
}

// GroupCount Execute the query based on the criteria, grouped by field and return count per group
func (s *inMemoryDatabaseQuery) GroupCount(field string, keys ...string) (out map[any]int64, total int64, err error) {
	return nil, 0, fmt.Errorf("not yet implemented")
}

// GroupAggregation Execute the query based on the criteria, order and pagination and return the aggregated value per group
// the data point is a calculation of the provided function on the selected field, each data point includes the number of documents and the calculated value
// the total is the sum of all calculated values in all the buckets
// supported functions: count : avg, sum, min, max
func (s *inMemoryDatabaseQuery) GroupAggregation(field string, function AggFunc, keys ...string) (out map[any]Tuple[int64, float64], total float64, err error) {
	return nil, 0, fmt.Errorf("not yet implemented")
}

// Histogram returns a time series data points based on the time field, supported intervals: Minute, Hour, Day, week, month
// the data point is a calculation of the provided function on the selected field, each data point includes the number of documents and the calculated value
// the total is the sum of all calculated values in all the buckets
// supported functions: count : avg, sum, min, max
func (s *inMemoryDatabaseQuery) Histogram(field string, function AggFunc, timeField string, interval time.Duration, keys ...string) (out map[Timestamp]Tuple[int64, float64], total float64, err error) {
	return nil, 0, fmt.Errorf("not yet implemented")
}

// Histogram2D returns a two-dimensional time series data points based on the time field, supported intervals: Minute, Hour, Day, week, month
// the data point is a calculation of the provided function on the selected field
// supported functions: count : avg, sum, min, max
func (s *inMemoryDatabaseQuery) Histogram2D(field string, function AggFunc, dim, timeField string, interval time.Duration, keys ...string) (out map[Timestamp]map[any]Tuple[int64, float64], total float64, err error) {
	return nil, 0, fmt.Errorf("not yet implemented")
}

// FindSingle execute a query based on the where criteria to get a single (the first) result
// After the marshaling the result shall be transformed via the query callback chain
func (s *inMemoryDatabaseQuery) FindSingle(keys ...string) (entity Entity, err error) {
	if list, _, fe := s.Find(keys...); fe != nil {
		return nil, fe
	} else {
		if len(list) == 0 {
			return nil, fmt.Errorf("not found")
		} else {
			return list[0], nil
		}
	}
}

// GetMap execute a query based on the criteria, order and pagination and return the results as a map of id->Entity
func (s *inMemoryDatabaseQuery) GetMap(keys ...string) (out map[string]Entity, err error) {
	out = make(map[string]Entity)
	if list, _, fe := s.Find(keys...); fe != nil {
		return nil, fe
	} else {
		for _, ent := range list {
			out[ent.ID()] = ent
		}
	}
	return
}

// GetIDs executes a query based on the where criteria, order and pagination and return the results as a list of Ids
func (s *inMemoryDatabaseQuery) GetIDs(keys ...string) (out []string, err error) {
	out = make([]string, 0)

	if list, _, fe := s.Find(keys...); fe != nil {
		return nil, fe
	} else {
		for _, ent := range list {
			out = append(out, ent.ID())
		}
	}
	return
}

// Delete executes a delete command based on the where criteria
func (s *inMemoryDatabaseQuery) Delete(keys ...string) (total int64, err error) {
	deleteIds := make([]string, 0)

	if list, _, fe := s.Find(keys...); fe != nil {
		return 0, fe
	} else {
		for _, ent := range list {
			deleteIds = append(deleteIds, ent.ID())
		}
	}

	if affected, fe := s.db.BulkDelete(s.factory, deleteIds, keys...); fe != nil {
		return 0, fe
	} else {
		return affected, nil
	}
}

// SetField updates a single field of all the documents meeting the criteria in a single transaction
func (s *inMemoryDatabaseQuery) SetField(field string, value any, keys ...string) (total int64, err error) {
	fields := make(map[string]any)
	fields[field] = value
	return s.SetFields(fields, keys...)
}

// SetFields updates multiple fields of all the documents meeting the criteria in a single transaction
func (s *inMemoryDatabaseQuery) SetFields(fields map[string]any, keys ...string) (total int64, err error) {
	changeList := make([]Entity, 0)

	list, _, fe := s.Find(keys...)
	if fe != nil {
		return 0, fe
	}

	for _, entity := range list {

		raw, er := utils.JsonUtils().ToJson(entity)
		if er != nil {
			continue
		}

		for f, v := range fields {
			raw[f] = v
		}

		if changed, _ := utils.JsonUtils().FromJson(s.factory, raw); changed != nil {
			changeList = append(changeList, changed)
		}
	}

	if result, err := s.db.BulkUpdate(changeList); fe != nil {
		return 0, err
	} else {
		return result, nil
	}
}

// endregion

// region QueryBuilder Internal Methods --------------------------------------------------------------------------------
// Filter entity based on conditions
func (s *inMemoryDatabaseQuery) filter(in Entity) (out Entity) {

	// convert entity to Json
	raw, fe := utils.JsonUtils().ToJson(in)
	if fe != nil {
		return in
	}

	// Apply All (AND) filters
	for _, andList := range s.allFilters {
		for _, fq := range andList {
			if testField(raw, fq) == false {
				return nil
			}
		}
	}

	or := func(list []QueryFilter) bool {
		for _, fq := range list {
			if testField(raw, fq) == true {
				return true
			}
		}
		return false
	}

	// Apply Any (OR) filters
	for _, orList := range s.anyFilters {
		if or(orList) == false {
			return nil
		}
	}

	return in
}

// processCallbacks transforms the entity through the chain of callbacks
func (s *inMemoryDatabaseQuery) processCallbacks(in Entity) (out Entity) {
	if len(s.callbacks) == 0 {
		out = in
		return
	}

	tmp := in
	for _, cb := range s.callbacks {
		out = cb(tmp)
		if out == nil {
			return nil
		} else {
			tmp = out
		}
	}
	return
}

// endregion

// region QueryBuilder ToString Methods --------------------------------------------------------------------------------

// ToString gets a string representation of the query
func (s *inMemoryDatabaseQuery) ToString() string {
	// Create Json representing the internal builder
	if bytes, err := Marshal(s); err != nil {
		return err.Error()
	} else {
		return string(bytes)
	}
}

// endregion
