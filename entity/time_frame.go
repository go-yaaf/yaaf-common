package entity

import (
	"fmt"
	"time"
)

// region TimeFrame ----------------------------------------------------------------------------------------------------

// TimeFrame represents a time interval with a start and end timestamp.
// @Data
type TimeFrame struct {
	From Timestamp `json:"from"` // From is the start timestamp
	To   Timestamp `json:"to"`   // To is the end timestamp
}

// NewTimeFrame creates a new TimeFrame from start and end timestamps.
func NewTimeFrame(from, to Timestamp) TimeFrame {
	return TimeFrame{From: from, To: to}
}

// GetTimeFrame creates a new TimeFrame from a start timestamp and a duration.
func GetTimeFrame(from Timestamp, duration time.Duration) TimeFrame {
	to := int64(from) + int64(duration/time.Millisecond)
	return TimeFrame{From: from, To: Timestamp(to)}
}

// String returns a string representation of the TimeFrame in the format "start - end".
func (tf *TimeFrame) String(format string) string {
	return fmt.Sprintf("%s - %s", tf.From.String(format), tf.To.String(format))
}

// Duration returns the duration of the TimeFrame.
func (tf *TimeFrame) Duration() time.Duration {
	millis := int64(tf.To) - int64(tf.From)
	return time.Duration(millis) * time.Millisecond
}

// Overlaps returns true if the provided TimeFrame overlaps with the current TimeFrame.
func (tf *TimeFrame) Overlaps(with TimeFrame, gap time.Duration) bool {
	overlapStart := tf.Max(tf.From, with.From)
	overlapEnd := tf.Min(tf.To, with.To)

	overlapMs := int64(overlapEnd) - int64(overlapStart)
	return overlapMs > gap.Milliseconds()
}

// Max find max value between two timestamps
func (tf *TimeFrame) Max(a, b Timestamp) Timestamp {
	if a > b {
		return a
	}
	return b
}

// Min find min value between two timestamps
func (tf *TimeFrame) Min(a, b Timestamp) Timestamp {
	if a < b {
		return a
	}
	return b
}

// GetTimeFrameOf creates a new TimeFrame from the beginning of the period to the end of the period
// if year is provided and month and day are 0, the period is one year.
// if year and month are provided and day is 0, the period is one month.
// if year, month and day are provided (not 0), the period is one day.
func GetTimeFrameOf(year, month, day int) TimeFrame {

	if year == 0 {
		to := time.Date(5000, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
		return TimeFrame{From: 0, To: NewTimestamp(to)}
	}

	if month == 0 {
		from := time.Date(year, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
		to := from.Add(time.Hour * 24 * 365)
		return TimeFrame{From: NewTimestamp(from), To: NewTimestamp(to)}
	}

	if day == 0 {
		from := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		to := from.Add(time.Hour * 24)
		return TimeFrame{From: NewTimestamp(from), To: NewTimestamp(to)}
	}

	from := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	to := time.Date(year, time.Month(month), day, 23, 59, 59, 0, time.UTC)
	return TimeFrame{From: NewTimestamp(from), To: NewTimestamp(to)}
}

// endregion
