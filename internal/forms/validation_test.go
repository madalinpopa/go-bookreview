package forms

import (
	"github.com/madalinpopa/go-bookreview/internal/testutil"
	"testing"
	"time"
)

// TestNotBlank verifies the behavior of the NotBlank function with various string inputs and their expected boolean results.
func TestNotBlank(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "empty string",
			input: "",
			want:  false,
		},
		{
			name:  "only whitespace",
			input: "   \t\n",
			want:  false,
		},
		{
			name:  "valid string",
			input: "hello",
			want:  true,
		},
		{
			name:  "string with spaces",
			input: "  hello  ",
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NotBlank(tt.input)
			testutil.Equal(t, got, tt.want)
		})
	}
}

// TestMaxChars validates the behavior of the MaxChars function with various input strings and character limits.
func TestMaxChars(t *testing.T) {
	tests := []struct {
		name  string
		input string
		limit int
		want  bool
	}{
		{
			name:  "empty string",
			input: "",
			limit: 5,
			want:  true,
		},
		{
			name:  "exactly at limit",
			input: "12345",
			limit: 5,
			want:  true,
		},
		{
			name:  "under limit",
			input: "123",
			limit: 5,
			want:  true,
		},
		{
			name:  "over limit",
			input: "123456",
			limit: 5,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaxChars(tt.input, tt.limit)
			testutil.Equal(t, got, tt.want)
		})
	}
}

// TestPermittedValue tests the PermittedValue function for both string and integer types with various input scenarios.
func TestPermittedValue(t *testing.T) {
	// String tests
	stringTests := []struct {
		name      string
		value     string
		permitted []string
		want      bool
	}{
		{
			name:      "string in permitted values",
			value:     "reading",
			permitted: []string{"want_to_read", "reading", "finished"},
			want:      true,
		},
		{
			name:      "string not in permitted values",
			value:     "unknown",
			permitted: []string{"want_to_read", "reading", "finished"},
			want:      false,
		},
		{
			name:      "empty permitted values",
			value:     "reading",
			permitted: []string{},
			want:      false,
		},
	}

	for _, tt := range stringTests {
		t.Run(tt.name, func(t *testing.T) {
			got := PermittedValue(tt.value, tt.permitted...)
			testutil.Equal(t, got, tt.want)
		})
	}

	// Integer tests
	intTests := []struct {
		name      string
		value     int
		permitted []int
		want      bool
	}{
		{
			name:      "int in permitted values",
			value:     1,
			permitted: []int{1, 2, 3},
			want:      true,
		},
		{
			name:      "int not in permitted values",
			value:     4,
			permitted: []int{1, 2, 3},
			want:      false,
		},
	}

	for _, tt := range intTests {
		t.Run(tt.name, func(t *testing.T) {
			got := PermittedValue(tt.value, tt.permitted...)
			testutil.Equal(t, got, tt.want)
		})
	}
}

// TestValidDate validates the behavior of the ValidDate function by checking various inputs and comparing them to expected results.
func TestValidDate(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want bool
	}{
		{
			name: "zero date",
			date: time.Time{}, // Zero value
			want: false,
		},
		{
			name: "valid date",
			date: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
			want: true,
		},
		{
			name: "current time",
			date: time.Now(),
			want: true,
		},
		{
			name: "past date",
			date: time.Date(2020, time.December, 25, 0, 0, 0, 0, time.UTC),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidDate(tt.date)
			testutil.Equal(t, got, tt.want)
		})
	}
}

// TestMinChars tests the MinChars function to ensure it correctly verifies if a string meets the minimum character requirement.
func TestMinChars(t *testing.T) {
	tests := []struct {
		name  string
		value string
		min   int
		want  bool
	}{
		{
			name:  "empty string",
			value: "",
			min:   1,
			want:  false,
		},
		{
			name:  "exactly minimum length",
			value: "12345",
			min:   5,
			want:  true,
		},
		{
			name:  "above minimum length",
			value: "123456",
			min:   5,
			want:  true,
		},
		{
			name:  "below minimum length",
			value: "1234",
			min:   5,
			want:  false,
		},
		{
			name:  "minimum zero",
			value: "",
			min:   0,
			want:  true,
		},
		{
			name:  "unicode characters",
			value: "Hello世界",
			min:   5,
			want:  true,
		},
		{
			name:  "spaces only",
			value: "     ",
			min:   3,
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MinChars(tt.value, tt.min)
			testutil.Equal(t, got, tt.want)
		})
	}
}

// TestMatches verifies the Matches function using a set of test cases for validating email address formats.
func TestMatches(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "valid email",
			value: "user@example.com",
			want:  true,
		},
		{
			name:  "email with subdomain",
			value: "user@sub.example.com",
			want:  true,
		},
		{
			name:  "email with plus",
			value: "user+tag@example.com",
			want:  true,
		},
		{
			name:  "missing @",
			value: "userexample.com",
			want:  false,
		},
		{
			name:  "missing domain",
			value: "user@",
			want:  false,
		},
		{
			name:  "missing local part",
			value: "@example.com",
			want:  false,
		},
		{
			name:  "empty string",
			value: "",
			want:  false,
		},
		{
			name:  "email with special characters",
			value: "user.name+tag@example-site.co.uk",
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Matches(tt.value, EmailRX)
			testutil.Equal(t, got, tt.want)
		})
	}
}

// TestValidNumber tests the ValidNumber function to ensure it correctly validates if a string represents a valid integer.
func TestValidNumber(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "valid integer",
			value: "123",
			want:  true,
		},
		{
			name:  "zero",
			value: "0",
			want:  true,
		},
		{
			name:  "negative number",
			value: "-123",
			want:  true,
		},
		{
			name:  "empty string",
			value: "",
			want:  false,
		},
		{
			name:  "decimal number",
			value: "123.45",
			want:  false,
		},
		{
			name:  "non-numeric string",
			value: "abc",
			want:  false,
		},
		{
			name:  "mixed string",
			value: "123abc",
			want:  false,
		},
		{
			name:  "spaces only",
			value: "   ",
			want:  false,
		},
		{
			name:  "number with spaces",
			value: "  123  ",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidNumber(tt.value)
			testutil.Equal(t, got, tt.want)
		})
	}
}

// TestMaxNumber verifies the behavior of the MaxNumber function with various input values and expected results.
func TestMaxNumber(t *testing.T) {
	tests := []struct {
		name  string
		value int
		max   int
		want  bool
	}{
		{
			name:  "below maximum",
			value: 3,
			max:   5,
			want:  true,
		},
		{
			name:  "equal to maximum",
			value: 5,
			max:   5,
			want:  true,
		},
		{
			name:  "above maximum",
			value: 7,
			max:   5,
			want:  false,
		},
		{
			name:  "zero value",
			value: 0,
			max:   5,
			want:  true,
		},
		{
			name:  "negative value",
			value: -5,
			max:   5,
			want:  true,
		},
		{
			name:  "zero maximum",
			value: 1,
			max:   0,
			want:  false,
		},
		{
			name:  "negative maximum",
			value: -10,
			max:   -5,
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaxNumber(tt.value, tt.max)
			testutil.Equal(t, got, tt.want)
		})
	}
}
