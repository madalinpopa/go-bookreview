package testutil

import (
	"errors"
	"testing"
)

// TestEqual verifies the Equal function for various test cases with different types to ensure proper equality comparison.
func TestEqual(t *testing.T) {
	tests := []struct {
		name string
		got  interface{}
		want interface{}
	}{
		{
			name: "strings are equal",
			got:  "hello",
			want: "hello",
		},
		{
			name: "integers are equal",
			got:  42,
			want: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &testing.T{}
			Equal(mockT, tt.got, tt.want)
			if mockT.Failed() {
				t.Errorf("Equal() failed for equal values")
			}
		})
	}
}

func TestNoError(t *testing.T) {
	mockT := &testing.T{}
	NoError(mockT, nil)
	if mockT.Failed() {
		t.Error("NoError() failed for nil error")
	}

	mockT = &testing.T{}
	NoError(mockT, errors.New("test error"))
	if !mockT.Failed() {
		t.Error("NoError() didn't fail for non-nil error")
	}
}
