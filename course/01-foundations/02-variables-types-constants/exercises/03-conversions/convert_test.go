package convert

import (
	"math"
	"testing"
)

func near(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestBMI(t *testing.T) {
	tests := []struct {
		w    float64
		h    int
		want float64
	}{
		{70, 175, 22.857142857142858},
		{90, 200, 22.5},
		{50, 150, 22.22222222222222},
		{100, 100, 100},
		{64.5, 168, 22.85289115646259},
	}
	for _, tt := range tests {
		if got := BMI(tt.w, tt.h); !near(got, tt.want) {
			t.Errorf("BMI(%v, %v) = %v, want %v", tt.w, tt.h, got, tt.want)
		}
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		f    float64
		want int
	}{
		{3.99, 3}, {3.0, 3}, {0.5, 0}, {-0.5, 0}, {-3.99, -3}, {-7, -7}, {1e9 + 0.75, 1000000000},
	}
	for _, tt := range tests {
		if got := Truncate(tt.f); got != tt.want {
			t.Errorf("Truncate(%v) = %v, want %v", tt.f, got, tt.want)
		}
	}
}

func TestPixels(t *testing.T) {
	tests := []struct {
		in   float64
		dpi  int
		want int
	}{
		{2.5, 96, 240}, {1.1, 72, 79}, {0, 300, 0}, {3.3, 3, 9}, {8.5, 300, 2550},
	}
	for _, tt := range tests {
		if got := Pixels(tt.in, tt.dpi); got != tt.want {
			t.Errorf("Pixels(%v, %v) = %v, want %v", tt.in, tt.dpi, got, tt.want)
		}
	}
}
