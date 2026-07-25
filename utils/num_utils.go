// Package utils provides a collection of general-purpose utility functions.
// This file contains utilities for string manipulation, such as wildcard matching.
package utils

import (
	"math"
)

// region Factory method -----------------------------------------------------------------------------------------------

// NumUtilsStruct provides a set of utility functions for working with numbers.
// It is used as a singleton to offer a centralized and efficient way to perform numeric operations.
type NumUtilsStruct struct {
}

// NumUtils is a factory method that returns a singleton instance of NumUtilsStruct.
// This provides a single point of access to the string utilities.
func NumUtils() *NumUtilsStruct {
	// Note: In a highly concurrent application, consider using sync.Once for thread-safe singleton initialization.
	return &NumUtilsStruct{}
}

// endregion

// region Public methods -----------------------------------------------------------------------------------------------

// Round rounds a float64 number to the number of digits after the decimal point
func (t *NumUtilsStruct) Round(input float64, digits int) float64 {
	factor := math.Pow(10, float64(digits))
	return math.Round(input*factor) / factor
}

// endregion
