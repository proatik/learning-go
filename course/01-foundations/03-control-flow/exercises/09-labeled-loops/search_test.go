package search

import "testing"

func TestFirstPair(t *testing.T) {
	tests := []struct{ target, limit, wantA, wantB int }{
		{12, 20, 1, 12},
		{12, 12, 1, 12},
		{12, 11, 2, 6},
		{12, 10, 2, 6},
		{12, 5, 3, 4},
		{12, 3, 0, 0},
		{49, 10, 7, 7},
		{42, 9, 6, 7},
		{1, 1, 1, 1},
		{13, 10, 0, 0},
		{13, 13, 1, 13},
		{0, 10, 0, 0},
		{-6, 10, 0, 0},
		{6, 0, 0, 0},
		{6, -5, 0, 0},
	}
	for _, tt := range tests {
		a, b := FirstPair(tt.target, tt.limit)
		if a != tt.wantA || b != tt.wantB {
			t.Errorf("FirstPair(%d, %d) = (%d, %d), want (%d, %d)",
				tt.target, tt.limit, a, b, tt.wantA, tt.wantB)
		}
	}
}

func TestCountPrimes(t *testing.T) {
	tests := []struct{ n, want int }{
		{10, 4}, {2, 1}, {3, 2}, {4, 2}, {11, 5}, {25, 9}, {49, 15},
		{100, 25}, {1000, 168},
		{1, 0}, {0, 0}, {-10, 0},
	}
	for _, tt := range tests {
		if got := CountPrimes(tt.n); got != tt.want {
			t.Errorf("CountPrimes(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}
