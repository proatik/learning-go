package stats_test

import (
	"math"
	"testing"

	"learngo/course/04-generics-packages-testing/04-packages/exercises/01-split-package/stats"
)

func TestSum(t *testing.T) {
	tests := []struct {
		in   []int
		want int
	}{
		{nil, 0},
		{[]int{}, 0},
		{[]int{5}, 5},
		{[]int{70, 80, 95}, 245},
		{[]int{-3, 3, -4}, -4},
	}
	for _, tt := range tests {
		if got := stats.Sum(tt.in); got != tt.want {
			t.Errorf("Sum(%v) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestMean(t *testing.T) {
	tests := []struct {
		in   []int
		want float64
	}{
		{nil, 0},
		{[]int{}, 0},
		{[]int{7}, 7},
		{[]int{70, 80, 95}, 245.0 / 3},
		{[]int{1, 2}, 1.5},
		{[]int{-2, -4}, -3},
	}
	for _, tt := range tests {
		if got := stats.Mean(tt.in); math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("Mean(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestMinMax(t *testing.T) {
	tests := []struct {
		in     []int
		lo, hi int
		ok     bool
	}{
		{nil, 0, 0, false},
		{[]int{}, 0, 0, false},
		{[]int{4}, 4, 4, true},
		{[]int{3, -1, 7}, -1, 7, true},
		{[]int{5, 5, 5}, 5, 5, true},
		{[]int{9, 8, 7, 6}, 6, 9, true},
	}
	for _, tt := range tests {
		lo, hi, ok := stats.MinMax(tt.in)
		if lo != tt.lo || hi != tt.hi || ok != tt.ok {
			t.Errorf("MinMax(%v) = %d, %d, %v, want %d, %d, %v", tt.in, lo, hi, ok, tt.lo, tt.hi, tt.ok)
		}
	}
}
