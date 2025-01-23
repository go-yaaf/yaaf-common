package entity

type TimePeriodCode int

type timePeriodCode struct {
	UNDEFINED TimePeriodCode `Undefined[0]`
	MINUTE    TimePeriodCode `Minute[1]`
	HOUR      TimePeriodCode `HOUR[2]`
	DAY       TimePeriodCode `DAY[3]`
	WEEK      TimePeriodCode `WEEK[4]`
	MONTH     TimePeriodCode `MONTH[5]`
}

var TimePeriodCodes = &timePeriodCode{
	UNDEFINED: 0, //Undefined[0]
	MINUTE:    1, //Minute[1]
	HOUR:      2, //HOUR[2]
	DAY:       3, //DAY[3]
	WEEK:      4, //WEEK[4]
	MONTH:     5, //MONTH[5]
}
