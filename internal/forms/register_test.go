package forms

import (
	"github.com/madalinpopa/go-bookreview/internal/testutil"
	"testing"
)

// TestRegisterForm_Validate is a unit test for the RegisterForm.Validate method ensuring validation logic correctness.
func TestRegisterForm_Validate(t *testing.T) {
	tests := []struct {
		name          string
		form          RegisterForm
		wantValid     bool
		wantFieldErrs map[string]string
	}{
		{
			name: "valid form",
			form: RegisterForm{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantValid:     true,
			wantFieldErrs: nil,
		},
		{
			name: "empty username",
			form: RegisterForm{
				Username: "",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantValid: false,
			wantFieldErrs: map[string]string{
				"username": "This field is required",
			},
		},
		{
			name: "empty email",
			form: RegisterForm{
				Username: "testuser",
				Email:    "",
				Password: "password123",
			},
			wantValid: false,
			wantFieldErrs: map[string]string{
				"email": "This field is required",
			},
		},
		{
			name: "invalid email format",
			form: RegisterForm{
				Username: "testuser",
				Email:    "invalid-email",
				Password: "password123",
			},
			wantValid: false,
			wantFieldErrs: map[string]string{
				"email": "The email address is not valid.",
			},
		},
		{
			name: "password too short",
			form: RegisterForm{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "short",
			},
			wantValid: false,
			wantFieldErrs: map[string]string{
				"password": "Password must be at least 8 characters.",
			},
		},
		{
			name: "all fields empty",
			form: RegisterForm{
				Username: "",
				Email:    "",
				Password: "",
			},
			wantValid: false,
			wantFieldErrs: map[string]string{
				"username": "This field is required",
				"email":    "This field is required",
				"password": "This field is required",
			},
		},
		{
			name: "valid complex email",
			form: RegisterForm{
				Username: "testuser",
				Email:    "user.name+tag@example-site.co.uk",
				Password: "password123",
			},
			wantValid:     true,
			wantFieldErrs: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.form.Validate()
			testutil.Equal(t, tt.form.Valid(), tt.wantValid)

			// Check field errors
			if tt.wantFieldErrs == nil {
				testutil.Equal(t, len(tt.form.FieldErrors), 0)
			} else {
				testutil.Equal(t, len(tt.form.FieldErrors), len(tt.wantFieldErrs))
				for k, want := range tt.wantFieldErrs {
					got, exists := tt.form.FieldErrors[k]
					testutil.Equal(t, exists, true)
					testutil.Equal(t, got, want)
				}
			}
		})
	}
}
