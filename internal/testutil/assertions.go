package testutil

import "testing"

// Equal checks if the `got` and `want` values of a comparable type are equal, and reports an error if they differ.
func Equal[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v; want %v", got, want)
	}
}

// NoError checks if the given error is nil, and reports a test failure if it is not nil.
func NoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {

		t.Errorf("unexpected error: %v", err)
	}
}

// Error fails the test if the provided error is nil, optionally logging a custom message with additional arguments.
func Error(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Errorf("expected error but got nil")
	}
}
