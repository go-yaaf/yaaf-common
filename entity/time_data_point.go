package entity

import (
	"fmt"
)

// region TimeDataPoint ------------------------------------------------------------------------------------------------

// TimeDataPoint represents a generic data point associated with a timestamp.
// @Data
type TimeDataPoint[V any] struct {
	Timestamp Timestamp `json:"timestamp"` // Timestamp of the data point
	Value     V         `json:"value"`     // Value of the data point
}

// NewTimeDataPoint creates a new TimeDataPoint instance.
func NewTimeDataPoint[V any](ts Timestamp, value V) TimeDataPoint[V] {
	return TimeDataPoint[V]{Timestamp: ts, Value: value}
}

// String returns a string representation of the TimeDataPoint.
func (tf *TimeDataPoint[V]) String(format string) string {
	return fmt.Sprintf("%s - %v", tf.Timestamp.String(format), tf.Value)
}

// endregion
