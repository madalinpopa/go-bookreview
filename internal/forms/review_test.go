package forms

import (
	"github.com/madalinpopa/go-bookreview/internal/testutil"
	"testing"
)

// TestBookReviewForm_Validate tests the validation logic for BookReviewForm, ensuring correct validation of rating and review text.
func TestBookReviewForm_Validate(t *testing.T) {
	tests := []struct {
		name          string
		form          BookReviewForm
		wantValid     bool
		wantFieldErrs map[string]string
	}{
		{
			name: "valid review",
			form: BookReviewForm{
				Rating:     4,
				ReviewText: "Great book, highly recommended!",
			},
			wantValid:     true,
			wantFieldErrs: nil,
		},
		{
			name: "rating too high",
			form: BookReviewForm{
				Rating:     6,
				ReviewText: "Great book!",
			},
			wantValid: false,
			wantFieldErrs: map[string]string{
				"rating": "Rating must be a number between 1 and 5",
			},
		},
		{
			name: "empty review text",
			form: BookReviewForm{
				Rating:     4,
				ReviewText: "",
			},
			wantValid: false,
			wantFieldErrs: map[string]string{
				"review_text": "Review text is required",
			},
		},
		{
			name: "whitespace only review",
			form: BookReviewForm{
				Rating:     4,
				ReviewText: "    ",
			},
			wantValid: false,
			wantFieldErrs: map[string]string{
				"review_text": "Review text is required",
			},
		},
		{
			name: "multiple errors",
			form: BookReviewForm{
				Rating:     6,
				ReviewText: "",
			},
			wantValid: false,
			wantFieldErrs: map[string]string{
				"rating":      "Rating must be a number between 1 and 5",
				"review_text": "Review text is required",
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
