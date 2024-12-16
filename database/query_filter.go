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

	// In - match one of the values
	In(values ...any) QueryFilter

	// NotIn - Not In
	NotIn(values ...any) QueryFilter

	// InSubQuery - match one of the values in the result of the sub query
	InSubQuery(field string, subQuery IQuery) QueryFilter

	// NotInSubQuery - exclude any record matching one of the values in the result of the sub query
	NotInSubQuery(field string, subQuery IQuery) QueryFilter

	// Between - equal or greater than the lower boundary and equal or less than the upper boundary
	Between(value1, value2 any) QueryFilter

	// Contains - a field of type array contains the provided value
	Contains(value any) QueryFilter

	// IsEmpty - field is null or empty
	IsEmpty() QueryFilter

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

	// GetSubQuery gets the underlying sub-query
	GetSubQuery() IQuery

	// GetSubQueryField gets the underlying sub-query field
	GetSubQueryField() string
}

// endregion

// region QueryFilter internal implementation --------------------------------------------------------------------------

// Query filter
type queryFilter struct {
	field         string
	operator      QueryOperator
	values        []any
	active        bool
	subQuery      IQuery
	subQueryField string
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
	q.active = len(fmt.Sprintf("%v", value)) > 0
	return q
}

// Neq - Not equal
func (q *queryFilter) Neq(value any) QueryFilter {
	q.operator = Neq
	q.values = append(q.values, value)
	q.active = len(fmt.Sprintf("%v", value)) > 0
	return q
}

// Like - similar
func (q *queryFilter) Like(value string) QueryFilter {
	q.operator = Like
	q.values = append(q.values, value)
	q.active = len(fmt.Sprintf("%v", value)) > 0
	return q
}

// Gt - Greater than
func (q *queryFilter) Gt(value any) QueryFilter {
	q.operator = Gt
	q.values = append(q.values, value)
	return q
}

// Gte - Greater or equal
func (q *queryFilter) Gte(value any) QueryFilter {
	q.operator = Gte
	q.values = append(q.values, value)
	return q
}

// Lt - Less than
func (q *queryFilter) Lt(value any) QueryFilter {
	q.operator = Lt
	q.values = append(q.values, value)
	return q
}

// Lte - Less or equal
func (q *queryFilter) Lte(value any) QueryFilter {
	q.operator = Lte
	q.values = append(q.values, value)
	return q
}

// In - mach one of the values
func (q *queryFilter) In(values ...any) QueryFilter {
	q.operator = In
	q.values = append(q.values, values...)
	q.active = len(values) > 0
	return q
}

// NotIn - Not in
func (q *queryFilter) NotIn(values ...any) QueryFilter {
	q.operator = NotIn
	q.values = append(q.values, values...)
	q.active = len(values) > 0
	return q
}

// InSubQuery - match one of the values in the result of the sub query
func (q *queryFilter) InSubQuery(field string, subQuery IQuery) QueryFilter {
	q.operator = InSQ
	q.subQuery = subQuery
	q.subQueryField = field
	return q
}

// NotInSubQuery - exclude any record matching one of the values in the result of the sub query
func (q *queryFilter) NotInSubQuery(field string, subQuery IQuery) QueryFilter {
	q.operator = NotInSQ
	q.subQuery = subQuery
	q.subQueryField = field
	return q
}

// Between - equal or greater than the lower boundary and equal or less than the upper boundary
func (q *queryFilter) Between(value1, value2 any) QueryFilter {
	q.operator = Between
	q.values = append(q.values, value1, value2)
	return q
}

// Contains - a field of type array contains the provided value
func (q *queryFilter) Contains(value any) QueryFilter {
	q.operator = Contains
	q.values = append(q.values, value)
	return q
}

// If - Include this filter only if condition is true
func (q *queryFilter) If(value bool) QueryFilter {
	q.active = value
	return q
}

// IsEmpty - Field does not include value (null or empty)
func (q *queryFilter) IsEmpty() QueryFilter {
	q.operator = Empty
	q.active = true
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

// GetSubQuery gets the underlying sub-query
func (q *queryFilter) GetSubQuery() IQuery {
	return q.subQuery
}

// GetSubQueryField gets the underlying sub-query field
func (q *queryFilter) GetSubQueryField() string {
	return q.subQueryField
}

// endregion
