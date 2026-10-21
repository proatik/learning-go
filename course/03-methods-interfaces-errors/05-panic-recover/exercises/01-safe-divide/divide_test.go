package divide

import (
	"math"
	"testing"
)

// noPanic runs f and reports a panic that escaped instead of crashing.
func noPanic(t *testing.T, name string, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("%s panicked: %v", name, r)
		}
	}()
	f()
}

func TestSafeDivide(t *testing.T) {
	tests := []struct {
		a, b int
		want int
	}{
		{7, 2, 3},
		{-7, 2, -3},
		{7, -2, -3},
		{0, 5, 0},
		{10, 10, 1},
		{1, 3, 0},
		{math.MinInt, -1, math.MinInt},
		{math.MaxInt, 1, math.MaxInt},
	}
	for _, tt := range tests {
		noPanic(t, "SafeDivide", func() {
			got, err := SafeDivide(tt.a, tt.b)
			if got != tt.want || err != nil {
				t.Errorf("SafeDivide(%d, %d) = %d, %v; want %d, <nil>", tt.a, tt.b, got, err, tt.want)
			}
		})
	}
}

func TestSafeDivideByZero(t *testing.T) {
	const want = "safe divide: runtime error: integer divide by zero"
	for _, a := range []int{5, 0, -3, math.MaxInt} {
		noPanic(t, "SafeDivide", func() {
			got, err := SafeDivide(a, 0)
			if err == nil {
				t.Errorf("SafeDivide(%d, 0) error = <nil>, want %q", a, want)
				return
			}
			if err.Error() != want {
				t.Errorf("SafeDivide(%d, 0) error = %q, want %q", a, err.Error(), want)
			}
			if got != 0 {
				t.Errorf("SafeDivide(%d, 0) = %d, want 0", a, got)
			}
		})
	}
}

func TestAverage(t *testing.T) {
	tests := []struct {
		nums []int
		want int
	}{
		{[]int{2, 4, 9}, 5},
		{[]int{7}, 7},
		{[]int{1, 2}, 1},
		{[]int{-1, -2}, -1},
		{[]int{-5, 5}, 0},
		{[]int{0, 0, 0}, 0},
	}
	for _, tt := range tests {
		noPanic(t, "Average", func() {
			got, err := Average(tt.nums)
			if got != tt.want || err != nil {
				t.Errorf("Average(%v) = %d, %v; want %d, <nil>", tt.nums, got, err, tt.want)
			}
		})
	}
}

func TestAverageEmpty(t *testing.T) {
	const want = "average: runtime error: integer divide by zero"
	for _, nums := range [][]int{nil, {}} {
		noPanic(t, "Average", func() {
			got, err := Average(nums)
			if err == nil || err.Error() != want || got != 0 {
				t.Errorf("Average(%v) = %d, %v; want 0, %q", nums, got, err, want)
			}
		})
	}
}
