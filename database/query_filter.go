package database

import "fmt"

// region QueryFilter Interface ----------------------------------------------------------------------------------------

// QueryFilter defines the interface for building query filters.
// It allows for constructing complex conditions for database queries.
type QueryFilter interface {

	// Eq adds an "Equal" condition.
	Eq(value any) QueryFilter

	// Neq adds a "Not Equal" condition.
	Neq(value any) QueryFilter

	// Like adds a "Like" condition (similar to SQL LIKE).
	Like(value string) QueryFilter

	// NotLike adds a "Not Like" condition.
	NotLike(value string) QueryFilter

	// Gt adds a "Greater Than" condition.
	Gt(value any) QueryFilter

	// Gte adds a "Greater Than or Equal" condition.
	Gte(value any) QueryFilter

	// Lt adds a "Less Than" condition.
	Lt(value any) QueryFilter

	// Lte adds a "Less Than or Equal" condition.
	Lte(value any) QueryFilter

	// In adds an "In" condition (match one of the values).
	In(values ...any) QueryFilter

	// NotIn adds a "Not In" condition.
	NotIn(values ...any) QueryFilter

	// InSubQuery adds a condition to match values in the result of a sub-query.
	InSubQuery(field string, subQuery IQuery) QueryFilter

	// NotInSubQuery adds a condition to exclude values in the result of a sub-query.
	NotInSubQuery(field string, subQuery IQuery) QueryFilter

	// Between adds a "Between" condition (inclusive).
	Between(value1, value2 any) QueryFilter

	// Contains adds a condition to check if an array field contains the value.
	Contains(value any) QueryFilter

	// NotContains adds a condition to check if an array field does not contain the value.
	NotContains(value any) QueryFilter

	// WithFlag adds a condition to check if an integer field has specific bit flags set.
	WithFlag(value int) QueryFilter

	// WithNoFlag adds a condition to check if an integer field does not have specific bit flags set.
	WithNoFlag(value int) QueryFilter

	// IsEmpty adds a condition to check if a field is null or empty.
	IsEmpty() QueryFilter

	// IsTrue adds a condition to check if a boolean field is true.
	IsTrue() QueryFilter

	// IsFalse adds a condition to check if a boolean field is false or null.
	IsFalse() QueryFilter

	// If enables or disables the filter based on the boolean condition.
	If(value bool) QueryFilter

	// IsActive returns true if the filter is active.
	IsActive() bool

	// GetField returns the field name being filtered.
	GetField() string

	// GetOperator returns the filter operator.
	GetOperator() QueryOperator

	// GetValues returns the filter values.
	GetValues() []any

	// GetStringValue returns the string representation of the value at the given index.
	GetStringValue(index int) string

	// GetSubQuery returns the underlying sub-query.
	GetSubQuery() IQuery

	// GetSubQueryField returns the field used in the sub-query.
	GetSubQueryField() string
}

// endregion

// region QueryFilter internal implementation --------------------------------------------------------------------------

// queryFilter implements the QueryFilter interface.
type queryFilter struct {
	field         string
	operator      QueryOperator
	values        []any
	active        bool
	subQuery      IQuery
	subQueryField string
}

// Filter creates a new QueryFilter for the specified field.
func Filter(field string) QueryFilter {
	return &queryFilter{
		field:    field,
		operator: Eq,
		active:   true,
	}
}

// F is a shorthand alias for Filter.
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

// NotLike - not similar
func (q *queryFilter) NotLike(value string) QueryFilter {
	q.operator = NotLike
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

// NotContains - a field of type array does not contain the provided value
func (q *queryFilter) NotContains(value any) QueryFilter {
	q.operator = NotContains
	q.values = append(q.values, value)
	return q
}

// WithFlag - a field of type integer representing bit flags include a flag or set of flags
func (q *queryFilter) WithFlag(value int) QueryFilter {
	q.operator = WithFlag
	q.values = append(q.values, value)
	return q
}

// WithNoFlag - a field of type integer representing bit flags does not include a flag or set of flags
func (q *queryFilter) WithNoFlag(value int) QueryFilter {
	q.operator = WithNoFlag
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

// IsTrue - boolean field is true
func (q *queryFilter) IsTrue() QueryFilter {
	q.operator = True
	q.active = true
	return q
}

// IsFalse - boolean field is null or false
func (q *queryFilter) IsFalse() QueryFilter {
	q.operator = False
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
