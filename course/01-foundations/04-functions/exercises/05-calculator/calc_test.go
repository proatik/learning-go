package calc

import "testing"

func TestLookup(t *testing.T) {
	tests := []struct {
		op   string
		a, b int
		want int
	}{
		{"+", 2, 3, 5},
		{"-", 2, 3, -1},
		{"*", -4, 3, -12},
		{"/", 7, 2, 3},
		{"%", 7, 2, 1},
		{"/", -7, 2, -3},
	}
	for _, tt := range tests {
		f, err := Lookup(tt.op)
		if err != nil {
			t.Errorf("Lookup(%q) returned error %v, want nil", tt.op, err)
			continue
		}
		if f == nil {
			t.Errorf("Lookup(%q) returned a nil function", tt.op)
			continue
		}
		if got := f(tt.a, tt.b); got != tt.want {
			t.Errorf("Lookup(%q)(%d, %d) = %d, want %d", tt.op, tt.a, tt.b, got, tt.want)
		}
	}
}

func TestLookupUnknown(t *testing.T) {
	for _, op := range []string{"^", "", "plus", "++", " +"} {
		f, err := Lookup(op)
		if err == nil {
			t.Errorf("Lookup(%q): err = nil, want an error", op)
		}
		if f != nil {
			t.Errorf("Lookup(%q) returned a non-nil function alongside the error", op)
		}
	}
}

func TestCalc(t *testing.T) {
	tests := []struct {
		op   string
		a, b int
		want int
	}{
		{"+", 2, 3, 5},
		{"+", -2, -3, -5},
		{"-", 2, 3, -1},
		{"*", 6, 7, 42},
		{"*", 0, 99, 0},
		{"/", 7, 2, 3},
		{"/", 0, 5, 0},
		{"%", 7, 3, 1},
		{"%", -7, 3, -1},
	}
	for _, tt := range tests {
		got, err := Calc(tt.op, tt.a, tt.b)
		if err != nil {
			t.Errorf("Calc(%q, %d, %d) returned error %v, want nil", tt.op, tt.a, tt.b, err)
			continue
		}
		if got != tt.want {
			t.Errorf("Calc(%q, %d, %d) = %d, want %d", tt.op, tt.a, tt.b, got, tt.want)
		}
	}
}

func TestCalcErrors(t *testing.T) {
	tests := []struct {
		op   string
		a, b int
	}{
		{"/", 7, 0},
		{"%", 7, 0},
		{"/", 0, 0},
		{"^", 2, 3},
		{"", 1, 1},
		{"?", 1, 0},
	}
	for _, tt := range tests {
		got, err := Calc(tt.op, tt.a, tt.b)
		if err == nil {
			t.Errorf("Calc(%q, %d, %d): err = nil, want an error", tt.op, tt.a, tt.b)
		}
		if got != 0 {
			t.Errorf("Calc(%q, %d, %d) = %d, want 0 alongside the error", tt.op, tt.a, tt.b, got)
		}
	}
}

func TestCalcMultiplyByZeroIsFine(t *testing.T) {
	got, err := Calc("*", 5, 0)
	if err != nil || got != 0 {
		t.Errorf("Calc(\"*\", 5, 0) = (%d, %v), want (0, nil)", got, err)
	}
	got, err = Calc("-", 5, 0)
	if err != nil || got != 5 {
		t.Errorf("Calc(\"-\", 5, 0) = (%d, %v), want (5, nil)", got, err)
	}
}
