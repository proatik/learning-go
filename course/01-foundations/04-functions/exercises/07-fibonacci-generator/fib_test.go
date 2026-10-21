package fib

import "testing"

func TestFibonacci(t *testing.T) {
	f := Fibonacci()
	if f == nil {
		t.Fatal("Fibonacci() returned nil")
	}
	want := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144}
	for i, w := range want {
		if got := f(); got != w {
			t.Fatalf("call %d of the generator returned %d, want %d", i+1, got, w)
		}
	}
}

func TestFibonacciIndependent(t *testing.T) {
	a, b := Fibonacci(), Fibonacci()
	if a == nil || b == nil {
		t.Fatal("Fibonacci() returned nil")
	}
	for range 5 {
		a()
	}
	if got := b(); got != 0 {
		t.Errorf("a fresh generator returned %d on its first call, want 0", got)
	}
	if got := a(); got != 5 {
		t.Errorf("the first generator returned %d on its 6th call, want 5", got)
	}
}

func TestNth(t *testing.T) {
	tests := []struct{ n, want int }{
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 2},
		{10, 55},
		{20, 6765},
		{50, 12586269025},
		{90, 2880067194370816120},
	}
	for _, tt := range tests {
		got, err := Nth(tt.n)
		if err != nil {
			t.Errorf("Nth(%d) returned error %v, want nil", tt.n, err)
			continue
		}
		if got != tt.want {
			t.Errorf("Nth(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

func TestNthNegative(t *testing.T) {
	for _, n := range []int{-1, -10} {
		got, err := Nth(n)
		if err == nil {
			t.Errorf("Nth(%d): err = nil, want an error", n)
		}
		if got != 0 {
			t.Errorf("Nth(%d) = %d, want 0 alongside the error", n, got)
		}
	}
}
