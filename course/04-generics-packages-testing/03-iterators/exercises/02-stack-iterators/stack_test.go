package stack_test

import (
	"iter"
	"slices"
	"testing"

	stack "learngo/course/04-generics-packages-testing/03-iterators/exercises/02-stack-iterators"
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

func build(items ...string) *stack.Stack[string] {
	s := &stack.Stack[string]{}
	for _, v := range items {
		s.Push(v)
	}
	return s
}

func TestAll(t *testing.T) {
	tests := []struct {
		pushed []string
		want   []string
	}{
		{[]string{"a", "b", "c"}, []string{"c", "b", "a"}},
		{[]string{"only"}, []string{"only"}},
		{nil, nil},
		{[]string{"日", "本"}, []string{"本", "日"}},
	}
	for _, tt := range tests {
		s := build(tt.pushed...)
		got, p := collect(s.All(), -1)
		if p != nil {
			t.Fatalf("range over All() after pushing %q panicked: %v", tt.pushed, p)
		}
		if !slices.Equal(got, tt.want) {
			t.Errorf("All() after pushing %q = %q, want %q", tt.pushed, got, tt.want)
		}
		again, _ := collect(s.All(), -1)
		if !slices.Equal(again, tt.want) || s.Len() != len(tt.pushed) {
			t.Errorf("All() a second time after pushing %q = %q (Len %d), want %q (Len %d): All must not modify the stack", tt.pushed, again, s.Len(), tt.want, len(tt.pushed))
		}
	}
}

func TestBackward(t *testing.T) {
	tests := [][]string{{"a", "b", "c"}, {"x"}, nil}
	for _, pushed := range tests {
		s := build(pushed...)
		got, p := collect(s.Backward(), -1)
		if p != nil {
			t.Fatalf("range over Backward() after pushing %q panicked: %v", pushed, p)
		}
		if !slices.Equal(got, pushed) {
			t.Errorf("Backward() after pushing %q = %q, want %q", pushed, got, pushed)
		}
		if s.Len() != len(pushed) {
			t.Errorf("Len() after Backward() on %q = %d, want %d: Backward must not modify the stack", pushed, s.Len(), len(pushed))
		}
	}
}

func TestReadsAtRangeTime(t *testing.T) {
	s := build("a")
	seq := s.All()
	s.Push("b")
	got, _ := collect(seq, -1)
	if want := []string{"b", "a"}; !slices.Equal(got, want) {
		t.Errorf("All() built before Push(%q), ranged after = %q, want %q", "b", got, want)
	}
}

func TestZeroValue(t *testing.T) {
	var s stack.Stack[int]
	for name, seq := range map[string]iter.Seq[int]{"All": s.All(), "Backward": s.Backward(), "Drain": s.Drain()} {
		if seq == nil {
			t.Fatalf("%s() on a zero Stack = nil, want a non-nil iterator", name)
		}
		got, p := collect(seq, -1)
		if p != nil || len(got) != 0 {
			t.Errorf("%s() on a zero Stack = %v (panic %v), want []", name, got, p)
		}
	}
}

func TestBreak(t *testing.T) {
	s := build("a", "b", "c", "d")
	got, p := collect(s.All(), 2)
	if p != nil {
		t.Fatalf("break after 2 values of All() panicked: %v (stop when yield returns false)", p)
	}
	if want := []string{"d", "c"}; !slices.Equal(got, want) {
		t.Errorf("first 2 values of All() on [a b c d] = %q, want %q", got, want)
	}
	got, p = collect(s.Backward(), 1)
	if p != nil {
		t.Fatalf("break after 1 value of Backward() panicked: %v (stop when yield returns false)", p)
	}
	if want := []string{"a"}; !slices.Equal(got, want) {
		t.Errorf("first value of Backward() on [a b c d] = %q, want %q", got, want)
	}
}

func TestDrain(t *testing.T) {
	s := build("a", "b", "c")
	got, p := collect(s.Drain(), -1)
	if p != nil {
		t.Fatalf("range over Drain() panicked: %v", p)
	}
	if want := []string{"c", "b", "a"}; !slices.Equal(got, want) {
		t.Errorf("Drain() after pushing [a b c] = %q, want %q", got, want)
	}
	if s.Len() != 0 {
		t.Errorf("Len() after a full Drain() = %d, want 0", s.Len())
	}
}

func TestDrainBreak(t *testing.T) {
	s := build("a", "b", "c", "d")
	got, p := collect(s.Drain(), 2)
	if p != nil {
		t.Fatalf("break after 2 values of Drain() panicked: %v (stop when yield returns false)", p)
	}
	if want := []string{"d", "c"}; !slices.Equal(got, want) {
		t.Errorf("first 2 values of Drain() on [a b c d] = %q, want %q", got, want)
	}
	if s.Len() != 2 {
		t.Fatalf("Len() after a break at 2 values of Drain() on [a b c d] = %d, want 2", s.Len())
	}
	rest, _ := collect(s.Drain(), -1)
	if want := []string{"b", "a"}; !slices.Equal(rest, want) {
		t.Errorf("Drain() of the rest after the break = %q, want %q", rest, want)
	}
}
