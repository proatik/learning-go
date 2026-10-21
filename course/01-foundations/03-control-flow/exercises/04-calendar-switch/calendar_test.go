package calendar

import "testing"

func TestDayType(t *testing.T) {
	tests := []struct{ day, want string }{
		{"Mon", "weekday"}, {"Tue", "weekday"}, {"Wed", "weekday"},
		{"Thu", "weekday"}, {"Fri", "weekday"},
		{"Sat", "weekend"}, {"Sun", "weekend"},
		{"sun", "invalid"}, {"MON", "invalid"}, {"Monday", "invalid"},
		{"", "invalid"}, {"Xyz", "invalid"},
	}
	for _, tt := range tests {
		if got := DayType(tt.day); got != tt.want {
			t.Errorf("DayType(%q) = %q, want %q", tt.day, got, tt.want)
		}
	}
}

func TestDaysInMonth(t *testing.T) {
	tests := []struct{ month, year, want int }{
		{1, 2023, 31}, {3, 2023, 31}, {5, 2023, 31}, {7, 2023, 31},
		{8, 2023, 31}, {10, 2023, 31}, {12, 2023, 31},
		{4, 2023, 30}, {6, 2023, 30}, {9, 2023, 30}, {11, 2023, 30},
		{2, 2023, 28}, {2, 2024, 29}, {2, 1900, 28}, {2, 2100, 28},
		{2, 2000, 29}, {2, 1600, 29},
		{0, 2024, 0}, {13, 2024, 0}, {-1, 2024, 0}, {100, 2024, 0},
	}
	for _, tt := range tests {
		if got := DaysInMonth(tt.month, tt.year); got != tt.want {
			t.Errorf("DaysInMonth(%d, %d) = %d, want %d", tt.month, tt.year, got, tt.want)
		}
	}
}
