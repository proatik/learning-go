package divide

import "testing"

func TestDivide(t *testing.T) {
	tests := []struct {
		a, b, q, r int
	}{
		{17, 5, 3, 2},
		{10, 2, 5, 0},
		{0, 4, 0, 0},
		{3, 7, 0, 3},
		{-7, 2, -3, -1},
		{7, -2, -3, 1},
		{-9, -3, 3, 0},
		{1, 1, 1, 0},
	}
	for _, tt := range tests {
		q, r, err := Divide(tt.a, tt.b)
		if err != nil {
			t.Errorf("Divide(%d, %d) returned error %v, want nil", tt.a, tt.b, err)
			continue
		}
		if q != tt.q || r != tt.r {
			t.Errorf("Divide(%d, %d) = (%d, %d), want (%d, %d)", tt.a, tt.b, q, r, tt.q, tt.r)
		}
	}
}

func TestDivideByZero(t *testing.T) {
	for _, a := range []int{5, 0, -3} {
		q, r, err := Divide(a, 0)
		if err == nil {
			t.Errorf("Divide(%d, 0): err = nil, want an error", a)
		}
		if q != 0 || r != 0 {
			t.Errorf("Divide(%d, 0) = (%d, %d), want (0, 0) alongside the error", a, q, r)
		}
	}
}
