package iterate

import "testing"

func TestCollatzSteps(t *testing.T) {
	tests := []struct{ n, want int }{
		{1, 0}, {2, 1}, {3, 7}, {4, 2}, {6, 8}, {7, 16},
		{27, 111}, {97, 118}, {1000000, 152},
		{0, -1}, {-1, -1}, {-27, -1},
	}
	for _, tt := range tests {
		if got := CollatzSteps(tt.n); got != tt.want {
			t.Errorf("CollatzSteps(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

func TestGCD(t *testing.T) {
	tests := []struct{ a, b, want int }{
		{48, 18, 6}, {18, 48, 6}, {17, 5, 1}, {12, 12, 12},
		{100, 75, 25}, {1071, 462, 21}, {1, 999, 1},
		{0, 9, 9}, {9, 0, 9}, {0, 0, 0},
		{-12, 18, 6}, {12, -18, 6}, {-12, -18, 6}, {-7, 0, 7},
	}
	for _, tt := range tests {
		if got := GCD(tt.a, tt.b); got != tt.want {
			t.Errorf("GCD(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}
