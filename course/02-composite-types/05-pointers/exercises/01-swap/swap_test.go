package swap

import "testing"

func TestSwap(t *testing.T) {
	tests := []struct{ x, y int }{
		{1, 2},
		{2, 1},
		{-5, 10},
		{0, 0},
		{42, -42},
	}
	for _, tt := range tests {
		x, y := tt.x, tt.y
		Swap(&x, &y)
		if x != tt.y || y != tt.x {
			t.Errorf("Swap(&%d, &%d): got x=%d y=%d, want x=%d y=%d", tt.x, tt.y, x, y, tt.y, tt.x)
		}
	}
}

func TestSwapSameVariable(t *testing.T) {
	x := 7
	Swap(&x, &x)
	if x != 7 {
		t.Errorf("Swap(&x, &x) with x=7: got x=%d, want 7", x)
	}
}

func TestSwapTwiceRestores(t *testing.T) {
	x, y := 3, 8
	Swap(&x, &y)
	Swap(&x, &y)
	if x != 3 || y != 8 {
		t.Errorf("Swap twice: got x=%d y=%d, want x=3 y=8", x, y)
	}
}

func TestSwapSliceElements(t *testing.T) {
	s := []int{10, 20, 30}
	Swap(&s[0], &s[2])
	if s[0] != 30 || s[1] != 20 || s[2] != 10 {
		t.Errorf("Swap(&s[0], &s[2]) on [10 20 30]: got %v, want [30 20 10]", s)
	}
}

func TestSwapIfGreater(t *testing.T) {
	tests := []struct {
		x, y         int
		wantX, wantY int
		want         bool
	}{
		{9, 4, 4, 9, true},
		{4, 9, 4, 9, false},
		{5, 5, 5, 5, false},
		{0, -1, -1, 0, true},
		{-3, -2, -3, -2, false},
	}
	for _, tt := range tests {
		x, y := tt.x, tt.y
		got := SwapIfGreater(&x, &y)
		if got != tt.want {
			t.Errorf("SwapIfGreater(&%d, &%d) = %v, want %v", tt.x, tt.y, got, tt.want)
		}
		if x != tt.wantX || y != tt.wantY {
			t.Errorf("SwapIfGreater(&%d, &%d): got x=%d y=%d, want x=%d y=%d", tt.x, tt.y, x, y, tt.wantX, tt.wantY)
		}
	}
}
