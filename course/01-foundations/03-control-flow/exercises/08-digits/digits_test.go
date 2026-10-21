package digits

import "testing"

func TestDigitSum(t *testing.T) {
	tests := []struct{ n, want int }{
		{1234, 10}, {0, 0}, {7, 7}, {10, 1}, {99999, 45}, {1000000, 1},
		{-123, 6}, {-9, 9},
	}
	for _, tt := range tests {
		if got := DigitSum(tt.n); got != tt.want {
			t.Errorf("DigitSum(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

func TestReverse(t *testing.T) {
	tests := []struct{ n, want int }{
		{1234, 4321}, {120, 21}, {1000, 1}, {7, 7}, {0, 0}, {12321, 12321},
		{907, 709}, {-56, -65}, {-120, -21}, {-5, -5},
	}
	for _, tt := range tests {
		if got := Reverse(tt.n); got != tt.want {
			t.Errorf("Reverse(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		n    int
		want bool
	}{
		{12321, true}, {1221, true}, {7, true}, {0, true}, {11, true},
		{10, false}, {123, false}, {1231, false}, {100, false},
		{-121, false}, {-7, false},
	}
	for _, tt := range tests {
		if got := IsPalindrome(tt.n); got != tt.want {
			t.Errorf("IsPalindrome(%d) = %t, want %t", tt.n, got, tt.want)
		}
	}
}
