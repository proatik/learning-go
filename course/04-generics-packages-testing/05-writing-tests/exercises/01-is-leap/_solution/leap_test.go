package leap

import "testing"

// TestIsLeap checks IsLeap with a table of years and the expected answers.
func TestIsLeap(t *testing.T) {
	tests := []struct {
		year int
		want bool
	}{
		{2023, false}, // not divisible by 4
		{2024, true},  // divisible by 4
		{1900, false}, // a century
		{2100, false}, // another century
		{2000, true},  // a 400th year
		{1600, true},
	}
	for _, tt := range tests {
		if got := IsLeap(tt.year); got != tt.want {
			t.Errorf("IsLeap(%d) = %t, want %t", tt.year, got, tt.want)
		}
	}
}
