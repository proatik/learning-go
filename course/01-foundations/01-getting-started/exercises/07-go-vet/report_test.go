package report

import "testing"

func TestSummary(t *testing.T) {
	if got, want := Summary("apples", 3), "apples: 3 items"; got != want {
		t.Errorf("Summary = %q, want %q", got, want)
	}
}

func TestPercent(t *testing.T) {
	tests := []struct {
		part, total int
		want        string
	}{
		{1, 4, "25.0%"},
		{1, 3, "33.3%"},
		{2, 2, "100.0%"},
	}
	for _, tt := range tests {
		if got := Percent(tt.part, tt.total); got != tt.want {
			t.Errorf("Percent(%d, %d) = %q, want %q", tt.part, tt.total, got, tt.want)
		}
	}
}
