package forms

import (
	"github.com/madalinpopa/go-bookreview/internal/testutil"
	"testing"
)

// TestBookNoteForm_Validate tests the validation logic of the BookNoteForm, ensuring appropriate field errors and validity status.
func TestBookNoteForm_Validate(t *testing.T) {
	tests := []struct {
		name          string
		form          BookNoteForm
		wantValid     bool
		wantFieldErrs map[string]string
	}{
		{
			name: "valid note",
			form: BookNoteForm{
				NoteText:   "This chapter explains concurrency patterns well.",
				PageNumber: 42,
			},
			wantValid:     true,
			wantFieldErrs: nil,
		},
		{
			name: "empty note text",
			form: BookNoteForm{
				NoteText:   "",
				PageNumber: 42,
			},
			wantValid: false,
			wantFieldErrs: map[string]string{
				"note_text": "Note text is required",
			},
		},
		{
			name: "whitespace only note",
			form: BookNoteForm{
				NoteText:   "    ",
				PageNumber: 42,
			},
			wantValid: false,
			wantFieldErrs: map[string]string{
				"note_text": "Note text is required",
			},
		},
		{
			name: "note with no page number",
			form: BookNoteForm{
				NoteText:   "General note about the book.",
				PageNumber: 0,
			},
			wantValid:     true, // page number is optional
			wantFieldErrs: nil,
		},
		{
			name: "negative page number",
			form: BookNoteForm{
				NoteText:   "Valid note text",
				PageNumber: -1,
			},
			wantValid:     true, // no validation on page number range
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
