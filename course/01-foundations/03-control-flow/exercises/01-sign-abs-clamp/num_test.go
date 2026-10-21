package num

import "testing"

func TestSign(t *testing.T) {
	tests := []struct{ n, want int }{
		{-7, -1}, {-1, -1}, {0, 0}, {1, 1}, {42, 1}, {-1000000, -1}, {1000000, 1},
	}
	for _, tt := range tests {
		if got := Sign(tt.n); got != tt.want {
			t.Errorf("Sign(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

func TestAbs(t *testing.T) {
	tests := []struct{ n, want int }{
		{-5, 5}, {3, 3}, {0, 0}, {-1, 1}, {1, 1}, {-123456, 123456},
	}
	for _, tt := range tests {
		if got := Abs(tt.n); got != tt.want {
			t.Errorf("Abs(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

func TestClamp(t *testing.T) {
	tests := []struct{ x, lo, hi, want int }{
		{15, 0, 10, 10},
		{-3, 0, 10, 0},
		{4, 0, 10, 4},
		{0, 0, 10, 0},
		{10, 0, 10, 10},
		{-20, -10, -5, -10},
		{-7, -10, -5, -7},
		{0, -10, -5, -5},
		{3, 5, 5, 5},
		{9, 5, 5, 5},
	}
	for _, tt := range tests {
		if got := Clamp(tt.x, tt.lo, tt.hi); got != tt.want {
			t.Errorf("Clamp(%d, %d, %d) = %d, want %d", tt.x, tt.lo, tt.hi, got, tt.want)
		}
	}
}
