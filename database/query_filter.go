package database

import "fmt"

// region QueryFilter Interface ----------------------------------------------------------------------------------------

// QueryFilter Query filter interface
type QueryFilter interface {

	// Eq - Equal
	Eq(value any) QueryFilter

	// Neq - Not equal
	Neq(value any) QueryFilter

	// Like - similar
	Like(value string) QueryFilter

	// Gt - Greater than
	Gt(value any) QueryFilter

	// Gte - Greater or equal
	Gte(value any) QueryFilter

	// Lt - Less than
	Lt(value any) QueryFilter

	// Lte - Less or equal
	Lte(value any) QueryFilter

	// In - mach one of the values
	In(values ...any) QueryFilter

	// NotIn - Not In
	NotIn(values ...any) QueryFilter

	// Between - equal or greater than the lower boundary and equal or less than the upper boundary
	Between(value1, value2 any) QueryFilter

	// If - Include this filter only if condition is true
	If(value bool) QueryFilter

	// IsActive Include this filter only if condition is true
	IsActive() bool

	// GetField Get the field name
	GetField() string

	// GetOperator Get the criteria operator
	GetOperator() QueryOperator

	// GetValues Get the criteria values
	GetValues() []any

	// GetStringValue Get string representation of the value
	GetStringValue(index int) string
}

// endregion

// region QueryFilter internal implementation --------------------------------------------------------------------------

// Query filter
type queryFilter struct {
	field    string
	operator QueryOperator
	values   []any
	active   bool
}

// Filter by field
func Filter(field string) QueryFilter {
	return &queryFilter{
		field:    field,
		operator: Eq,
		active:   true,
	}
}

// F filter by field (synonym to Filter)
func F(field string) QueryFilter {
	return &queryFilter{
		field:    field,
		operator: Eq,
		active:   true,
	}
}

// Eq - Equal
func (q *queryFilter) Eq(value any) QueryFilter {
	q.operator = Eq
	q.values = append(q.values, value)
	//q.values = []any{value}
	q.active = len(fmt.Sprintf("%v", value)) > 0
	return q
}

// Neq - Not equal
func (q *queryFilter) Neq(value any) QueryFilter {
	q.operator = Neq
	q.values = append(q.values, value)
	//q.values = []any{value}
	q.active = len(fmt.Sprintf("%v", value)) > 0
	return q
}

// Like - similar
func (q *queryFilter) Like(value string) QueryFilter {
	q.operator = Like
	q.values = append(q.values, value)
	//q.values = []any{value}
	q.active = len(fmt.Sprintf("%v", value)) > 0
	return q
}

// Gt - Greater than
func (q *queryFilter) Gt(value any) QueryFilter {
	q.operator = Gt
	q.values = append(q.values, value)
	//q.values = []any{value}
	return q
}

// Gte - Greater or equal
func (q *queryFilter) Gte(value any) QueryFilter {
	q.operator = Gte
	q.values = append(q.values, value)
	//q.values = []any{value}
	return q
}

// Lt - Less than
func (q *queryFilter) Lt(value any) QueryFilter {
	q.operator = Lt
	q.values = append(q.values, value)
	//q.values = []any{value}
	return q
}

// Lte - Less or equal
func (q *queryFilter) Lte(value any) QueryFilter {
	q.operator = Lte
	q.values = append(q.values, value)
	//q.values = []any{value}
	return q
}

// In - mach one of the values
func (q *queryFilter) In(values ...any) QueryFilter {
	q.operator = In
	q.values = append(q.values, values...)
	//q.values = []any{values}
	q.active = len(values) > 0
	return q
}

// NotIn - Not in
func (q *queryFilter) NotIn(values ...any) QueryFilter {
	q.operator = NotIn
	q.values = append(q.values, values...)
	//q.values = []any{values}
	q.active = len(values) > 0
	return q
}

// Between - equal or greater than the lower boundary and equal or less than the upper boundary
func (q *queryFilter) Between(value1, value2 any) QueryFilter {
	q.operator = Between
	q.values = append(q.values, value1, value2)
	//q.values = []any{value1, value2}
	return q
}

// If - Include this filter only if condition is true
func (q *queryFilter) If(value bool) QueryFilter {
	q.active = value
	return q
}

// IsActive Is the filter active?
func (q *queryFilter) IsActive() bool {
	return q.active
}

// GetField Get filtered field name
func (q *queryFilter) GetField() string {
	return q.field
}

// GetOperator Get the criteria operator
func (q *queryFilter) GetOperator() QueryOperator {
	return q.operator
}

// GetValues Get values
func (q *queryFilter) GetValues() []any {
	return q.values
}

// GetStringValue Get string representation of the value
func (q *queryFilter) GetStringValue(index int) string {
	if len(q.values) > index {
		return fmt.Sprintf("%v", q.values[index])
	} else {
		return ""
	}
}

// endregion
