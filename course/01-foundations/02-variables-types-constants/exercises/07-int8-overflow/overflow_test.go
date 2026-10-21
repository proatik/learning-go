package overflow

import "testing"

func TestAddInt8(t *testing.T) {
	tests := []struct {
		a, b     int8
		sum      int8
		overflow bool
	}{
		{1, 2, 3, false},
		{0, 0, 0, false},
		{100, 27, 127, false},
		{100, 28, -128, true},
		{127, 1, -128, true},
		{127, 127, -2, true},
		{-100, -28, -128, false},
		{-100, -29, 127, true},
		{-128, -1, 127, true},
		{-128, -128, 0, true},
		{-128, 127, -1, false},
		{-5, 3, -2, false},
	}
	for _, tt := range tests {
		sum, ovf := AddInt8(tt.a, tt.b)
		if sum != tt.sum || ovf != tt.overflow {
			t.Errorf("AddInt8(%d, %d) = %d, %v, want %d, %v", tt.a, tt.b, sum, ovf, tt.sum, tt.overflow)
		}
	}
}

// TestAddInt8Exhaustive checks every one of the 65 536 possible inputs.
func TestAddInt8Exhaustive(t *testing.T) {
	failures := 0
	for x := -128; x <= 127; x++ {
		for y := -128; y <= 127; y++ {
			a, b := int8(x), int8(y)
			want := a + b // Go's own wrapping addition
			wantOvf := x+y > 127 || x+y < -128
			sum, ovf := AddInt8(a, b)
			if (sum != want || ovf != wantOvf) && failures < 5 {
				failures++
				t.Errorf("AddInt8(%d, %d) = %d, %v, want %d, %v", a, b, sum, ovf, want, wantOvf)
			}
		}
	}
}

func TestSaturatingAdd(t *testing.T) {
	tests := []struct{ a, b, want int8 }{
		{1, 2, 3},
		{0, 0, 0},
		{100, 27, 127},
		{100, 28, 127},
		{127, 127, 127},
		{-100, -28, -128},
		{-100, -29, -128},
		{-128, -128, -128},
		{-128, 127, -1},
		{-5, 3, -2},
		{50, -60, -10},
	}
	for _, tt := range tests {
		if got := SaturatingAdd(tt.a, tt.b); got != tt.want {
			t.Errorf("SaturatingAdd(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}
