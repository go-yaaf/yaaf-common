package database

import (
	"time"

	e "github.com/go-yaaf/yaaf-common/entity"
)

type IQueryAnalytic interface {
	Sum(fieldName string) IQueryAnalytic
	Min(fieldName string) IQueryAnalytic
	Max(fieldName string) IQueryAnalytic
	Avg(fieldName string) IQueryAnalytic
	GroupBy(fieldName string) IQueryAnalytic
	GroupByTimePeriod(fieldName string, period time.Duration) IQueryAnalytic
	Compute() (out []e.Entity, err error)
}

type IAdvancedQuery interface {
	IQuery
	IQueryAnalytic
}
