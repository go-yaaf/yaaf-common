package utils

import (
	"bytes"
	"strings"
	"text/template"
	"time"

	. "github.com/go-yaaf/yaaf-common/entity"
)

// timeUtils internal helper
type timeUtils struct {
	baseTime Timestamp
	SECOND   uint64
	MINUTE   uint64
	HOUR     uint64
	DAY      uint64
}

// TimeUtils is a factory method
func TimeUtils(ts Timestamp) *timeUtils {
	return &timeUtils{baseTime: ts, SECOND: 1000, MINUTE: 60 * 1000, HOUR: 60 * 60 * 1000, DAY: 24 * 60 * 60 * 1000}
}

// Get returns the current timestamp
func (t *timeUtils) Get() Timestamp {
	return t.baseTime
}

// ConvertISO8601Format converts ISO6801 datetime format to Go RFC3339 format (used by Go)
func (t *timeUtils) ConvertISO8601Format(format string) string {

	data := struct {
		V1  string
		V2  string
		V3  string
		V4  string
		V5  string
		V6  string
		V7  string
		V8  string
		V9  string
		V10 string
		V11 string
		V12 string
		V13 string
		V14 string
		V15 string
		V16 string
		V17 string
		V18 string
		V19 string
	}{
		"2006", "06",
		"January", "Jan", "01", "1",
		"Monday", "Mon", "02", "2",
		"15", "3",
		"04", "4",
		"05", "5",
		"MST", "-0700",
		"PM",
	}

	tmpl := format

	tmpl = strings.ReplaceAll(tmpl, "YYYY", "{{.V1}}")
	tmpl = strings.ReplaceAll(tmpl, "yyyy", "{{.V1}}")
	tmpl = strings.ReplaceAll(tmpl, "YY", "{{.V2}}")
	tmpl = strings.ReplaceAll(tmpl, "yy", "{{.V2}}")
	tmpl = strings.ReplaceAll(tmpl, "MMMM", "{{.V3}}")
	tmpl = strings.ReplaceAll(tmpl, "MMM", "{{.V4}}")
	tmpl = strings.ReplaceAll(tmpl, "MM", "{{.V5}}")
	tmpl = strings.ReplaceAll(tmpl, "M", "{{.V6}}")

	tmpl = strings.ReplaceAll(tmpl, "dddd", "{{.V7}}")
	tmpl = strings.ReplaceAll(tmpl, "DDDD", "{{.V7}}")

	tmpl = strings.ReplaceAll(tmpl, "ddd", "{{.V8}}")
	tmpl = strings.ReplaceAll(tmpl, "DDD", "{{.V8}}")

	tmpl = strings.ReplaceAll(tmpl, "dd", "{{.V9}}")
	tmpl = strings.ReplaceAll(tmpl, "DD", "{{.V9}}")

	tmpl = strings.ReplaceAll(tmpl, "d", "{{.V10}}")
	tmpl = strings.ReplaceAll(tmpl, "D", "{{.V10}}")

	tmpl = strings.ReplaceAll(tmpl, "HH", "{{.V11}}")
	tmpl = strings.ReplaceAll(tmpl, "hh", "{{.V11}}")
	tmpl = strings.ReplaceAll(tmpl, "H", "{{.V12}}")
	tmpl = strings.ReplaceAll(tmpl, "h", "{{.V12}}")

	tmpl = strings.ReplaceAll(tmpl, "mm", "{{.V13}}")
	tmpl = strings.ReplaceAll(tmpl, "m", "{{.V14}}")
	tmpl = strings.ReplaceAll(tmpl, "ss", "{{.V15}}")
	tmpl = strings.ReplaceAll(tmpl, "s", "{{.V16}}")

	tmpl = strings.ReplaceAll(tmpl, "TZD", "{{.V17}}")
	tmpl = strings.ReplaceAll(tmpl, "z", "{{.V17}}")
	tmpl = strings.ReplaceAll(tmpl, "Z", "{{.V18}}")
	tmpl = strings.ReplaceAll(tmpl, "a", "{{.V19}}")
	tmpl = strings.ReplaceAll(tmpl, "A", "{{.V19}}")

	tm := template.Must(template.New("format").Parse(tmpl))

	var buff bytes.Buffer
	if err := tm.Execute(&buff, data); err != nil {
		return ""
	} else {
		return buff.String()
	}
}

// Format converts Epoch milliseconds timestamp to readable string, if format is empty, the default layout (RFC3339) is used
func (t *timeUtils) Format(format string) string {
	if len(format) == 0 {
		return time.UnixMilli(int64(t.baseTime)).Format(time.RFC3339)
	} else {
		layout := t.ConvertISO8601Format(format)
		return time.UnixMilli(int64(t.baseTime)).Format(layout)
	}
}

// SetInterval create periodic time triggered function call
func (t *timeUtils) SetInterval(someFunc func(), milliseconds int, async bool) chan bool {

	// How often to fire the passed in function
	// in milliseconds
	interval := time.Duration(milliseconds) * time.Millisecond

	// Setup the ticket and the channel to signal
	// the ending of the interval
	ticker := time.NewTicker(interval)
	clear := make(chan bool)

	// Put the selection in a go routine
	// so that the for loop is none blocking
	go func() {
		for {
			select {
			case <-ticker.C:
				if async {
					// This won't block
					go someFunc()
				} else {
					// This will block
					someFunc()
				}
			case <-clear:
				ticker.Stop()
				return
			}
		}
	}()

	// We return the channel so we can pass in
	// a value to it to clear the interval
	return clear
}

// LowerBound return the floor value of the timestamp to the lowest time duration
// Supported duration values:
// * time.Minute - get the lower bound by minute
// * time.Hour - get the lower bound by hour
// * time.Hour * 24 - get the lower bound by day
func (t *timeUtils) LowerBound(duration time.Duration) *timeUtils {
	tm := int64(t.baseTime) * int64(time.Millisecond)
	rem := tm - (tm % int64(duration))
	t.baseTime = Timestamp(rem / int64(time.Millisecond))
	return t
}

// UpperBound return the ceiling value of the timestamp to the next time duration
// Supported duration values:
// * time.Minute - get the upper bound by minute
// * time.Hour - get the upper bound by hour
// * time.Hour * 24 - get the upper bound by day
func (t *timeUtils) UpperBound(duration time.Duration) *timeUtils {
	tm := int64(t.baseTime) * int64(time.Millisecond)
	rem := tm - (tm % int64(duration)) + int64(duration)
	t.baseTime = Timestamp(rem / int64(time.Millisecond))
	return t
}

// GetSeries creates a time series from the base time to the end time with the given interval
func (t *timeUtils) GetSeries(end Timestamp, interval time.Duration) (series []Timestamp) {

	if interval == 0 {
		return series
	}

	from := int64(t.baseTime)
	// Add 1 to include the last slot
	to := int64(end)
	step := int64(interval / time.Millisecond)

	if from < to {
		eot := to + 1
		for ts := from; ts < eot; ts += step {
			series = append(series, Timestamp(ts))
		}
	} else {
		for ts := from; ts > to; ts -= step {
			series = append(series, Timestamp(ts))
		}
	}
	return series
}

// GetSeries2 creates a time series from the base time to the end time with the given interval
// This is fix  for GetSeries function which misses "last slot" due to incorrect "for"
// loop completion condition. For example, if we need series from day 01 to day 06,
// it will return series from day 01 to day 05 ( see line "eot := to + 1", it should be rather
// eot := to + step, or loop completion condition should be "ts =< to "). Since GetSeries function
// is widely used, to avoid regressions, I've created GetSeries2 until further clarifications
func (t *timeUtils) GetSeries2(end Timestamp, interval time.Duration) (series []Timestamp) {

	if interval == 0 {
		return series
	}

	from := int64(t.baseTime)
	to := int64(end)
	step := int64(interval / time.Millisecond)

	if from < to {
		eot := to + step
		for ts := from; ts < eot; ts += step {
			series = append(series, Timestamp(ts))
		}
	} else {
		for ts := from; ts > to; ts -= step {
			series = append(series, Timestamp(ts))
		}
	}
	return series
}

// GetSeriesMap creates a time series from the base time to the end time with the given interval as a map
func (t *timeUtils) GetSeriesMap(end Timestamp, interval time.Duration) map[Timestamp]int {

	series := make(map[Timestamp]int)
	if interval == 0 {
		return series
	}

	from := int64(t.baseTime)
	to := int64(end)
	step := int64(interval / time.Millisecond)

	if from < to {
		for ts := from; ts < to; ts += step {
			series[Timestamp(ts)] = 0
		}
	} else {
		for ts := from; ts > to; ts -= step {
			series[Timestamp(ts)] = 0
		}
	}
	return series
}

// GetTimeFrames creates time frames from the base time to the end time with the given interval with delay between slots
func (t *timeUtils) GetTimeFrames(end Timestamp, interval time.Duration) (series []TimeFrame) {

	if interval == 0 {
		return series
	}

	from := int64(t.baseTime)
	to := int64(end)
	step := int64(interval / time.Millisecond)

	prev := int64(-1)

	if from < to {
		for ts := from; ts < to; ts += step {
			if prev < 0 {
				prev = ts
			} else {
				series = append(series, NewTimeFrame(Timestamp(prev), Timestamp(ts)))
				prev = ts
			}
		}
	} else {
		for ts := from; ts > to; ts -= step {
			if prev < 0 {
				prev = ts
			} else {
				series = append(series, NewTimeFrame(Timestamp(ts), Timestamp(prev)))
				prev = ts
			}
		}
	}
	return series
}

// GetTimeFramesMap creates time frames from the base time to the end time with the given interval as a map
func (t *timeUtils) GetTimeFramesMap(end Timestamp, interval time.Duration) map[Timestamp]TimeFrame {

	frames := make(map[Timestamp]TimeFrame)
	if interval == 0 {
		return frames
	}

	from := int64(t.baseTime)
	to := int64(end)
	step := int64(interval / time.Millisecond)

	prev := int64(-1)

	if from < to {
		for ts := from; ts < to; ts += step {
			if prev < 0 {
				prev = ts
			} else {
				frames[Timestamp(prev)] = NewTimeFrame(Timestamp(prev), Timestamp(ts))
				prev = ts
			}
		}
	} else {
		for ts := from; ts > to; ts -= step {
			if prev < 0 {
				prev = ts
			} else {
				frames[Timestamp(ts)] = NewTimeFrame(Timestamp(ts), Timestamp(prev))
				prev = ts
			}
		}
	}
	return frames
}
