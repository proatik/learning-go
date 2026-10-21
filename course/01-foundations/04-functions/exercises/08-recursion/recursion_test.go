package recursion

import "testing"

func TestFactorial(t *testing.T) {
	tests := []struct{ n, want int }{
		{0, 1},
		{1, 1},
		{2, 2},
		{5, 120},
		{10, 3628800},
		{20, 2432902008176640000},
	}
	for _, tt := range tests {
		got, err := Factorial(tt.n)
		if err != nil {
			t.Errorf("Factorial(%d) returned error %v, want nil", tt.n, err)
			continue
		}
		if got != tt.want {
			t.Errorf("Factorial(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

func TestFactorialErrors(t *testing.T) {
	for _, n := range []int{-1, -20, 21, 100} {
		got, err := Factorial(n)
		if err == nil {
			t.Errorf("Factorial(%d): err = nil, want an error", n)
		}
		if got != 0 {
			t.Errorf("Factorial(%d) = %d, want 0 alongside the error", n, got)
		}
	}
}

func TestPower(t *testing.T) {
	tests := []struct{ base, exp, want int }{
		{2, 0, 1},
		{7, 0, 1},
		{0, 0, 1},
		{0, 5, 0},
		{2, 1, 2},
		{2, 10, 1024},
		{3, 5, 243},
		{-3, 3, -27},
		{-2, 4, 16},
		{10, 18, 1_000_000_000_000_000_000},
		{2, 62, 1 << 62},
		{1, 1_000_000, 1},
		{-1, 1_000_001, -1},
	}
	for _, tt := range tests {
		if got := Power(tt.base, tt.exp); got != tt.want {
			t.Errorf("Power(%d, %d) = %d, want %d", tt.base, tt.exp, got, tt.want)
		}
	}
}

func TestGCD(t *testing.T) {
	tests := []struct{ a, b, want int }{
		{48, 18, 6},
		{18, 48, 6},
		{17, 5, 1},
		{0, 5, 5},
		{5, 0, 5},
		{0, 0, 0},
		{12, 12, 12},
		{270, 192, 6},
		{1071, 462, 21},
	}
	for _, tt := range tests {
		if got := GCD(tt.a, tt.b); got != tt.want {
			t.Errorf("GCD(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestSumDigits(t *testing.T) {
	tests := []struct{ n, want int }{
		{0, 0},
		{7, 7},
		{10, 1},
		{1234, 10},
		{99999, 45},
		{-56, 11},
		{-9, 9},
		{1_000_000_001, 2},
	}
	for _, tt := range tests {
		if got := SumDigits(tt.n); got != tt.want {
			t.Errorf("SumDigits(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}
