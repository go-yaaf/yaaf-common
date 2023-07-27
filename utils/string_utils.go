package utils

import (
	"regexp"
	"strings"
)

// region Factory method -----------------------------------------------------------------------------------------------

// General string utils
type stringUtils struct {
}

// StringUtils factory method
func StringUtils() *stringUtils {
	return &stringUtils{}
}

// endregion

// region Public methods -----------------------------------------------------------------------------------------------

// WildCardToRegexp converts a wildcard string (including * and ?) to regular expression
func (t *stringUtils) WildCardToRegexp(pattern string) string {
	components := strings.Split(pattern, "*")
	if len(components) == 1 {
		// if len is 1, there are no *'s, return exact match pattern
		return "^" + pattern + "$"
	}
	var result strings.Builder
	for i, literal := range components {

		// Replace * with .*
		if i > 0 {
			result.WriteString(".*")
		}

		// Quote any regular expression meta characters in the
		// literal text.
		result.WriteString(regexp.QuoteMeta(literal))
	}
	return "^" + result.String() + "$"
}

// WildCardMatch returns true if the source string matches the wildcard pattern
func (t *stringUtils) WildCardMatch(source string, pattern string) bool {
	result, _ := regexp.MatchString(t.WildCardToRegexp(pattern), source)
	return result
}

// endregion
