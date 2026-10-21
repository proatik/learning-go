package counters

import "testing"

func TestIncrement(t *testing.T) {
	for _, start := range []int{0, 5, -1, 99} {
		n := start
		Increment(&n)
		if n != start+1 {
			t.Errorf("Increment(&n) with n=%d: got %d, want %d", start, n, start+1)
		}
	}
}

func TestIncrementRepeated(t *testing.T) {
	n := 0
	for range 10 {
		Increment(&n)
	}
	if n != 10 {
		t.Errorf("Increment called 10 times from 0: got %d, want 10", n)
	}
}

func TestIncrementSharedPointer(t *testing.T) {
	n := 1
	p, q := &n, &n
	Increment(p)
	Increment(q)
	if n != 3 {
		t.Errorf("Increment through two pointers to the same int: got %d, want 3", n)
	}
}

func TestReset(t *testing.T) {
	for _, start := range []int{0, 5, -7} {
		n := start
		Reset(&n)
		if n != 0 {
			t.Errorf("Reset(&n) with n=%d: got %d, want 0", start, n)
		}
	}
}

func TestNilDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Increment(nil) or Reset(nil) panicked: %v", r)
		}
	}()
	Increment(nil)
	Reset(nil)
}

func TestSafeDeref(t *testing.T) {
	tests := []struct {
		name string
		p    *int
		def  int
		want int
	}{
		{"nil uses default", nil, 9, 9},
		{"nil with zero default", nil, 0, 0},
		{"set to zero beats default", new(0), 9, 0},
		{"set value", new(42), 9, 42},
		{"negative value", new(-3), 100, -3},
	}
	for _, tt := range tests {
		if got := SafeDeref(tt.p, tt.def); got != tt.want {
			t.Errorf("%s: SafeDeref(..., %d) = %d, want %d", tt.name, tt.def, got, tt.want)
		}
	}
}

func TestSum(t *testing.T) {
	tests := []struct {
		name string
		ps   []*int
		want int
	}{
		{"nil slice", nil, 0},
		{"empty slice", []*int{}, 0},
		{"all nil", []*int{nil, nil}, 0},
		{"mixed", []*int{new(1), nil, new(3)}, 4},
		{"negatives", []*int{new(-5), new(2)}, -3},
	}
	for _, tt := range tests {
		if got := Sum(tt.ps); got != tt.want {
			t.Errorf("%s: Sum(...) = %d, want %d", tt.name, got, tt.want)
		}
	}
}

func TestSumSeesLiveValues(t *testing.T) {
	a, b := 1, 2
	ps := []*int{&a, &b}
	a = 10 // the slice holds pointers, so Sum must see the new value
	if got := Sum(ps); got != 12 {
		t.Errorf("Sum after a=10, b=2 = %d, want 12", got)
	}
}
