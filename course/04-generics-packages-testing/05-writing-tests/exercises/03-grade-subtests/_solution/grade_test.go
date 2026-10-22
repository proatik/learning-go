package grade

import "testing"

// TestGrade checks Grade, one named subtest per case.
func TestGrade(t *testing.T) {
	tests := []struct {
		name  string
		score int
		want  string
	}{
		{"perfect", 100, "A"},
		{"lowest A", 90, "A"},
		{"highest B", 89, "B"},
		{"lowest B", 80, "B"},
		{"lowest C", 70, "C"},
		{"lowest D", 60, "D"},
		{"highest F", 59, "F"},
		{"zero", 0, "F"},
		{"above 100", 101, "invalid"},
		{"negative", -1, "invalid"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Grade(tt.score); got != tt.want {
				t.Errorf("Grade(%d) = %q, want %q", tt.score, got, tt.want)
			}
		})
	}
}
