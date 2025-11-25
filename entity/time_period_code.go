package entity

// TimePeriodCode represents a code for a time period (e.g., Minute, Hour, Day).
// @Enum
type TimePeriodCode int

// timePeriodCode defines the available time period codes.
// @EnumValuesFor: TimePeriodCode
type timePeriodCode struct {
	// Undefined [0]
	UNDEFINED TimePeriodCode `value:"0"`

	// Minute [1]
	MINUTE TimePeriodCode `value:"1"`

	// Hour [2]
	HOUR TimePeriodCode `value:"2"`

	// Day [3]
	DAY TimePeriodCode `value:"3"`

	// Week [4]
	WEEK TimePeriodCode `value:"4"`

	// Month [5]
	MONTH TimePeriodCode `value:"5"`
}

// TimePeriodCodes is a singleton instance containing all available time period codes.
var TimePeriodCodes = &timePeriodCode{
	UNDEFINED: 0, // Undefined [0]
	MINUTE:    1, // Minute [1]
	HOUR:      2, // Hour[2]
	DAY:       3, // Day [3]
	WEEK:      4, // Week[4]
	MONTH:     5, // Month [5]
}
