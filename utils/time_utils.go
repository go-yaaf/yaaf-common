package utils

import (
	"bytes"
	"strings"
	"text/template"
	"time"

	. "github.com/go-yaaf/yaaf-common/entity"
)

// TimeUtilsStruct internal helper
type TimeUtilsStruct struct {
	baseTime Timestamp
	SECOND   uint64
	MINUTE   uint64
	HOUR     uint64
	DAY      uint64
}

// TimeUtils is a factory method
func TimeUtils(ts Timestamp) *TimeUtilsStruct {
	return &TimeUtilsStruct{baseTime: ts, SECOND: 1000, MINUTE: 60 * 1000, HOUR: 60 * 60 * 1000, DAY: 24 * 60 * 60 * 1000}
}

// Get returns the current timestamp
func (t *TimeUtilsStruct) Get() Timestamp {
	return t.baseTime
}

// ConvertISO8601Format converts ISO6801 datetime format to Go RFC3339 format (used by Go)
func (t *TimeUtilsStruct) ConvertISO8601Format(format string) string {

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
func (t *TimeUtilsStruct) Format(format string) string {
	if len(format) == 0 {
		return time.UnixMilli(int64(t.baseTime)).Format(time.RFC3339)
	} else {
		layout := t.ConvertISO8601Format(format)
		return time.UnixMilli(int64(t.baseTime)).Format(layout)
	}
}

// SetInterval create periodic time triggered function call
func (t *TimeUtilsStruct) SetInterval(someFunc func(), milliseconds int, async bool) chan bool {

	// How often to fire the passed in function
	// in milliseconds
	interval := time.Duration(milliseconds) * time.Millisecond

	// Setup the ticket and the channel to signal
	// the ending of the interval
	ticker := time.NewTicker(interval)
	clearChan := make(chan bool)

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
			case <-clearChan:
				ticker.Stop()
				return
			}
		}
	}()

	// We return the channel so we can pass in
	// a value to it to clear the interval
	return clearChan
}

// LowerBound return the floor value of the timestamp to the lowest time duration
// Supported duration values:
// * time.Minute - get the lower bound by minute
// * time.Hour - get the lower bound by hour
// * time.Hour * 24 - get the lower bound by day
func (t *TimeUtilsStruct) LowerBound(duration time.Duration) *TimeUtilsStruct {
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
func (t *TimeUtilsStruct) UpperBound(duration time.Duration) *TimeUtilsStruct {
	tm := int64(t.baseTime) * int64(time.Millisecond)
	rem := tm - (tm % int64(duration)) + int64(duration)
	t.baseTime = Timestamp(rem / int64(time.Millisecond))
	return t
}

// GetSeries creates a time series from the base time to the end time with the given interval
func (t *TimeUtilsStruct) GetSeries(end Timestamp, interval time.Duration) (series []Timestamp) {

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
func (t *TimeUtilsStruct) GetSeriesMap(end Timestamp, interval time.Duration) map[Timestamp]int {

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
func (t *TimeUtilsStruct) GetTimeFrames(end Timestamp, interval time.Duration) (series []TimeFrame) {

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
func (t *TimeUtilsStruct) GetTimeFramesMap(end Timestamp, interval time.Duration) map[Timestamp]TimeFrame {

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

// region Time related math functions ----------------------------------------------------------------------------------

// DayRange create the boundaries of a day (start to end)
func (t *TimeUtilsStruct) DayRange() TimeFrame {
	current := t.baseTime.Time()
	sod := time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, current.Location())
	eod := time.Date(current.Year(), current.Month(), current.Day(), 23, 59, 59, 999, current.Location())
	return TimeFrame{
		From: NewTimestamp(sod),
		To:   NewTimestamp(eod),
	}
}

// MonthRange gets the local timestamp boundaries of a month
func (t *TimeUtilsStruct) MonthRange() TimeFrame {
	current := t.baseTime.Time()
	som := time.Date(current.Year(), current.Month(), 1, 0, 0, 0, 0, current.Location())
	firstDayOfNextMonth := som.AddDate(0, 1, 0)
	eom := firstDayOfNextMonth.Add(-time.Second)

	return TimeFrame{
		From: NewTimestamp(som),
		To:   NewTimestamp(eom),
	}
}

// endregion

// region General time utils functions ---------------------------------------------------------------------------------

// CalculateInterval returns the preferred time interval based on the time period
func CalculateInterval(from, to Timestamp) (rFrom, rTo Timestamp, interval time.Duration) {

	tuf := TimeUtils(from)
	tut := TimeUtils(to)

	delta := int64(to) - int64(from)
	if delta < 0 {
		delta = delta * (-1)
	}

	// If the period is 1 hour or less, the interval is minutes
	if delta <= int64(tuf.HOUR) {
		rFrom = tuf.LowerBound(time.Minute).Get()
		rTo = tut.LowerBound(time.Minute).Get()
		return rFrom, rTo, time.Minute
	}

	// If the period is 2 days or less, the interval is hours
	if delta <= int64(tuf.DAY)*2 {
		rFrom = tuf.LowerBound(time.Hour).Get()
		rTo = tut.LowerBound(time.Hour).Get()
		return rFrom, rTo, time.Hour
	}

	// If the period is 30 days or less, the interval is days
	if delta <= int64(tuf.DAY)*30 {
		rFrom = tuf.LowerBound(time.Hour * 24).Get()
		rTo = tut.LowerBound(time.Hour * 24).Get()
		return rFrom, rTo, time.Hour * 24
	}

	// If the period is 60 days or less, the interval is week
	if delta <= int64(tuf.DAY)*60 {
		rFrom = tuf.LowerBound(time.Hour * 24 * 7).Get()
		rTo = tut.LowerBound(time.Hour * 24 * 7).Get()
		return rFrom, rTo, time.Hour * 24 * 7
	}

	// If none of the above, the interval is month
	rFrom = tuf.LowerBound(time.Hour * 24 * 30).Get()
	rTo = tut.LowerBound(time.Hour * 24 * 30).Get()
	return rFrom, rTo, time.Hour * 24 * 30

}

// CreateTimeSeries creates a blank time series for the provided time period based on the time interval
func CreateTimeSeries[T any](name string, from, to Timestamp, interval time.Duration, initValue T) (series Entity) {

	intervalMs := int64(interval) / int64(time.Millisecond)

	period := TimeFrame{From: from, To: to}
	result := &TimeSeries[T]{Name: name, Range: period, Values: make([]TimeDataPoint[T], 0)}

	// truncate to the next interval
	start := int64(period.From)
	mod := start % intervalMs
	start = start - mod
	if mod > 0 {
		start = start + intervalMs
	}

	for ts := start; ts <= int64(period.To); ts += intervalMs {
		tdp := TimeDataPoint[T]{Timestamp: Timestamp(ts), Value: initValue}
		result.Values = append(result.Values, tdp)
	}
	return result
}

// HistogramTimeSeries converts histogram to a sorted and consecutive time series
func HistogramTimeSeries(name string, from, to Timestamp, interval time.Duration, hist map[Timestamp]Tuple[int64, float64]) (series Entity) {
	timeSeries := CreateTimeSeries[float64](name, from, to, interval, 0)
	for ts, val := range hist {
		timeSeries.(*TimeSeries[float64]).SetDataPoint(ts, val.Value)
	}
	return timeSeries
}

// endregion
