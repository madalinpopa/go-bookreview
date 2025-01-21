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
