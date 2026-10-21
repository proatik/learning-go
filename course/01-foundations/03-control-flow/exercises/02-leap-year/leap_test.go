package leap

import "testing"

var cases = []struct {
	year int
	leap bool
}{
	{2024, true},
	{2023, false},
	{2019, false},
	{1996, true},
	{1900, false},
	{2100, false},
	{1800, false},
	{2000, true},
	{1600, true},
	{2400, true},
	{4, true},
	{1, false},
	{0, true}, // year 0 is divisible by 400
}

func TestIsLeap(t *testing.T) {
	for _, tt := range cases {
		if got := IsLeap(tt.year); got != tt.leap {
			t.Errorf("IsLeap(%d) = %t, want %t", tt.year, got, tt.leap)
		}
	}
}

func TestDaysInYear(t *testing.T) {
	for _, tt := range cases {
		want := 365
		if tt.leap {
			want = 366
		}
		if got := DaysInYear(tt.year); got != want {
			t.Errorf("DaysInYear(%d) = %d, want %d", tt.year, got, want)
		}
	}
}
