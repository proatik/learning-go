package sums

import "testing"

func TestSumTo(t *testing.T) {
	tests := []struct{ n, want int }{
		{4, 10}, {1, 1}, {2, 3}, {10, 55}, {100, 5050},
		{0, 0}, {-1, 0}, {-10, 0},
	}
	for _, tt := range tests {
		if got := SumTo(tt.n); got != tt.want {
			t.Errorf("SumTo(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

func TestSumOfMultiples(t *testing.T) {
	tests := []struct{ limit, a, b, want int }{
		{10, 3, 5, 23},
		{1000, 3, 5, 233168},
		{20, 7, 7, 21},
		{16, 3, 5, 60}, // 3+5+6+9+10+12+15: 15 counted once
		{4, 3, 5, 3},
		{3, 3, 5, 0}, // limit is excluded
		{6, 1, 100, 15},
		{1, 3, 5, 0},
		{0, 3, 5, 0},
		{-10, 3, 5, 0},
	}
	for _, tt := range tests {
		if got := SumOfMultiples(tt.limit, tt.a, tt.b); got != tt.want {
			t.Errorf("SumOfMultiples(%d, %d, %d) = %d, want %d", tt.limit, tt.a, tt.b, got, tt.want)
		}
	}
}

func TestCountDivisors(t *testing.T) {
	tests := []struct{ n, want int }{
		{12, 6}, {1, 1}, {2, 2}, {7, 2}, {16, 5}, {36, 9}, {100, 9},
		{0, 0}, {-6, 0},
	}
	for _, tt := range tests {
		if got := CountDivisors(tt.n); got != tt.want {
			t.Errorf("CountDivisors(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}
