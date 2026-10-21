package clock

import "testing"

func TestAddMinutes(t *testing.T) {
	tests := []struct {
		h, m, delta int
		want        string
	}{
		{10, 0, 3, "10:03"},
		{0, 0, 0, "00:00"},
		{9, 5, 0, "09:05"},
		{23, 59, 1, "00:00"},
		{23, 59, 2, "00:01"},
		{0, 0, -1, "23:59"},
		{0, 0, -1441, "23:59"},
		{12, 30, 24 * 60, "12:30"},
		{1, 0, -3 * 24 * 60, "01:00"},
		{6, 15, 1000, "22:55"},
		{6, 15, -1000, "13:35"},
		{25, 160, 0, "03:40"},
		{-1, 0, 0, "23:00"},
		{0, -1, 0, "23:59"},
		{-25, -61, -1, "21:58"},
		{0, 0, 1_000_000, "10:40"},
	}
	for _, tt := range tests {
		if got := AddMinutes(tt.h, tt.m, tt.delta); got != tt.want {
			t.Errorf("AddMinutes(%d, %d, %d) = %q, want %q", tt.h, tt.m, tt.delta, got, tt.want)
		}
	}
}

func TestParseClock(t *testing.T) {
	valid := []struct {
		s    string
		h, m int
	}{
		{"00:00", 0, 0},
		{"09:05", 9, 5},
		{"12:30", 12, 30},
		{"23:59", 23, 59},
		{"10:00", 10, 0},
	}
	for _, tt := range valid {
		h, m, err := ParseClock(tt.s)
		if err != nil || h != tt.h || m != tt.m {
			t.Errorf("ParseClock(%q) = %d, %d, %v, want %d, %d, nil", tt.s, h, m, err, tt.h, tt.m)
		}
	}

	invalid := []string{
		"", "9:05", "09:5", "0905", "09-05", "09:055", "24:00", "23:60", "99:99",
		"+1:00", " 9:00", "09: 5", "ab:cd", "-1:00", "09:-5", "١٢:٣٠",
	}
	for _, s := range invalid {
		h, m, err := ParseClock(s)
		if err == nil {
			t.Errorf("ParseClock(%q) = %d, %d, nil, want an error", s, h, m)
		} else if h != 0 || m != 0 {
			t.Errorf("ParseClock(%q) = %d, %d, %v, want 0, 0 with the error", s, h, m, err)
		}
	}
}
