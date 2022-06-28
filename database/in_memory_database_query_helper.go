// Copyright 2022. Motty Cohen
//
// In-memory database query helpers (used for mock)
//
package database

import (
	"fmt"
	"strconv"
	"strings"
)

var operators map[queryOperator]FilterFunction

func init() {
	operators = make(map[queryOperator]FilterFunction)

	operators[Eq] = eq
	operators[Neq] = neq
	operators[Like] = like
	operators[Gt] = gt
	operators[Gte] = gte
	operators[Lt] = lt
	operators[Lte] = lte
	operators[In] = in
	operators[NotIn] = nin
	operators[Between] = between
	operators[Contains] = contains

}

// Signature of a filter function
type FilterFunction func(raw map[string]any, filter QueryFilter) bool

// Get filter function implementing an operator and a value
func testField(raw map[string]any, filter QueryFilter) bool {
	if filter.condition == false {
		return true
	}
	return operators[filter.operator](raw, filter)
}

// equal
func eq(raw map[string]any, filter QueryFilter) bool {
	if entityVal, ok := raw[filter.field]; ok {
		v1 := fmt.Sprintf("%v", entityVal)
		v2 := fmt.Sprintf("%v", filter.values[0])
		return v1 == v2
	} else {
		return false
	}
}

// not equal
func neq(raw map[string]any, filter QueryFilter) bool {
	if entityVal, ok := raw[filter.field]; ok {
		v1 := fmt.Sprintf("%v", entityVal)
		v2 := fmt.Sprintf("%v", filter.values[0])
		return v1 != v2
	} else {
		return false
	}
}

// like
func like(raw map[string]any, filter QueryFilter) bool {
	entityVal, ok := raw[filter.field]
	if !ok {
		return false
	}
	v1 := fmt.Sprintf("%v", entityVal)
	v2 := fmt.Sprintf("%v", filter.values[0])
	return strings.Contains(v1, v2)
}

// Greater than
func gt(raw map[string]any, filter QueryFilter) bool {
	entityVal, ok := raw[filter.field]
	if !ok {
		return false
	}
	v1 := fmt.Sprintf("%v", entityVal)
	v2 := fmt.Sprintf("%v", filter.values[0])

	n1, e1 := strconv.ParseFloat(v1, 64)
	n2, e2 := strconv.ParseFloat(v2, 64)

	if e1 != nil || e2 != nil {
		return false
	} else {
		return n1 > n2
	}
}

// less than
func lt(raw map[string]any, filter QueryFilter) bool {
	entityVal, ok := raw[filter.field]
	if !ok {
		return false
	}
	v1 := fmt.Sprintf("%v", entityVal)
	v2 := fmt.Sprintf("%v", filter.values[0])

	n1, e1 := strconv.ParseFloat(v1, 64)
	n2, e2 := strconv.ParseFloat(v2, 64)

	if e1 != nil || e2 != nil {
		return false
	} else {
		return n1 < n2
	}
}

// Greater than or equal
func gte(raw map[string]any, filter QueryFilter) bool {
	entityVal, ok := raw[filter.field]
	if !ok {
		return false
	}
	v1 := fmt.Sprintf("%v", entityVal)
	v2 := fmt.Sprintf("%v", filter.values[0])

	n1, e1 := strconv.ParseFloat(v1, 64)
	n2, e2 := strconv.ParseFloat(v2, 64)

	if e1 != nil || e2 != nil {
		return false
	} else {
		return n1 >= n2
	}
}

// less than or equal
func lte(raw map[string]any, filter QueryFilter) bool {
	entityVal, ok := raw[filter.field]
	if !ok {
		return false
	}
	v1 := fmt.Sprintf("%v", entityVal)
	v2 := fmt.Sprintf("%v", filter.values[0])

	n1, e1 := strconv.ParseFloat(v1, 64)
	n2, e2 := strconv.ParseFloat(v2, 64)

	if e1 != nil || e2 != nil {
		return false
	} else {
		return n1 <= n2
	}
}

// in (value should be an array)
func in(raw map[string]any, filter QueryFilter) bool {
	entityVal, ok := raw[filter.field]
	if !ok {
		return false
	}

	v1 := fmt.Sprintf("%v", entityVal)

	// Test for int array
	if arr, ok := filter.values[0].([]int); ok {
		if n, e := strconv.Atoi(v1); e != nil {
			return false
		} else {
			for _, t := range arr {
				if n == t {
					return true
				}
			}
			return false
		}
	}

	// Test for string array
	if arr, ok := filter.values[0].([]string); ok {
		for _, t := range arr {
			if v1 == t {
				return true
			}
		}
		return false
	}
	return false
}

// not in (value should be an array)
func nin(raw map[string]any, filter QueryFilter) bool {
	entityVal, ok := raw[filter.field]
	if !ok {
		return false
	}

	v1 := fmt.Sprintf("%v", entityVal)

	// Test for int array
	if arr, ok := filter.values[0].([]int); ok {
		if n, e := strconv.Atoi(v1); e != nil {
			return false
		} else {
			for _, t := range arr {
				if n == t {
					return false
				}
			}
			return true
		}
	}

	// Test for string array
	if arr, ok := filter.values[0].([]string); ok {
		for _, t := range arr {
			if v1 == t {
				return false
			}
		}
		return true
	}
	return false
}

// array field contains the tested value
func contains(raw map[string]any, filter QueryFilter) bool {
	entityVal, ok := raw[filter.field]
	if !ok {
		return false
	}

	v1 := fmt.Sprintf("%v", filter.values[0])

	// Test for int array field
	if arr, ok := entityVal.([]int); ok {
		if n, e := strconv.Atoi(v1); e != nil {
			return false
		} else {
			for _, t := range arr {
				if n == t {
					return true
				}
			}
			return false
		}
	}

	// Test for string array field
	if arr, ok := filter.values[0].([]string); ok {
		for _, t := range arr {
			if v1 == t {
				return true
			}
		}
		return false
	}
	return false
}

// between (the expected value is comma-separated list of 2 values
func between(raw map[string]any, filter QueryFilter) bool {

	entityVal, ok := raw[filter.field]
	if !ok {
		return false
	}
	v1 := fmt.Sprintf("%v", entityVal)

	val1 := fmt.Sprintf("%v", filter.values[0])
	val2 := fmt.Sprintf("%v", filter.values[1])

	n1, e1 := strconv.Atoi(strings.TrimSpace(val1))
	n2, e2 := strconv.Atoi(strings.TrimSpace(val2))
	t0, e3 := strconv.Atoi(strings.TrimSpace(v1))

	if e1 != nil || e2 != nil || e3 != nil {
		return false
	} else {
		return n1 <= t0 && t0 <= n2
	}
}

// endregion
