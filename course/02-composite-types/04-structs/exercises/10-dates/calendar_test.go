package calendar

import (
	"slices"
	"testing"
	"time"
)

// viaTime computes the expected result of AddDays with the standard library.
func viaTime(d Date, n int) Date {
	t := time.Date(d.Year, time.Month(d.Month), d.Day, 12, 0, 0, 0, time.UTC).AddDate(0, 0, n)
	return Date{Year: t.Year(), Month: int(t.Month()), Day: t.Day()}
}

func TestIsLeap(t *testing.T) {
	tests := []struct {
		year int
		want bool
	}{
		{2024, true}, {2023, false}, {1900, false}, {2000, true},
		{2100, false}, {2400, true}, {4, true}, {1, false}, {1600, true}, {1700, false},
	}
	for _, tt := range tests {
		if got := IsLeap(tt.year); got != tt.want {
			t.Errorf("IsLeap(%d) = %v, want %v", tt.year, got, tt.want)
		}
	}
}

func TestDaysInMonth(t *testing.T) {
	want := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	for m := 1; m <= 12; m++ {
		if got := DaysInMonth(2023, m); got != want[m-1] {
			t.Errorf("DaysInMonth(2023, %d) = %d, want %d", m, got, want[m-1])
		}
	}
	tests := []struct{ year, month, want int }{
		{2024, 2, 29}, {1900, 2, 28}, {2000, 2, 29},
		{2024, 0, 0}, {2024, 13, 0}, {2024, -1, 0},
	}
	for _, tt := range tests {
		if got := DaysInMonth(tt.year, tt.month); got != tt.want {
			t.Errorf("DaysInMonth(%d, %d) = %d, want %d", tt.year, tt.month, got, tt.want)
		}
	}
}

func TestValid(t *testing.T) {
	tests := []struct {
		d    Date
		want bool
	}{
		{Date{2024, 2, 29}, true},
		{Date{2023, 2, 29}, false},
		{Date{2023, 12, 31}, true},
		{Date{2023, 4, 31}, false},
		{Date{2023, 1, 0}, false},
		{Date{2023, 0, 10}, false},
		{Date{2023, 13, 1}, false},
		{Date{1, 1, 1}, true},
		{Date{9999, 12, 31}, true},
		{Date{0, 1, 1}, false},
		{Date{10000, 1, 1}, false},
		{Date{}, false},
	}
	for _, tt := range tests {
		if got := Valid(tt.d); got != tt.want {
			t.Errorf("Valid(%+v) = %v, want %v", tt.d, got, tt.want)
		}
	}
}

func TestCompareAndIsBefore(t *testing.T) {
	tests := []struct {
		a, b Date
		want int
	}{
		{Date{2024, 1, 31}, Date{2024, 2, 1}, -1},
		{Date{2024, 2, 1}, Date{2024, 1, 31}, 1},
		{Date{2023, 12, 31}, Date{2024, 1, 1}, -1},
		{Date{2024, 5, 5}, Date{2024, 5, 5}, 0},
		{Date{2024, 5, 4}, Date{2024, 5, 5}, -1},
		{Date{2025, 1, 1}, Date{2024, 12, 31}, 1},
		{Date{1, 1, 1}, Date{9999, 12, 31}, -1},
	}
	for _, tt := range tests {
		if got := Compare(tt.a, tt.b); got != tt.want {
			t.Errorf("Compare(%+v, %+v) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
		if got := IsBefore(tt.a, tt.b); got != (tt.want < 0) {
			t.Errorf("IsBefore(%+v, %+v) = %v, want %v", tt.a, tt.b, got, tt.want < 0)
		}
	}
}

func TestSortWithCompare(t *testing.T) {
	dates := []Date{{2024, 3, 1}, {2023, 12, 31}, {2024, 2, 29}, {2024, 3, 1}, {1999, 1, 1}}
	slices.SortFunc(dates, Compare)
	want := []Date{{1999, 1, 1}, {2023, 12, 31}, {2024, 2, 29}, {2024, 3, 1}, {2024, 3, 1}}
	if !slices.Equal(dates, want) {
		t.Errorf("slices.SortFunc(dates, Compare) = %+v, want %+v", dates, want)
	}
}

func TestAddDays(t *testing.T) {
	tests := []struct {
		d    Date
		n    int
		want Date
	}{
		{Date{2024, 2, 28}, 1, Date{2024, 2, 29}},
		{Date{2023, 2, 28}, 1, Date{2023, 3, 1}},
		{Date{2024, 12, 31}, 1, Date{2025, 1, 1}},
		{Date{2024, 3, 1}, -1, Date{2024, 2, 29}},
		{Date{2024, 1, 1}, -1, Date{2023, 12, 31}},
		{Date{2024, 1, 1}, 366, Date{2025, 1, 1}},
		{Date{2024, 5, 17}, 0, Date{2024, 5, 17}},
		{Date{1, 1, 1}, 0, Date{1, 1, 1}},
		{Date{1, 1, 2}, -1, Date{1, 1, 1}},
		{Date{9999, 12, 30}, 1, Date{9999, 12, 31}},
		{Date{1900, 2, 28}, 1, Date{1900, 3, 1}},
		{Date{2000, 2, 28}, 1, Date{2000, 2, 29}},
	}
	for _, tt := range tests {
		if got := AddDays(tt.d, tt.n); got != tt.want {
			t.Errorf("AddDays(%+v, %d) = %+v, want %+v", tt.d, tt.n, got, tt.want)
		}
	}
}

func TestAddDaysAgainstTime(t *testing.T) {
	starts := []Date{{2024, 2, 29}, {2000, 1, 1}, {1970, 7, 15}, {1600, 12, 31}, {2100, 3, 1}, {5000, 6, 30}}
	steps := []int{1, -1, 30, -30, 365, -365, 1461, 10_000, -10_000, 146_097, -146_097, 500_000, -500_000, 534_567}
	for _, d := range starts {
		for _, n := range steps {
			want := viaTime(d, n)
			if got := AddDays(d, n); got != want {
				t.Errorf("AddDays(%+v, %d) = %+v, want %+v", d, n, got, want)
			}
		}
	}
	// Every day of a leap and a non-leap year.
	for _, start := range []Date{{2023, 12, 31}, {2024, 1, 1}} {
		for n := range 400 {
			if got, want := AddDays(start, n), viaTime(start, n); got != want {
				t.Fatalf("AddDays(%+v, %d) = %+v, want %+v", start, n, got, want)
			}
		}
	}
}

func TestDaysBetween(t *testing.T) {
	tests := []struct {
		a, b Date
		want int
	}{
		{Date{2024, 1, 1}, Date{2025, 1, 1}, 366},
		{Date{2023, 1, 1}, Date{2024, 1, 1}, 365},
		{Date{2024, 3, 1}, Date{2024, 2, 1}, -29},
		{Date{2024, 5, 5}, Date{2024, 5, 5}, 0},
		{Date{1, 1, 1}, Date{401, 1, 1}, 146_097},
		{Date{1970, 1, 1}, Date{2000, 1, 1}, 10_957},
		{Date{2000, 1, 1}, Date{1970, 1, 1}, -10_957},
		{Date{1, 1, 1}, Date{9999, 12, 31}, 3_652_058},
	}
	for _, tt := range tests {
		if got := DaysBetween(tt.a, tt.b); got != tt.want {
			t.Errorf("DaysBetween(%+v, %+v) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
	d := Date{1987, 6, 5}
	for _, n := range []int{-200_000, -1, 1, 59, 400_000} {
		if got := DaysBetween(d, viaTime(d, n)); got != n {
			t.Errorf("DaysBetween(%+v, %+v) = %d, want %d", d, viaTime(d, n), got, n)
		}
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		d    Date
		want string
	}{
		{Date{2024, 2, 29}, "2024-02-29"},
		{Date{987, 6, 5}, "0987-06-05"},
		{Date{1, 1, 1}, "0001-01-01"},
		{Date{9999, 12, 31}, "9999-12-31"},
		{Date{2024, 10, 20}, "2024-10-20"},
	}
	for _, tt := range tests {
		if got := Format(tt.d); got != tt.want {
			t.Errorf("Format(%+v) = %q, want %q", tt.d, got, tt.want)
		}
	}
}
