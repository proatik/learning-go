package average

import (
	"math"
	"testing"
)

func near(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestIntAverage(t *testing.T) {
	tests := []struct{ a, b, c, want int }{
		{1, 2, 2, 1},
		{3, 3, 3, 3},
		{1, 2, 3, 2},
		{0, 0, 2, 0},
		{-1, -2, -2, -1},
		{-3, 0, 0, -1},
		{-1, 0, 0, 0},
		{10, 20, 31, 20},
	}
	for _, tt := range tests {
		if got := IntAverage(tt.a, tt.b, tt.c); got != tt.want {
			t.Errorf("IntAverage(%d, %d, %d) = %d, want %d", tt.a, tt.b, tt.c, got, tt.want)
		}
	}
}

func TestAverage(t *testing.T) {
	tests := []struct {
		a, b, c int
		want    float64
	}{
		{1, 2, 2, 5.0 / 3},
		{3, 3, 3, 3},
		{0, 0, 1, 1.0 / 3},
		{-1, 0, 0, -1.0 / 3},
		{1, 2, 3, 2},
		{10, 20, 31, 61.0 / 3},
		{0, 0, 0, 0},
	}
	for _, tt := range tests {
		if got := Average(tt.a, tt.b, tt.c); !near(got, tt.want) {
			t.Errorf("Average(%d, %d, %d) = %v, want %v", tt.a, tt.b, tt.c, got, tt.want)
		}
	}
}

func TestPercent(t *testing.T) {
	tests := []struct {
		part, total int
		want        float64
	}{
		{1, 3, 100.0 / 3},
		{1, 4, 25},
		{7, 4, 175},
		{0, 9, 0},
		{-1, 2, -50},
		{1, 8, 12.5},
	}
	for _, tt := range tests {
		if got := Percent(tt.part, tt.total); !near(got, tt.want) {
			t.Errorf("Percent(%d, %d) = %v, want %v", tt.part, tt.total, got, tt.want)
		}
	}
}
