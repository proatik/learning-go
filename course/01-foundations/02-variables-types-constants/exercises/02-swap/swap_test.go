package swap

import "testing"

func TestSwap(t *testing.T) {
	tests := []struct{ a, b int }{{1, 2}, {-5, 7}, {0, 42}, {3, 3}, {-1, 0}}
	for _, tt := range tests {
		x, y := Swap(tt.a, tt.b)
		if x != tt.b || y != tt.a {
			t.Errorf("Swap(%d, %d) = %d, %d, want %d, %d", tt.a, tt.b, x, y, tt.b, tt.a)
		}
	}
}

func TestRotate(t *testing.T) {
	tests := []struct{ a, b, c int }{{1, 2, 3}, {-1, 0, 1}, {9, 9, 1}, {0, 0, 7}}
	for _, tt := range tests {
		x, y, z := Rotate(tt.a, tt.b, tt.c)
		if x != tt.b || y != tt.c || z != tt.a {
			t.Errorf("Rotate(%d, %d, %d) = %d, %d, %d, want %d, %d, %d",
				tt.a, tt.b, tt.c, x, y, z, tt.b, tt.c, tt.a)
		}
	}
}

func TestSwapStrings(t *testing.T) {
	tests := []struct{ a, b string }{{"go", "lang"}, {"", "x"}, {"héllo", "wörld"}}
	for _, tt := range tests {
		x, y := SwapStrings(tt.a, tt.b)
		if x != tt.b || y != tt.a {
			t.Errorf("SwapStrings(%q, %q) = %q, %q, want %q, %q", tt.a, tt.b, x, y, tt.b, tt.a)
		}
	}
}
