package grade

import "testing"

func TestGrade(t *testing.T) {
	tests := []struct {
		score int
		want  string
	}{
		{100, "A"}, {95, "A"}, {90, "A"},
		{89, "B"}, {80, "B"},
		{79, "C"}, {75, "C"}, {70, "C"},
		{69, "D"}, {60, "D"},
		{59, "F"}, {30, "F"}, {1, "F"}, {0, "F"},
		{101, "invalid"}, {1000, "invalid"}, {-1, "invalid"}, {-50, "invalid"},
	}
	for _, tt := range tests {
		if got := Grade(tt.score); got != tt.want {
			t.Errorf("Grade(%d) = %q, want %q", tt.score, got, tt.want)
		}
	}
}
