package countdown

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

// yieldCalls calls seq directly with a yield that returns false on call
// number stopAt, and reports how many times yield was called in total.
func yieldCalls[V any](seq iter.Seq[V], stopAt int) int {
	calls := 0
	seq(func(V) bool {
		calls++
		return calls < stopAt
	})
	return calls
}

func TestCountdown(t *testing.T) {
	tests := []struct {
		n    int
		want []int
	}{
		{3, []int{3, 2, 1}},
		{1, []int{1}},
		{5, []int{5, 4, 3, 2, 1}},
		{0, nil},
		{-4, nil},
	}
	for _, tt := range tests {
		seq := Countdown(tt.n)
		if seq == nil {
			t.Fatalf("Countdown(%d) = nil, want a non-nil iterator", tt.n)
		}
		got, p := collect(seq, -1)
		if p != nil {
			t.Fatalf("range over Countdown(%d) panicked: %v", tt.n, p)
		}
		if !slices.Equal(got, tt.want) {
			t.Errorf("Countdown(%d) = %v, want %v", tt.n, got, tt.want)
		}
	}
}

func TestCountdownRunsAgain(t *testing.T) {
	seq := Countdown(2)
	first, _ := collect(seq, -1)
	second, _ := collect(seq, -1)
	if !slices.Equal(second, []int{2, 1}) || !slices.Equal(first, second) {
		t.Errorf("Countdown(2) ranged twice = %v then %v, want [2 1] both times", first, second)
	}
}

func TestCountdownBreak(t *testing.T) {
	got, p := collect(Countdown(1_000_000_000), 3)
	if p != nil {
		t.Fatalf("break after 3 values of Countdown(1e9) panicked: %v (stop when yield returns false)", p)
	}
	if want := []int{1_000_000_000, 999_999_999, 999_999_998}; !slices.Equal(got, want) {
		t.Errorf("first 3 values of Countdown(1e9) = %v, want %v", got, want)
	}
	for _, stopAt := range []int{1, 2, 5} {
		if got := yieldCalls(Countdown(10), stopAt); got != stopAt {
			t.Errorf("Countdown(10) with yield returning false on call %d: yield calls = %d, want %d", stopAt, got, stopAt)
		}
	}
}

func TestRepeat(t *testing.T) {
	gotS, p := collect(Repeat("go", 3), -1)
	if p != nil {
		t.Fatalf("range over Repeat(%q, 3) panicked: %v", "go", p)
	}
	if want := []string{"go", "go", "go"}; !slices.Equal(gotS, want) {
		t.Errorf("Repeat(%q, 3) = %q, want %q", "go", gotS, want)
	}
	for _, n := range []int{0, -1} {
		got, _ := collect(Repeat(7, n), -1)
		if len(got) != 0 {
			t.Errorf("Repeat(7, %d) = %v, want []", n, got)
		}
	}
	got, _ := collect(Repeat(true, 1), -1)
	if !slices.Equal(got, []bool{true}) {
		t.Errorf("Repeat(true, 1) = %v, want [true]", got)
	}
}

func TestRepeatBreak(t *testing.T) {
	got, p := collect(Repeat('x', 1_000_000_000), 2)
	if p != nil {
		t.Fatalf("break after 2 values of Repeat('x', 1e9) panicked: %v (stop when yield returns false)", p)
	}
	if len(got) != 2 {
		t.Errorf("first 2 values of Repeat('x', 1e9) = %q, want %q", got, []rune{'x', 'x'})
	}
	for _, stopAt := range []int{1, 3} {
		if got := yieldCalls(Repeat("a", 10), stopAt); got != stopAt {
			t.Errorf("Repeat(%q, 10) with yield returning false on call %d: yield calls = %d, want %d", "a", stopAt, got, stopAt)
		}
	}
}
