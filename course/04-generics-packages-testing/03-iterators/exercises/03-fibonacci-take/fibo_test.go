package fibo

import (
	"iter"
	"slices"
	"testing"
)

// collect ranges over seq and returns its values, breaking after max values
// when max > 0 (and always after 10000, so a runaway sequence can't hang the
// test). A panic during the loop is returned instead of crashing.
func collect[V any](seq iter.Seq[V], max int) (got []V, panicked any) {
	defer func() { panicked = recover() }()
	for v := range seq {
		got = append(got, v)
		if len(got) == max || len(got) == 10000 {
			break
		}
	}
	return got, nil
}

// counting yields 1, 2, …, n (forever if n < 0) and records how many values
// were pulled and whether it was called at all.
type counting struct {
	n      int
	calls  int
	pulled int
}

func (c *counting) seq() iter.Seq[int] {
	return func(yield func(int) bool) {
		c.calls++
		for i := 1; c.n < 0 || i <= c.n; i++ {
			c.pulled++
			if !yield(i) {
				return
			}
		}
	}
}

func TestFibonacci(t *testing.T) {
	seq := Fibonacci()
	if seq == nil {
		t.Fatalf("Fibonacci() = nil, want a non-nil iterator")
	}
	got, p := collect(seq, 12)
	if p != nil {
		t.Fatalf("break after 12 values of Fibonacci() panicked: %v (stop when yield returns false)", p)
	}
	want := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89}
	if !slices.Equal(got, want) {
		t.Errorf("first 12 values of Fibonacci() = %v, want %v", got, want)
	}
	got, _ = collect(seq, 3)
	if !slices.Equal(got, want[:3]) {
		t.Errorf("Fibonacci() ranged a second time = %v, want %v (each range starts again)", got, want[:3])
	}
}

func TestFibonacciStops(t *testing.T) {
	calls := 0
	func() {
		defer func() { recover() }() // the runaway guard below
		Fibonacci()(func(v int) bool {
			calls++
			if calls > 1000 {
				panic("runaway")
			}
			return v < 50
		})
	}()
	if calls != 11 {
		t.Errorf("Fibonacci() with yield returning false at 55: yield calls = %d, want 11", calls)
	}
}

func TestTake(t *testing.T) {
	tests := []struct {
		src  int // length of the source; -1 means infinite
		n    int
		want []int
	}{
		{-1, 3, []int{1, 2, 3}},
		{-1, 1, []int{1}},
		{5, 5, []int{1, 2, 3, 4, 5}},
		{2, 5, []int{1, 2}},
		{0, 3, nil},
		{-1, 0, nil},
		{5, -2, nil},
	}
	for _, tt := range tests {
		c := &counting{n: tt.src}
		seq := Take(c.seq(), tt.n)
		if seq == nil {
			t.Fatalf("Take(source of %d, %d) = nil, want a non-nil iterator", tt.src, tt.n)
		}
		got, p := collect(seq, -1)
		if p != nil {
			t.Fatalf("range over Take(source of %d, %d) panicked: %v", tt.src, tt.n, p)
		}
		if !slices.Equal(got, tt.want) {
			t.Errorf("Take(source of %d, %d) = %v, want %v", tt.src, tt.n, got, tt.want)
		}
		if c.pulled != len(tt.want) {
			t.Errorf("Take(source of %d, %d): values pulled from the source = %d, want %d", tt.src, tt.n, c.pulled, len(tt.want))
		}
	}
}

func TestTakeFibonacci(t *testing.T) {
	got, p := collect(Take(Fibonacci(), 8), -1)
	if p != nil {
		t.Fatalf("range over Take(Fibonacci(), 8) panicked: %v", p)
	}
	if want := []int{0, 1, 1, 2, 3, 5, 8, 13}; !slices.Equal(got, want) {
		t.Errorf("Take(Fibonacci(), 8) = %v, want %v", got, want)
	}
	words, _ := collect(Take(slices.Values([]string{"a", "b"}), 5), -1)
	if want := []string{"a", "b"}; !slices.Equal(words, want) {
		t.Errorf("Take([a b], 5) = %q, want %q", words, want)
	}
}

func TestTakeIsLazy(t *testing.T) {
	c := &counting{n: -1}
	seq := Take(c.seq(), 3)
	if c.calls != 0 {
		t.Errorf("Take(source, 3) before ranging: source calls = %d, want 0", c.calls)
	}
	collect(seq, -1)
	c2 := &counting{n: -1}
	collect(Take(c2.seq(), 0), -1)
	if c2.calls != 0 {
		t.Errorf("Take(source, 0) ranged: source calls = %d, want 0", c2.calls)
	}
}

func TestTakeBreak(t *testing.T) {
	c := &counting{n: -1}
	got, p := collect(Take(c.seq(), 100), 2)
	if p != nil {
		t.Fatalf("break after 2 values of Take(source, 100) panicked: %v (stop when yield returns false)", p)
	}
	if want := []int{1, 2}; !slices.Equal(got, want) {
		t.Errorf("first 2 values of Take(source, 100) = %v, want %v", got, want)
	}
	if c.pulled != 2 {
		t.Errorf("break after 2 values of Take(source, 100): values pulled from the source = %d, want 2", c.pulled)
	}
}
