package forms

import (
	"github.com/madalinpopa/go-bookreview/internal/testutil"
	"testing"
)

// TestBase_Valid verifies the Valid method of the Base struct by testing various combinations of field and non-field errors.
func TestBase_Valid(t *testing.T) {
	tests := []struct {
		name           string
		fieldErrors    map[string]string
		nonFieldErrors []string
		want           bool
	}{
		{
			name:           "no errors",
			fieldErrors:    nil,
			nonFieldErrors: nil,
			want:           true,
		},
		{
			name: "with field errors",
			fieldErrors: map[string]string{
				"username": "This field is required",
			},
			nonFieldErrors: nil,
			want:           false,
		},
		{
			name:           "with non-field errors",
			fieldErrors:    nil,
			nonFieldErrors: []string{"Invalid credentials"},
			want:           false,
		},
		{
			name: "with both error types",
			fieldErrors: map[string]string{
				"email": "Invalid email format",
			},
			nonFieldErrors: []string{"Server error"},
			want:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := &Base{
				FieldErrors:    tt.fieldErrors,
				NonFieldErrors: tt.nonFieldErrors,
			}
			got := form.Valid()
			testutil.Equal(t, got, tt.want)
		})
	}
}

// TestBase_AddFieldError tests the AddFieldError method to ensure it correctly handles field error insertion scenarios.
func TestBase_AddFieldError(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		message  string
		existing map[string]string
		want     map[string]string
	}{
		{
			name:     "add to empty errors",
			key:      "username",
			message:  "This field is required",
			existing: nil,
			want: map[string]string{
				"username": "This field is required",
			},
		},
		{
			name:    "add to existing errors",
			key:     "email",
			message: "Invalid email",
			existing: map[string]string{
				"username": "This field is required",
			},
			want: map[string]string{
				"username": "This field is required",
				"email":    "Invalid email",
			},
		},
		{
			name:    "don't override existing key",
			key:     "username",
			message: "New error",
			existing: map[string]string{
				"username": "Original error",
			},
			want: map[string]string{
				"username": "Original error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := &Base{
				FieldErrors: tt.existing,
			}
			form.AddFieldError(tt.key, tt.message)
			testutil.Equal(t, len(form.FieldErrors), len(tt.want))
			for k, v := range tt.want {
				got, exists := form.FieldErrors[k]
				testutil.Equal(t, exists, true)
				testutil.Equal(t, got, v)
			}
		})
	}
}

// TestBase_AddNonFieldError tests the AddNonFieldError method to ensure it appends messages to NonFieldErrors correctly.
func TestBase_AddNonFieldError(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		existing []string
		want     []string
	}{
		{
			name:     "add to empty errors",
			message:  "Server error",
			existing: nil,
			want:     []string{"Server error"},
		},
		{
			name:     "add to existing errors",
			message:  "Database error",
			existing: []string{"Server error"},
			want:     []string{"Server error", "Database error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := &Base{
				NonFieldErrors: tt.existing,
			}
			form.AddNonFieldError(tt.message)
			testutil.Equal(t, len(form.NonFieldErrors), len(tt.want))
			for i, v := range tt.want {
				testutil.Equal(t, form.NonFieldErrors[i], v)
			}
		})
	}
}

// TestBase_CheckField tests the CheckField method of the Base struct to ensure it handles field validation correctly.
func TestBase_CheckField(t *testing.T) {
	tests := []struct {
		name    string
		ok      bool
		key     string
		message string
		want    bool
	}{
		{
			name:    "valid field",
			ok:      true,
			key:     "username",
			message: "This field is required",
			want:    false, // no error should be added
		},
		{
			name:    "invalid field",
			ok:      false,
			key:     "email",
			message: "Invalid email format",
			want:    true, // error should be added
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := &Base{}
			form.CheckField(tt.ok, tt.key, tt.message)
			_, got := form.FieldErrors[tt.key]
			testutil.Equal(t, got, tt.want)
			if tt.want {
				testutil.Equal(t, form.FieldErrors[tt.key], tt.message)
			}
		})
	}
}
