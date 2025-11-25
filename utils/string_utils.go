// Package utils provides a collection of general-purpose utility functions.
// This file contains utilities for string manipulation, such as wildcard matching.
package utils

import (
	"regexp"
	"strings"
)

// region Factory method -----------------------------------------------------------------------------------------------

// StringUtilsStruct provides a set of utility functions for working with strings.
// It is used as a singleton to offer a centralized and efficient way to perform string operations.
type StringUtilsStruct struct {
}

// StringUtils is a factory method that returns a singleton instance of StringUtilsStruct.
// This provides a single point of access to the string utilities.
func StringUtils() *StringUtilsStruct {
	// Note: In a highly concurrent application, consider using sync.Once for thread-safe singleton initialization.
	return &StringUtilsStruct{}
}

// endregion

// region Public methods -----------------------------------------------------------------------------------------------

// WildCardToRegexp converts a wildcard pattern string into a regular expression string.
// This implementation supports the '*' wildcard, which matches any sequence of characters.
// Other characters are treated as literals. The resulting regex is anchored to match the entire string.
//
// For example:
//
//	"data*.log" -> "^data.*\.log$"
//	"report.*"  -> "^report\..*$"
//
// Parameters:
//
//	pattern: The wildcard string to convert.
//
// Returns:
//
//	A string containing the equivalent regular expression.
func (t *StringUtilsStruct) WildCardToRegexp(pattern string) string {
	var result strings.Builder
	result.WriteString("^")

	// Split the pattern by the '*' wildcard.
	parts := strings.Split(pattern, "*")
	for i, part := range parts {
		if i > 0 {
			// Replace '*' with '.*' which matches any sequence of zero or more characters.
			result.WriteString(".*")
		}
		// Escape any special regex characters in the literal parts of the pattern
		// to ensure they are treated as literal characters.
		result.WriteString(regexp.QuoteMeta(part))
	}

	result.WriteString("$")
	return result.String()
}

// WildCardMatch checks if a given source string matches a wildcard pattern.
// The pattern can contain '*' to match any sequence of characters.
//
// Parameters:
//
//	source: The string to test.
//	pattern: The wildcard pattern to match against.
//
// Returns:
//
//	`true` if the source string matches the pattern, `false` otherwise.
func (t *StringUtilsStruct) WildCardMatch(source string, pattern string) bool {
	regexPattern := t.WildCardToRegexp(pattern)
	// MatchString compiles the regex and returns true if it matches the source string.
	// We can ignore the error here as WildCardToRegexp is designed to produce a valid regex.
	matched, _ := regexp.MatchString(regexPattern, source)
	return matched
}

// endregion
