package database

import (
	e "github.com/go-yaaf/yaaf-common/entity"
)

// IAnalyticQuery defines the interface for building and executing analytic queries.
// It supports aggregation functions (Sum, Min, Max, Avg, Count) and grouping.
type IAnalyticQuery interface {

	// Sum calculates the sum of the specified field.
	Sum(fieldName string) IAnalyticQuery

	// Min calculates the minimum value of the specified field.
	Min(fieldName string) IAnalyticQuery

	// Max calculates the maximum value of the specified field.
	Max(fieldName string) IAnalyticQuery

	// Avg calculates the average value of the specified field.
	Avg(fieldName string) IAnalyticQuery

	// CountAll counts all occurrences of the specified field.
	CountAll(fieldName string) IAnalyticQuery

	// CountUnique counts unique occurrences of the specified field.
	CountUnique(fieldName string) IAnalyticQuery

	// GroupBy groups the results by the specified field and time period.
	GroupBy(fieldName string, timePeriod e.TimePeriodCode) IAnalyticQuery

	// Compute executes the analytic query and returns the results.
	Compute() (out []e.Entity, err error)
}

// IAdvancedQuery defines a composed interface that combines IQuery and IAnalyticQuery.
// It is intended for advanced querying capabilities, including both standard data retrieval and analytics.
type IAdvancedQuery interface {
	IQuery
	IAnalyticQuery
}
