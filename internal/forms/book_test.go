package forms

import (
	"github.com/madalinpopa/go-bookreview/internal/testutil"
	"testing"
)

// TestBookForm_Validate tests the validation logic of the BookForm struct, ensuring correct field validations and error handling.
func TestBookForm_Validate(t *testing.T) {
	tests := []struct {
		name          string
		form          BookForm
		wantValid     bool
		wantFieldErrs map[string]string
	}{
		{
			name: "valid form",
			form: BookForm{
				Title:           "The Go Programming Language",
				Author:          "Alan A. A. Donovan",
				ISBN:            "978-0134190440",
				PublicationYear: 2015,
			},
			wantValid:     true,
			wantFieldErrs: nil,
		},
		{
			name: "empty title",
			form: BookForm{
				Title:           "",
				Author:          "Alan A. A. Donovan",
				ISBN:            "978-0134190440",
				PublicationYear: 2015,
			},
			wantValid: false,
			wantFieldErrs: map[string]string{
				"title": "Title is required",
			},
		},
		{
			name: "empty author",
			form: BookForm{
				Title:           "The Go Programming Language",
				Author:          "",
				ISBN:            "978-0134190440",
				PublicationYear: 2015,
			},
			wantValid: false,
			wantFieldErrs: map[string]string{
				"author": "Author is required",
			},
		},
		{
			name: "empty isbn",
			form: BookForm{
				Title:           "The Go Programming Language",
				Author:          "Alan A. A. Donovan",
				ISBN:            "",
				PublicationYear: 2015,
			},
			wantValid: false,
			wantFieldErrs: map[string]string{
				"isbn": "ISBN is required",
			},
		},
		{
			name: "all fields empty",
			form: BookForm{
				Title:           "",
				Author:          "",
				ISBN:            "",
				PublicationYear: 0,
			},
			wantValid: false,
			wantFieldErrs: map[string]string{
				"title":  "Title is required",
				"author": "Author is required",
				"isbn":   "ISBN is required",
			},
		},
		{
			name: "whitespace fields",
			form: BookForm{
				Title:           "   ",
				Author:          "   ",
				ISBN:            "   ",
				PublicationYear: 0,
			},
			wantValid: false,
			wantFieldErrs: map[string]string{
				"title":  "Title is required",
				"author": "Author is required",
				"isbn":   "ISBN is required",
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
