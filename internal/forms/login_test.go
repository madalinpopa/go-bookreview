package forms

import (
	"github.com/madalinpopa/go-bookreview/internal/testutil"
	"testing"
)

// TestUserLoginForm_Validate tests the validation logic for UserLoginForm, ensuring proper handling of various input cases.
func TestUserLoginForm_Validate(t *testing.T) {
	tests := []struct {
		name          string
		form          UserLoginForm
		wantValid     bool
		wantFieldErrs map[string]string
	}{
		{
			name: "valid form",
			form: UserLoginForm{
				Username: "testuser",
				Password: "password123",
			},
			wantValid:     true,
			wantFieldErrs: nil,
		},
		{
			name: "empty username",
			form: UserLoginForm{
				Username: "",
				Password: "password123",
			},
			wantValid: false,
			wantFieldErrs: map[string]string{
				"username": "This field is required",
			},
		},
		{
			name: "empty password",
			form: UserLoginForm{
				Username: "testuser",
				Password: "",
			},
			wantValid: false,
			wantFieldErrs: map[string]string{
				"password": "This field is required",
			},
		},
		{
			name: "both fields empty",
			form: UserLoginForm{
				Username: "",
				Password: "",
			},
			wantValid: false,
			wantFieldErrs: map[string]string{
				"username": "This field is required",
				"password": "This field is required",
			},
		},
		{
			name: "whitespace only",
			form: UserLoginForm{
				Username: "   ",
				Password: "   ",
			},
			wantValid: false,
			wantFieldErrs: map[string]string{
				"username": "This field is required",
				"password": "This field is required",
			},
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
