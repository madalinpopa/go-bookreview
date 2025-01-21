package testutil

import "testing"

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
			// Using a testing.T wrapper to catch the error calls
			mockT := &testing.T{}
			Equal(mockT, tt.got, tt.want)
			if mockT.Failed() {
				t.Errorf("Equal() failed for equal values")
			}
		})
	}
}
