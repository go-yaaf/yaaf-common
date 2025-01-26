package database

import (
	e "github.com/go-yaaf/yaaf-common/entity"
)

type IAnalyticQuery interface {
	Sum(fieldName string) IAnalyticQuery
	Min(fieldName string) IAnalyticQuery
	Max(fieldName string) IAnalyticQuery
	Avg(fieldName string) IAnalyticQuery
	CountAll(fieldName string) IAnalyticQuery
	CountUnique(fieldName string) IAnalyticQuery
	GroupBy(fieldName string, timePeriod e.TimePeriodCode) IAnalyticQuery
	Compute() (out []e.Entity, err error)
}

type IAdvancedQuery interface {
	IQuery
	IAnalyticQuery
}
