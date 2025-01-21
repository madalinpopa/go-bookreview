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
