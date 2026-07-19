package entity

import (
	"bytes"
	"strings"
	"text/template"
	"time"
)

// region Timestamp ----------------------------------------------------------------------------------------------------

// Timestamp represents Epoch milliseconds timestamp.
// It is the primary time representation in the system, allowing for easy serialization and arithmetic.
// @Data
type Timestamp int64

// EpochNowMillis returns the current time as Epoch time in milliseconds, with an optional delta.
//
// Parameters:
//   - delta: A duration in milliseconds to add to the current time.
//
// Returns:
//   - The calculated Timestamp.
func EpochNowMillis(delta int64) Timestamp {
	return Timestamp((time.Now().UnixNano() / 1000000) + delta)
}

// Now returns the current time as Epoch time in milliseconds.
func Now() Timestamp {
	return EpochNowMillis(0)
}

// NewTimestamp creates a Timestamp from a standard Go time.Time object.
func NewTimestamp(t time.Time) Timestamp {
	return Timestamp(t.UnixNano() / 1000000)
}

// GetTimestamp creates a Timestamp from the provided date values
func GetTimestamp(year, month, day, hour, minute, second int) Timestamp {
	tm := time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC)
	return Timestamp(tm.UnixNano() / 1000000)
}

// Add adds a duration to the Timestamp and returns a new Timestamp.
func (ts Timestamp) Add(delta time.Duration) Timestamp {
	return Timestamp(int64(ts) + delta.Milliseconds())
}

// Time converts the Timestamp to a standard Go time.Time object.
func (ts Timestamp) Time() (result time.Time) {
	return time.UnixMilli(int64(ts))
}

// StartOfHour returns a new Timestamp representing the start of the hour for the current timestamp.
func (ts Timestamp) StartOfHour() Timestamp {
	t := ts.Time()
	sod := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
	return Timestamp(sod.UnixNano() / 1000000)
}

// EndOfHour returns a new Timestamp representing the end of the hour for the current timestamp.
func (ts Timestamp) EndOfHour() Timestamp {
	t := ts.Time()
	eod := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 59, 59, 9.99e+8, t.Location())
	return Timestamp(eod.UnixNano() / 1000000)
}

// StartOfDay returns a new Timestamp representing the start of the day (00:00:00) for the current timestamp.
func (ts Timestamp) StartOfDay() Timestamp {
	t := ts.Time()
	sod := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return Timestamp(sod.UnixNano() / 1000000)
}

// EndOfDay returns a new Timestamp representing the end of the day (23:59:59.999) for the current timestamp.
func (ts Timestamp) EndOfDay() Timestamp {
	t := ts.Time()
	eod := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 9.99e+8, t.Location())
	return Timestamp(eod.UnixNano() / 1000000)
}

// StartOfMonth returns a new Timestamp representing the start of the month for the current timestamp.
func (ts Timestamp) StartOfMonth() Timestamp {
	t := ts.Time()
	som := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return Timestamp(som.UnixNano() / 1000000)
}

// EndOfMonth returns a new Timestamp representing the end of the month for the current timestamp.
func (ts Timestamp) EndOfMonth() Timestamp {
	t := ts.Time()
	som := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	firstDayOfNextMonth := som.AddDate(0, 1, 0)
	eom := firstDayOfNextMonth.Add(-time.Second)
	return Timestamp(eom.UnixNano() / 1000000)
}

// convertISO8601Format converts a custom date format string (e.g., "YYYY-MM-DD") to Go's reference time format.
// This helper function allows using more familiar format strings instead of Go's specific reference date.
//
// Parameters:
//   - format: The format string using ISO 8601 like placeholders (YYYY, MM, DD, etc.).
//
// Returns:
//   - The Go reference time format string.
func (ts Timestamp) convertISO8601Format(format string) string {

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

// String converts the Timestamp to a string using the specified format.
// If format is empty, it uses RFC3339.
// It supports custom format strings via convertISO8601Format.
func (ts Timestamp) String(format string) string {
	if len(format) == 0 {
		return ts.Time().Format(time.RFC3339)
	} else {
		layout := ts.convertISO8601Format(format)
		return ts.Time().Format(layout)
	}
}

// LocalString converts the Timestamp to a string in a specific timezone.
//
// Parameters:
//   - format: The format string.
//   - tz: The timezone identifier (e.g., "America/New_York").
//
// Returns:
//   - The formatted time string in the specified timezone.
func (ts Timestamp) LocalString(format string, tz string) string {

	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.UTC
	}
	layout := ts.convertISO8601Format(format)
	return ts.Time().In(loc).Format(layout)
}

// endregion
