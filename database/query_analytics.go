package database

import (
	e "github.com/go-yaaf/yaaf-common/entity"
)

type IQueryAnalytic interface {
	Sum(fieldName string) IQueryAnalytic
	Min(fieldName string) IQueryAnalytic
	Max(fieldName string) IQueryAnalytic
	Avg(fieldName string) IQueryAnalytic
	CountAll(fieldName string) IQueryAnalytic
	CountUnique(fieldName string) IQueryAnalytic
	GroupBy(fieldName string, timePeriod e.TimePeriodCode) IQueryAnalytic
	Compute() (out []e.Entity, err error)
}

type IAdvancedQuery interface {
	IQuery
	IQueryAnalytic
}
