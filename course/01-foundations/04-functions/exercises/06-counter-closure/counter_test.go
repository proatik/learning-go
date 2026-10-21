package counter

import "testing"

func TestCounter(t *testing.T) {
	c := Counter()
	if c == nil {
		t.Fatal("Counter() returned nil")
	}
	for want := 1; want <= 5; want++ {
		if got := c(); got != want {
			t.Errorf("call %d of Counter(): got %d, want %d", want, got, want)
		}
	}
}

func TestCounterIndependent(t *testing.T) {
	a, b := Counter(), Counter()
	if a == nil || b == nil {
		t.Fatal("Counter() returned nil")
	}
	a()
	a()
	a()
	if got := b(); got != 1 {
		t.Errorf("a fresh counter returned %d after another counter was used, want 1", got)
	}
	if got := a(); got != 4 {
		t.Errorf("first counter returned %d on its 4th call, want 4", got)
	}
}

func TestCountBy(t *testing.T) {
	tests := []struct {
		start, step int
		want        []int
	}{
		{10, 5, []int{10, 15, 20, 25}},
		{3, -1, []int{3, 2, 1, 0, -1}},
		{0, 0, []int{0, 0, 0}},
		{-4, 2, []int{-4, -2, 0, 2}},
		{1, 1, []int{1, 2, 3}},
	}
	for _, tt := range tests {
		c := CountBy(tt.start, tt.step)
		if c == nil {
			t.Fatalf("CountBy(%d, %d) returned nil", tt.start, tt.step)
		}
		for i, want := range tt.want {
			if got := c(); got != want {
				t.Errorf("CountBy(%d, %d): call %d returned %d, want %d", tt.start, tt.step, i+1, got, want)
			}
		}
	}
}

func TestCountByIndependent(t *testing.T) {
	a, b := CountBy(0, 10), CountBy(0, 10)
	if a == nil || b == nil {
		t.Fatal("CountBy(0, 10) returned nil")
	}
	a()
	a()
	if got := b(); got != 0 {
		t.Errorf("a fresh CountBy(0, 10) returned %d on its first call, want 0", got)
	}
}

func TestAccumulator(t *testing.T) {
	acc := Accumulator(100)
	if acc == nil {
		t.Fatal("Accumulator(100) returned nil")
	}
	steps := []struct{ add, want int }{
		{10, 110},
		{-30, 80},
		{0, 80},
		{20, 100},
	}
	for _, s := range steps {
		if got := acc(s.add); got != s.want {
			t.Errorf("acc(%d) = %d, want %d", s.add, got, s.want)
		}
	}

	other := Accumulator(0)
	if other == nil {
		t.Fatal("Accumulator(0) returned nil")
	}
	if got := other(5); got != 5 {
		t.Errorf("a second Accumulator(0)(5) = %d, want 5 (state must not be shared)", got)
	}
}
