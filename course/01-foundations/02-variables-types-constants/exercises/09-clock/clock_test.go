package clock

import "testing"

func TestSplit(t *testing.T) {
	tests := []struct{ total, h, m, s int }{
		{0, 0, 0, 0},
		{59, 0, 0, 59},
		{60, 0, 1, 0},
		{3599, 0, 59, 59},
		{3600, 1, 0, 0},
		{3661, 1, 1, 1},
		{86399, 23, 59, 59},
		{90061, 25, 1, 1},
		{-59, 0, 0, -59},
		{-3661, -1, -1, -1},
		{-7200, -2, 0, 0},
	}
	for _, tt := range tests {
		h, m, s := Split(tt.total)
		if h != tt.h || m != tt.m || s != tt.s {
			t.Errorf("Split(%d) = %d, %d, %d, want %d, %d, %d", tt.total, h, m, s, tt.h, tt.m, tt.s)
		}
	}
}

func TestClock(t *testing.T) {
	tests := []struct {
		total int
		want  string
	}{
		{0, "00:00:00"},
		{5, "00:00:05"},
		{3661, "01:01:01"},
		{45296, "12:34:56"},
		{86399, "23:59:59"},
		{90061, "25:01:01"},
		{360000, "100:00:00"},
		{-90, "-00:01:30"},
		{-3661, "-01:01:01"},
	}
	for _, tt := range tests {
		if got := Clock(tt.total); got != tt.want {
			t.Errorf("Clock(%d) = %q, want %q", tt.total, got, tt.want)
		}
	}
}

func TestTimeOfDay(t *testing.T) {
	tests := []struct {
		total int
		want  string
	}{
		{0, "00:00:00"},
		{3600, "01:00:00"},
		{45296, "12:34:56"},
		{86399, "23:59:59"},
		{86400, "00:00:00"},
		{86400 + 61, "00:01:01"},
		{3*86400 + 7200, "02:00:00"},
		{-1, "23:59:59"},
		{-60, "23:59:00"},
		{-86400, "00:00:00"},
		{-86400 - 3600, "23:00:00"},
	}
	for _, tt := range tests {
		if got := TimeOfDay(tt.total); got != tt.want {
			t.Errorf("TimeOfDay(%d) = %q, want %q", tt.total, got, tt.want)
		}
	}
}
