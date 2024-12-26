package forms

import (
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// EmailRX is a compiled regular expression to validate the format of email addresses according to standard RFC 5322 rules.
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// NotBlank checks if the provided string is not empty or whitespace-only
// and returns true if it contains non-whitespace characters.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars checks if the number of characters in a string is
// less than or equal to a specified limit.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermittedValue checks if a given value is present within a
// list of permitted values and returns true if found.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// ValidDate checks if the provided date string represents a valid,
// non-zero date and returns true if valid, otherwise false.
func ValidDate(date time.Time) bool {
	if date.IsZero() {
		return false
	}
	return true
}

// MinChars checks if the string `value` contains at least `n` characters,
// returning true if the condition is met.
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// Matches checks if the given string value matches
// the specified regular expression pattern.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// ValidNumber checks if the given string represents a valid integer and returns true if it does, otherwise false.
func ValidNumber(value string) bool {
	_, err := strconv.ParseInt(value, 10, 64)
	return err == nil
}

// MaxNumber checks if the numerical value in the given string does not exceed the specified maximum integer value.
func MaxNumber(value int, max int) bool {
	return value <= max
}
