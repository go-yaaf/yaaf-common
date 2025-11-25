// Package utils provides a collection of utility functions, including helpers for time manipulation.
package utils

import (
	"strings"
	"time"

	. "github.com/go-yaaf/yaaf-common/entity"
)

// TimeUtilsStruct provides a fluent interface for time-related operations.
// It encapsulates a base time and offers methods for formatting, manipulation, and time series generation.
type TimeUtilsStruct struct {
	baseTime Timestamp
	SECOND   uint64
	MINUTE   uint64
	HOUR     uint64
	DAY      uint64
}

// TimeUtils is a factory function that creates a new TimeUtilsStruct with a given base time.
//
// Parameters:
//
//	ts: The base timestamp for time operations.
//
// Returns:
//
//	A new TimeUtilsStruct instance.
func TimeUtils(ts Timestamp) *TimeUtilsStruct {
	return &TimeUtilsStruct{
		baseTime: ts,
		SECOND:   1000,
		MINUTE:   60 * 1000,
		HOUR:     60 * 60 * 1000,
		DAY:      24 * 60 * 60 * 1000,
	}
}

// Get returns the current base timestamp of the TimeUtilsStruct.
func (t *TimeUtilsStruct) Get() Timestamp {
	return t.baseTime
}

// ConvertISO8601Format converts a date-time format string from ISO 8601 style to Go's reference time format.
// This allows for more intuitive format strings, such as "YYYY-MM-DD" instead of "2006-01-02".
//
// Parameters:
//
//	format: The ISO 8601-style format string.
//
// Returns:
//
//	The equivalent Go reference time format string.
func (t *TimeUtilsStruct) ConvertISO8601Format(format string) string {
	replacements := map[string]string{
		"YYYY": "2006", "yyyy": "2006", "YY": "06", "yy": "06",
		"MMMM": "January", "MMM": "Jan", "MM": "01", "M": "1",
		"dddd": "Monday", "DDDD": "Monday", "ddd": "Mon", "DDD": "Mon",
		"dd": "02", "DD": "02", "d": "2", "D": "2",
		"HH": "15", "hh": "03", "H": "15", "h": "3",
		"mm": "04", "m": "4",
		"ss": "05", "s": "5",
		"TZD": "MST", "z": "MST", "Z": "-0700",
		"a": "pm", "A": "PM",
	}

	for k, v := range replacements {
		format = strings.ReplaceAll(format, k, v)
	}
	return format
}

// Format converts the base timestamp to a formatted string.
// If no format is provided, it defaults to RFC3339.
//
// Parameters:
//
//	format: The desired output format, e.g., "YYYY-MM-DD hh:mm:ss".
//
// Returns:
//
//	The formatted time string.
func (t *TimeUtilsStruct) Format(format string) string {
	if format == "" {
		return t.baseTime.Time().Format(time.RFC3339)
	}
	layout := t.ConvertISO8601Format(format)
	return t.baseTime.Time().Format(layout)
}

// SetInterval executes a function at a specified interval.
// It can run the function synchronously or asynchronously.
//
// Parameters:
//
//	someFunc: The function to execute.
//	milliseconds: The interval in milliseconds.
//	async: If true, the function is executed in a new goroutine.
//
// Returns:
//
//	A channel that can be used to stop the interval by sending a boolean value.
func (t *TimeUtilsStruct) SetInterval(someFunc func(), milliseconds int, async bool) chan bool {
	interval := time.Duration(milliseconds) * time.Millisecond
	ticker := time.NewTicker(interval)
	clearChan := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				if async {
					go someFunc()
				} else {
					someFunc()
				}
			case <-clearChan:
				ticker.Stop()
				return
			}
		}
	}()
	return clearChan
}

// LowerBound rounds the base time down to the nearest specified duration.
//
// Parameters:
//
//	duration: The duration to round down to (e.g., time.Minute, time.Hour).
//
// Returns:
//
//	The TimeUtilsStruct instance for chaining.
func (t *TimeUtilsStruct) LowerBound(duration time.Duration) *TimeUtilsStruct {
	t.baseTime = Timestamp(t.baseTime.Time().Truncate(duration).UnixMilli())
	return t
}

// UpperBound rounds the base time up to the nearest specified duration.
//
// Parameters:
//
//	duration: The duration to round up to (e.g., time.Minute, time.Hour).
//
// Returns:
//
//	The TimeUtilsStruct instance for chaining.
func (t *TimeUtilsStruct) UpperBound(duration time.Duration) *TimeUtilsStruct {
	truncated := t.baseTime.Time().Truncate(duration)
	if truncated.Equal(t.baseTime.Time()) {
		t.baseTime = Timestamp(truncated.UnixMilli())
	} else {
		t.baseTime = Timestamp(truncated.Add(duration).UnixMilli())
	}
	return t
}

// GetSeries creates a slice of timestamps from the base time to a specified end time, at a given interval.
//
// Parameters:
//
//	end: The end time of the series.
//	interval: The duration between each timestamp in the series.
//
// Returns:
//
//	A slice of timestamps.
func (t *TimeUtilsStruct) GetSeries(end Timestamp, interval time.Duration) []Timestamp {
	var series []Timestamp
	if interval <= 0 {
		return series
	}

	current := t.baseTime.Time()
	endTime := end.Time()

	if current.Before(endTime) {
		for !current.After(endTime) {
			series = append(series, NewTimestamp(current))
			current = current.Add(interval)
		}
	} else {
		for !current.Before(endTime) {
			series = append(series, NewTimestamp(current))
			current = current.Add(-interval)
		}
	}
	return series
}

// GetSeriesMap creates a map of timestamps from the base time to a specified end time, at a given interval.
// The map values are initialized to zero.
//
// Parameters:
//
//	end: The end time of the series.
//	interval: The duration between each timestamp.
//
// Returns:
//
//	A map where keys are timestamps and values are 0.
func (t *TimeUtilsStruct) GetSeriesMap(end Timestamp, interval time.Duration) map[Timestamp]int {
	series := make(map[Timestamp]int)
	if interval <= 0 {
		return series
	}

	for _, ts := range t.GetSeries(end, interval) {
		series[ts] = 0
	}
	return series
}

// GetTimeFrames creates a slice of time frames from the base time to a specified end time, at a given interval.
//
// Parameters:
//
//	end: The end time of the series.
//	interval: The duration of each time frame.
//
// Returns:
//
//	A slice of TimeFrame structs.
func (t *TimeUtilsStruct) GetTimeFrames(end Timestamp, interval time.Duration) []TimeFrame {
	var series []TimeFrame
	if interval <= 0 {
		return series
	}

	timestamps := t.GetSeries(end, interval)
	for i := 0; i < len(timestamps)-1; i++ {
		series = append(series, NewTimeFrame(timestamps[i], timestamps[i+1]))
	}
	return series
}

// GetTimeFramesMap creates a map of time frames from the base time to a specified end time, at a given interval.
//
// Parameters:
//
//	end: The end time of the series.
//	interval: The duration of each time frame.
//
// Returns:
//
//	A map where keys are the start of the time frame and values are the TimeFrame structs.
func (t *TimeUtilsStruct) GetTimeFramesMap(end Timestamp, interval time.Duration) map[Timestamp]TimeFrame {
	frames := make(map[Timestamp]TimeFrame)
	for _, tf := range t.GetTimeFrames(end, interval) {
		frames[tf.From] = tf
	}
	return frames
}

// DayRange returns a TimeFrame representing the start and end of the day for the base time.
func (t *TimeUtilsStruct) DayRange() TimeFrame {
	year, month, day := t.baseTime.Time().Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, t.baseTime.Time().Location())
	endOfDay := startOfDay.AddDate(0, 0, 1).Add(-time.Nanosecond)
	return NewTimeFrame(NewTimestamp(startOfDay), NewTimestamp(endOfDay))
}

// MonthRange returns a TimeFrame representing the start and end of the month for the base time.
func (t *TimeUtilsStruct) MonthRange() TimeFrame {
	year, month, _ := t.baseTime.Time().Date()
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, t.baseTime.Time().Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return NewTimeFrame(NewTimestamp(startOfMonth), NewTimestamp(endOfMonth))
}

// CalculateInterval determines an appropriate time interval for a given time range.
// The interval is chosen based on the duration of the range (e.g., minutes for an hour, hours for a day).
func CalculateInterval(from, to Timestamp) (rFrom, rTo Timestamp, interval time.Duration) {

	delta := time.Duration(int64(to) - int64(from))
	if delta < 0 {
		delta = delta * (-1)
	}

	switch {
	case delta <= time.Hour:
		interval = time.Minute
	case delta <= 2*24*time.Hour:
		interval = time.Hour
	case delta <= 30*24*time.Hour:
		interval = 24 * time.Hour
	case delta <= 60*24*time.Hour:
		interval = 7 * 24 * time.Hour
	default:
		interval = 30 * 24 * time.Hour
	}

	rFrom = Timestamp(from.Time().Truncate(interval).UnixMilli())
	rTo = Timestamp(to.Time().Truncate(interval).UnixMilli())
	return
}

// CreateTimeSeries creates a new TimeSeries entity with a given name, time range, and interval.
// The values of the time series are initialized with a provided initial value.
func CreateTimeSeries[T any](name string, from, to Timestamp, interval time.Duration, initValue T) *TimeSeries[T] {
	series := &TimeSeries[T]{
		Name:   name,
		Range:  NewTimeFrame(from, to),
		Values: make([]TimeDataPoint[T], 0),
	}

	if interval <= 0 {
		return series
	}

	current := from.Time().Truncate(interval)
	if from.Time().After(current) {
		current = current.Add(interval)
	}

	for !current.After(to.Time()) {
		series.Values = append(series.Values, TimeDataPoint[T]{Timestamp: NewTimestamp(current), Value: initValue})
		current = current.Add(interval)
	}
	return series
}

// HistogramTimeSeries converts a histogram (a map of timestamps to values) into a sorted and continuous TimeSeries.
func HistogramTimeSeries(name string, from, to Timestamp, interval time.Duration, hist map[Timestamp]Tuple[int64, float64]) *TimeSeries[float64] {
	timeSeries := CreateTimeSeries[float64](name, from, to, interval, 0.0)
	for ts, val := range hist {
		timeSeries.SetDataPoint(ts, val.Value)
	}
	return timeSeries
}
