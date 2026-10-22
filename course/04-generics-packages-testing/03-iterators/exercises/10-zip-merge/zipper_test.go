package zipper_test

import (
	"fmt"
	"iter"
	"slices"
	"testing"

	zipper "learngo/course/04-generics-packages-testing/03-iterators/exercises/10-zip-merge"
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

// collect2 is collect for a Seq2; each pair is formatted as "k:v".
func collect2[K, V any](seq iter.Seq2[K, V], max int) (got []string, panicked any) {
	defer func() { panicked = recover() }()
	for k, v := range seq {
		got = append(got, fmt.Sprintf("%v:%v", k, v))
		if len(got) == max || len(got) == 10000 {
			break
		}
	}
	return got, nil
}

// source yields vals, or start, start+step, … forever when vals is nil.
// open counts walks that started but haven't run their cleanup yet.
type source struct {
	vals        []int
	start, step int
	started     bool
	open        int
	pulled      int
}

func (s *source) seq() iter.Seq[int] {
	return func(yield func(int) bool) {
		s.started = true
		s.open++
		defer func() { s.open-- }()
		if s.vals != nil {
			for _, v := range s.vals {
				s.pulled++
				if !yield(v) {
					return
				}
			}
			return
		}
		for v := s.start; ; v += s.step {
			s.pulled++
			if !yield(v) {
				return
			}
		}
	}
}

func list(vals ...int) *source      { return &source{vals: append([]int{}, vals...)} }
func count(start, step int) *source { return &source{start: start, step: step} }

func checkClosed(t *testing.T, what string, srcs ...*source) {
	t.Helper()
	for i, s := range srcs {
		if s.open != 0 {
			t.Errorf("%s: input %d still open after the loop = true, want false (defer stop() after iter.Pull)", what, i+1)
		}
	}
}

func TestZip(t *testing.T) {
	tests := []struct {
		as, bs []int
		want   []string
	}{
		{[]int{1, 2, 3}, []int{10, 20}, []string{"1:10", "2:20"}},
		{[]int{1}, []int{10, 20, 30}, []string{"1:10"}},
		{[]int{1, 2}, []int{10, 20}, []string{"1:10", "2:20"}},
		{[]int{}, []int{10}, nil},
		{[]int{1}, []int{}, nil},
	}
	for _, tt := range tests {
		a, b := list(tt.as...), list(tt.bs...)
		seq := zipper.Zip(a.seq(), b.seq())
		if seq == nil {
			t.Fatalf("Zip(%v, %v) = nil, want a non-nil iterator", tt.as, tt.bs)
		}
		got, p := collect2(seq, -1)
		if p != nil {
			t.Fatalf("range over Zip(%v, %v) panicked: %v", tt.as, tt.bs, p)
		}
		if !slices.Equal(got, tt.want) {
			t.Errorf("Zip(%v, %v) = %q, want %q", tt.as, tt.bs, got, tt.want)
		}
		checkClosed(t, fmt.Sprintf("Zip(%v, %v)", tt.as, tt.bs), a, b)
	}
}

func TestZipTypes(t *testing.T) {
	got, p := collect2(zipper.Zip(slices.Values([]string{"x", "y", "z"}), count(0, 1).seq()), -1)
	if p != nil {
		t.Fatalf("range over Zip([x y z], 0, 1, 2, …) panicked: %v", p)
	}
	if want := []string{"x:0", "y:1", "z:2"}; !slices.Equal(got, want) {
		t.Errorf("Zip([x y z], 0, 1, 2, …) = %q, want %q", got, want)
	}
}

func TestZipIsLazy(t *testing.T) {
	a, b := count(0, 1), count(100, 1)
	seq := zipper.Zip(a.seq(), b.seq())
	if a.started || b.started {
		t.Errorf("Zip(a, b) before ranging: started a = %v, b = %v, want false, false", a.started, b.started)
	}
	got, p := collect2(seq, 3)
	if p != nil {
		t.Fatalf("break after 3 pairs of Zip(0, 1, …, 100, 101, …) panicked: %v (stop when yield returns false)", p)
	}
	if want := []string{"0:100", "1:101", "2:102"}; !slices.Equal(got, want) {
		t.Errorf("first 3 pairs of Zip(0, 1, …, 100, 101, …) = %q, want %q", got, want)
	}
	if a.pulled != 3 || b.pulled != 3 {
		t.Errorf("break after 3 pairs of Zip(infinite, infinite): values pulled = %d and %d, want 3 and 3", a.pulled, b.pulled)
	}
	checkClosed(t, "break in Zip(infinite, infinite)", a, b)

	a, b = list(1, 2, 3), count(0, 1)
	collect2(zipper.Zip(a.seq(), b.seq()), -1)
	if b.pulled != 3 {
		t.Errorf("Zip([1 2 3], infinite): values pulled from the infinite input = %d, want 3", b.pulled)
	}
	checkClosed(t, "Zip([1 2 3], infinite)", a, b)
}

func TestMerge(t *testing.T) {
	tests := []struct {
		a, b []int
		want []int
	}{
		{[]int{1, 4, 9}, []int{2, 3, 10, 11}, []int{1, 2, 3, 4, 9, 10, 11}},
		{[]int{1, 2}, []int{2, 3}, []int{1, 2, 2, 3}},
		{[]int{}, []int{5}, []int{5}},
		{[]int{5}, []int{}, []int{5}},
		{[]int{}, []int{}, nil},
		{[]int{1, 1, 1}, []int{1}, []int{1, 1, 1, 1}},
		{[]int{-3, 0}, []int{-5, -4, 7}, []int{-5, -4, -3, 0, 7}},
		{[]int{10, 20, 30}, []int{1, 2, 3}, []int{1, 2, 3, 10, 20, 30}},
	}
	for _, tt := range tests {
		a, b := list(tt.a...), list(tt.b...)
		seq := zipper.Merge(a.seq(), b.seq())
		if seq == nil {
			t.Fatalf("Merge(%v, %v) = nil, want a non-nil iterator", tt.a, tt.b)
		}
		got, p := collect(seq, -1)
		if p != nil {
			t.Fatalf("range over Merge(%v, %v) panicked: %v", tt.a, tt.b, p)
		}
		if !slices.Equal(got, tt.want) {
			t.Errorf("Merge(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
		checkClosed(t, fmt.Sprintf("Merge(%v, %v)", tt.a, tt.b), a, b)
	}
}

func TestMergeStrings(t *testing.T) {
	a := slices.Values([]string{"apple", "fig", "pear"})
	b := slices.Values([]string{"banana", "fig", "zoo"})
	got, _ := collect(zipper.Merge(a, b), -1)
	if want := []string{"apple", "banana", "fig", "fig", "pear", "zoo"}; !slices.Equal(got, want) {
		t.Errorf("Merge([apple fig pear], [banana fig zoo]) = %q, want %q", got, want)
	}
}

// TestMergeTies checks that a's value comes first on a tie by watching
// which input is pulled from: after yielding the tied value from a, Merge
// reads a's next value, not b's.
func TestMergeTies(t *testing.T) {
	a, b := list(5, 6), list(5, 9)
	got, p := collect(zipper.Merge(a.seq(), b.seq()), 2)
	if p != nil {
		t.Fatalf("break after 2 values of Merge([5 6], [5 9]) panicked: %v", p)
	}
	if want := []int{5, 5}; !slices.Equal(got, want) {
		t.Errorf("first 2 values of Merge([5 6], [5 9]) = %v, want %v", got, want)
	}
	if a.pulled != 2 || b.pulled != 1 {
		t.Errorf("break after 2 values of Merge([5 6], [5 9]): values pulled = %d and %d, want 2 and 1 (on a tie, yield a's value first)", a.pulled, b.pulled)
	}
}

func TestMergeIsLazy(t *testing.T) {
	evens, odds := count(0, 2), count(1, 2)
	seq := zipper.Merge(evens.seq(), odds.seq())
	if evens.started || odds.started {
		t.Errorf("Merge(a, b) before ranging: started a = %v, b = %v, want false, false", evens.started, odds.started)
	}
	got, p := collect(seq, 5)
	if p != nil {
		t.Fatalf("break after 5 values of Merge(evens, odds) panicked: %v (stop when yield returns false)", p)
	}
	if want := []int{0, 1, 2, 3, 4}; !slices.Equal(got, want) {
		t.Errorf("first 5 values of Merge(evens, odds) = %v, want %v", got, want)
	}
	if total := evens.pulled + odds.pulled; total > 7 {
		t.Errorf("break after 5 values of Merge(evens, odds): values pulled = %d, want at most 7", total)
	}
	checkClosed(t, "break in Merge(evens, odds)", evens, odds)

	short, long := list(1, 2), count(0, 10)
	got, p = collect(zipper.Merge(short.seq(), long.seq()), 4)
	if p != nil {
		t.Fatalf("break after 4 values of Merge([1 2], 0, 10, 20, …) panicked: %v", p)
	}
	if want := []int{0, 1, 2, 10}; !slices.Equal(got, want) {
		t.Errorf("first 4 values of Merge([1 2], 0, 10, 20, …) = %v, want %v", got, want)
	}
	checkClosed(t, "break in Merge([1 2], 0, 10, 20, …)", short, long)
}
