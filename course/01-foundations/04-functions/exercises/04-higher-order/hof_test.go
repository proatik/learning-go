package hof

import "testing"

func inc(x int) int    { return x + 1 }
func dbl(x int) int    { return x * 2 }
func square(x int) int { return x * x }
func neg(x int) int    { return -x }

func TestApplyN(t *testing.T) {
	tests := []struct {
		name string
		f    IntFunc
		n, x int
		want int
	}{
		{"dbl", dbl, 3, 1, 8},
		{"dbl", dbl, 10, 1, 1024},
		{"inc", inc, 0, 42, 42},
		{"inc", inc, -3, 42, 42},
		{"inc", inc, 5, -2, 3},
		{"square", square, 2, 3, 81},
		{"neg", neg, 3, 7, -7},
		{"neg", neg, 4, 7, 7},
	}
	for _, tt := range tests {
		if got := ApplyN(tt.f, tt.n, tt.x); got != tt.want {
			t.Errorf("ApplyN(%s, %d, %d) = %d, want %d", tt.name, tt.n, tt.x, got, tt.want)
		}
	}
}

func TestApplyNCallCount(t *testing.T) {
	calls := 0
	counting := func(x int) int {
		calls++
		return x
	}
	ApplyN(counting, 6, 0)
	if calls != 6 {
		t.Errorf("ApplyN(f, 6, 0) called f %d times, want 6", calls)
	}
}

func TestCompose(t *testing.T) {
	tests := []struct {
		name string
		f, g IntFunc
		x    int
		want int
	}{
		{"Compose(inc, dbl)", inc, dbl, 5, 11},
		{"Compose(dbl, inc)", dbl, inc, 5, 12},
		{"Compose(square, inc)", square, inc, 2, 9},
		{"Compose(inc, square)", inc, square, 2, 5},
		{"Compose(neg, neg)", neg, neg, 9, 9},
	}
	for _, tt := range tests {
		h := Compose(tt.f, tt.g)
		if h == nil {
			t.Fatalf("%s returned nil", tt.name)
		}
		if got := h(tt.x); got != tt.want {
			t.Errorf("%s(%d) = %d, want %d", tt.name, tt.x, got, tt.want)
		}
	}
}

func TestComposeIsLazy(t *testing.T) {
	calls := 0
	counting := func(x int) int {
		calls++
		return x
	}
	h := Compose(counting, counting)
	if h == nil {
		t.Fatal("Compose returned nil")
	}
	if calls != 0 {
		t.Fatalf("Compose called its arguments %d times before the result was used, want 0", calls)
	}
	h(1)
	h(2)
	if calls != 4 {
		t.Errorf("calling the composed function twice made %d calls, want 4", calls)
	}
}

func TestTwice(t *testing.T) {
	tests := []struct {
		name string
		f    IntFunc
		x    int
		want int
	}{
		{"Twice(inc)", inc, 5, 7},
		{"Twice(dbl)", dbl, 3, 12},
		{"Twice(square)", square, 3, 81},
		{"Twice(neg)", neg, 4, 4},
	}
	for _, tt := range tests {
		h := Twice(tt.f)
		if h == nil {
			t.Fatalf("%s returned nil", tt.name)
		}
		if got := h(tt.x); got != tt.want {
			t.Errorf("%s(%d) = %d, want %d", tt.name, tt.x, got, tt.want)
		}
	}

	inner := Twice(dbl)
	if inner == nil {
		t.Fatal("Twice(dbl) returned nil")
	}
	if outer := Twice(inner); outer == nil {
		t.Error("Twice(Twice(dbl)) returned nil")
	} else if got := outer(1); got != 16 {
		t.Errorf("Twice(Twice(dbl))(1) = %d, want 16", got)
	}
}
