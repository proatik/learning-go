package bounds

import "testing"

func TestMinMaxSingle(t *testing.T) {
	for _, n := range []int{7, 0, -12} {
		lo, hi := MinMax(n)
		if lo != n || hi != n {
			t.Errorf("MinMax(%d) = (%d, %d), want (%d, %d)", n, lo, hi, n, n)
		}
	}
}

func TestMinMax(t *testing.T) {
	tests := []struct {
		first  int
		rest   []int
		lo, hi int
	}{
		{3, []int{9, -2, 5}, -2, 9},
		{4, []int{4, 4}, 4, 4},
		{1, []int{2, 3, 4, 5}, 1, 5},
		{5, []int{4, 3, 2, 1}, 1, 5},
		{0, []int{-1}, -1, 0},
		{-5, []int{-9, -1, -7}, -9, -1},
		{10, []int{20, 10, 30, 0, 15}, 0, 30},
		{100, []int{}, 100, 100},
	}
	for _, tt := range tests {
		lo, hi := MinMax(tt.first, tt.rest...)
		if lo != tt.lo || hi != tt.hi {
			t.Errorf("MinMax(%d, %v...) = (%d, %d), want (%d, %d)", tt.first, tt.rest, lo, hi, tt.lo, tt.hi)
		}
	}
}
