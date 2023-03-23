package entity

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"
)

// region Timestamp ----------------------------------------------------------------------------------------------------

// Timestamp represents Epoch milliseconds timestamp
type Timestamp int64

// EpochNowMillis return current time as Epoch time milliseconds with delta in millis
func EpochNowMillis(delta int64) Timestamp {
	return Timestamp((time.Now().UnixNano() / 1000000) + delta)
}

// Now return current time as Epoch time milliseconds with delta in millis
func Now() Timestamp {
	return EpochNowMillis(0)
}

// Add time and return a new timestamp
func (ts *Timestamp) Add(delta time.Duration) Timestamp {
	return Timestamp(int64(*ts) + delta.Milliseconds())
}

// Time returns the Go primitive  time.Time object
func (ts *Timestamp) Time() (result time.Time) {
	return time.UnixMilli(int64(*ts))
}

// Convert ISO6801 datetime format to Go RFC3339 format (used by Go)
// @param format ISO 8601 format
// @return RFC3339 format (using magic date sample: Jan 02 3:04:05 2006 -0700)
func (ts *Timestamp) convertISO8601Format(format string) string {

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

// String convert Epoch milliseconds timestamp to readable string
func (ts *Timestamp) String(format string) string {
	if len(format) == 0 {
		return ts.Time().Format(time.RFC3339)
	} else {
		layout := ts.convertISO8601Format(format)
		return ts.Time().Format(layout)
	}
}

// LocalString convert Epoch milliseconds timestamp with timezone (IANA) to readable string
func (ts *Timestamp) LocalString(format string, tz string) string {

	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.UTC
	}
	layout := ts.convertISO8601Format(format)
	return ts.Time().In(loc).Format(layout)
}

// endregion

// region TimeFrame ----------------------------------------------------------------------------------------------------

// TimeFrame represents a slot in time
type TimeFrame struct {
	From Timestamp
	To   Timestamp
}

// NewTimeFrame return new time slot using start and end time
func NewTimeFrame(from, to Timestamp) TimeFrame {
	return TimeFrame{From: from, To: to}
}

// GetTimeFrame return new time slot using start and duration
func GetTimeFrame(from Timestamp, duration time.Duration) TimeFrame {
	to := int64(from) + int64(duration/time.Millisecond)
	return TimeFrame{From: from, To: Timestamp(to)}
}

// String convert Epoch milliseconds timestamp to readable string
func (tf *TimeFrame) String(format string) string {
	return fmt.Sprintf("%s - %s", tf.From.String(format), tf.To.String(format))
}

// Duration of the timeframe
func (tf *TimeFrame) Duration() time.Duration {
	millis := int64(tf.To) - int64(tf.From)
	return time.Duration(millis) * time.Millisecond
}

// endregion
