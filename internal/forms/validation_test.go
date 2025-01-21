package forms

import (
	"github.com/madalinpopa/go-bookreview/internal/testutil"
	"testing"
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

// TestPermittedValue verifies the behavior of the PermittedValue function by testing various scenarios and expected results.
func TestPermittedValue(t *testing.T) {
	tests := []struct {
		name      string
		value     interface{}
		permitted []interface{}
		want      bool
	}{
		{
			name:      "string in permitted values",
			value:     "reading",
			permitted: []interface{}{"want_to_read", "reading", "finished"},
			want:      true,
		},
		{
			name:      "string not in permitted values",
			value:     "unknown",
			permitted: []interface{}{"want_to_read", "reading", "finished"},
			want:      false,
		},
		{
			name:      "empty permitted values",
			value:     "reading",
			permitted: []interface{}{},
			want:      false,
		},
		{
			name:      "int in permitted values",
			value:     1,
			permitted: []interface{}{1, 2, 3},
			want:      true,
		},
		{
			name:      "int not in permitted values",
			value:     4,
			permitted: []interface{}{1, 2, 3},
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PermittedValue(tt.value, tt.permitted...)
			testutil.Equal(t, got, tt.want)
		})
	}
}
