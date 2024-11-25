package entity

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// region LocalTimes ----------------------------------------------------------------------------------------------------

// LocalTime represents the time (year, month, day, hour, minute, second) as number in the following format: YYYYMMDDhhmmss
type LocalTime int64

// LocalTimestamp converts LocalTime to Timestamp
func LocalTimestamp(lt ...LocalTime) Timestamp {
	if len(lt) > 0 {
		return Timestamp(lt[0])
	}
	utc := time.Now().UTC()
	result := utc.Year() * 10000000000
	result += int(utc.Month()) * 100000000
	result += utc.Day() * 1000000
	result += utc.Hour() * 10000
	result += utc.Minute() * 100
	result += utc.Second()

	return Timestamp(result)
}

// FromTime convert time to LocalTime
func FromTime(t time.Time) LocalTime {
	utc := t.UTC()
	result := utc.Year() * 10000000000
	result += int(utc.Month()) * 100000000
	result += utc.Day() * 1000000
	result += utc.Hour() * 10000
	result += utc.Minute() * 100
	result += utc.Second()

	return LocalTime(result)
}

// FromTimestamp convert Timestamp to LocalTime
func FromTimestamp(ts Timestamp) LocalTime {
	return FromTime(ts.Time())
}

// LocalNow return current time as YYYYMMDDhhmmss
func LocalNow() LocalTime {
	utc := time.Now().UTC()
	result := utc.Year() * 10000000000
	result += int(utc.Month()) * 100000000
	result += utc.Day() * 1000000
	result += utc.Hour() * 10000
	result += utc.Minute() * 100
	result += utc.Second()

	return LocalTime(result)
}

// Add time and return a new timestamp
func (lt *LocalTime) Add(delta time.Duration) LocalTime {
	return LocalTime(int64(*lt) + delta.Milliseconds())
}

// Time returns the Go primitive  time.Time object
func (lt *LocalTime) Time() (result time.Time) {
	year, month, day, hours, minutes, seconds := lt.Split()
	return time.Date(year, time.Month(month), day, hours, minutes, seconds, 0, time.UTC)
}

// Split timestamp string in format of: hh:mm to hour and minute
func (lt *LocalTime) Split() (year, month, day, hours, minutes, seconds int) {

	str := fmt.Sprintf("%d", *lt)
	if len(str) >= 4 {
		year, _ = strconv.Atoi(str[:4])
	}
	if len(str) >= 6 {
		month, _ = strconv.Atoi(str[4:6])
	}
	if len(str) >= 8 {
		day, _ = strconv.Atoi(str[6:8])
	}
	if len(str) >= 10 {
		hours, _ = strconv.Atoi(str[8:10])
	}
	if len(str) >= 12 {
		minutes, _ = strconv.Atoi(str[10:12])
	}
	if len(str) >= 14 {
		seconds, _ = strconv.Atoi(str[12:14])
	}
	return
}

// LocalString convert Epoch milliseconds timestamp with timezone (IANA) to readable string
func (lt *LocalTime) String(format string) string {
	layout := ConvertISO8601Format(format)
	return lt.Time().Format(layout)
}

// ConvertISO8601Format converts ISO6801 datetime format to Go RFC3339 format (used by Go)
func ConvertISO8601Format(format string) string {

	if len(format) == 0 {
		return "2006-01-02 15:04:05"
	}
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

// Timestamp converts local time to Timestamp
func (lt *LocalTime) Timestamp() Timestamp {
	return Timestamp(*lt)
}

// endregion
