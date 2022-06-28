// Copyright 2022. Motty Cohen
//
// Database Query filter
//
package database

// region QueryFilter Interface ----------------------------------------------------------------------------------------

// Query filter interface
type QueryFilter interface {

	// Equal
	Eq(value any) QueryFilter

	// Not equal
	Neq(value any) QueryFilter

	// Like
	Like(value string) QueryFilter

	// Greater
	Gt(value any) QueryFilter

	// Greater or equal
	Gte(value any) QueryFilter

	// Less than
	Lt(value any) QueryFilter

	// Less or equal
	Lte(value any) QueryFilter

	// In
	In(values ...any) QueryFilter

	// Not In
	NotIn(values ...any) QueryFilter

	// Between
	Between(value1, value2 any) QueryFilter

	// Include this filter only if condition is true
	If(value bool) QueryFilter

	// Include this filter only if condition is true
	IsActive() bool

	Field() string

	Values() []any
}

// endregion

// region QueryFilter internal implementation --------------------------------------------------------------------------

// Query filter
type queryFilter struct {
	field     string
	operator  queryOperator
	values    []any
	condition bool
}

// Filter by field
func F(field string) QueryFilter {
	return &queryFilter{
		field:    field,
		operator: Eq,
	}
}

// Equal
func (q *queryFilter) Eq(value any) QueryFilter {
	q.operator = Eq
	q.values = []any{value}
	return q
}

// Not equal
func (q *queryFilter) Neq(value any) QueryFilter {
	q.operator = Neq
	q.values = []any{value}
	return q
}

// Like
func (q *queryFilter) Like(value string) QueryFilter {
	q.operator = Like
	q.values = []any{value}
	return q
}

// Greater
func (q *queryFilter) Gt(value any) QueryFilter {
	q.operator = Gt
	q.values = []any{value}
	return q
}

// Greater or equal
func (q *queryFilter) Gte(value any) QueryFilter {
	q.operator = Gte
	q.values = []any{value}
	return q
}

// Less
func (q *queryFilter) Lt(value any) QueryFilter {
	q.operator = Lt
	q.values = []any{value}
	return q
}

// Less or equal
func (q *queryFilter) Lte(value any) QueryFilter {
	q.operator = Lte
	q.values = []any{value}
	return q
}

// In
func (q *queryFilter) In(values ...any) QueryFilter {
	q.operator = In
	q.values = []any{values}
	return q
}

// Not in
func (q *queryFilter) NotIn(values ...any) QueryFilter {
	q.operator = NotIn
	q.values = []any{values}
	return q
}

// Between
func (q *queryFilter) Between(value1, value2 any) QueryFilter {
	q.operator = Between
	q.values = []any{value1, value2}
	return q
}

// Include this filter only if condition is true
func (q *queryFilter) If(value bool) QueryFilter {
	q.condition = value
	return q
}

// Is the filter active?
func (q *queryFilter) IsActive() bool {
	return q.condition
}

// Get filtered field name
func (q *queryFilter) Field() string {
	return q.field
}

// Get values
func (q *queryFilter) Values() []any {
	return q.values
}

// endregion
