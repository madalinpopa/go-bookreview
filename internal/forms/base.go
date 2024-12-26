package forms

import "errors"

var (

	// ErrInvalidFileType represents an error indicating that the provided file type is not supported or invalid.
	ErrInvalidFileType = errors.New("invalid file type")

	// ErrFormBadRequest represents an error indicating a bad request occurred while processing a form.
	ErrFormBadRequest = errors.New("bad request while processing form")
)

// Former defines an interface for input validation, tracking field errors, and managing form-related error states.
type Former interface {
	Valid() bool
	AddFieldError(key, message string)
	AddNonFieldError(message string)
	CheckField(ok bool, key, message string)
}

// Base represents a structure for managing field-specific and non-field-specific validation errors.
// NonFieldErrors hold validation error messages not tied to specific fields.
// FieldErrors maps field names to their corresponding validation error messages.
type Base struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

// Valid checks if the validator's FieldErrors map is empty, indicating that no validation errors are present.
func (b *Base) Valid() bool {
	return len(b.FieldErrors) == 0 && len(b.NonFieldErrors) == 0
}

// AddFieldError adds a validation error message for a specific field if it does not already exist in the FieldErrors map.
func (b *Base) AddFieldError(key, message string) {

	// Check if the FieldErrors map is nil.
	// If it is nil, initialize it to a new empty map to avoid runtime errors when adding values.
	if b.FieldErrors == nil {
		b.FieldErrors = make(map[string]string)
	}

	if _, exists := b.FieldErrors[key]; !exists {
		b.FieldErrors[key] = message
	}
}

// AddNonFieldError appends a non-field-specific error message to the NonFieldErrors slice.
func (b *Base) AddNonFieldError(message string) {
	b.NonFieldErrors = append(b.NonFieldErrors, message)
}

// CheckField adds a validation error for the specified
// field if the condition is not met.
func (b *Base) CheckField(ok bool, key, message string) {
	if !ok {
		b.AddFieldError(key, message)
	}
}
