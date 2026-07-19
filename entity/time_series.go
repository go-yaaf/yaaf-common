package entity

// region TimeSeries ---------------------------------------------------------------------------------------------------

// TimeSeries represents a named collection of data points over a specific time range.
// @Data
type TimeSeries[T any] struct {
	Name   string             `json:"name"`   // Name of the time series
	Range  TimeFrame          `json:"range"`  // Range covers the start and end time of the series
	Values []TimeDataPoint[T] `json:"values"` // Values contains the data points
}

// ID returns the time series name as its ID.
func (ts *TimeSeries[T]) ID() string { return ts.Name }

// TABLE returns an empty string as TimeSeries is typically not a database table itself.
func (ts *TimeSeries[T]) TABLE() string { return "" }

// NAME returns the time series name.
func (ts *TimeSeries[T]) NAME() string { return ts.Name }

// KEY returns an empty string.
func (ts *TimeSeries[T]) KEY() string { return "" }

// SetDataPoint updates the value of a data point at a specific timestamp.
// It returns true if the data point was found and updated, false otherwise.
func (ts *TimeSeries[T]) SetDataPoint(t Timestamp, val T) bool {
	if len(ts.Values) == 0 {
		return false
	}
	for i := range ts.Values {
		if ts.Values[i].Timestamp == t {
			ts.Values[i].Value = val
			return true
		}
	}
	return false
}

// endregion
