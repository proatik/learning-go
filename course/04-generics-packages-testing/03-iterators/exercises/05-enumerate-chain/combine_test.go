package combine

import (
	"fmt"
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

// tracked yields vals and records whether it was started and how many
// values it produced.
type tracked struct {
	vals    []string
	started bool
	pulled  int
}

func (tr *tracked) seq() iter.Seq[string] {
	return func(yield func(string) bool) {
		tr.started = true
		for _, v := range tr.vals {
			tr.pulled++
			if !yield(v) {
				return
			}
		}
	}
}

func TestEnumerate(t *testing.T) {
	tests := []struct {
		in   []string
		want []string
	}{
		{[]string{"a", "b"}, []string{"0:a", "1:b"}},
		{[]string{"x"}, []string{"0:x"}},
		{nil, nil},
		{[]string{"é", "", "z"}, []string{"0:é", "1:", "2:z"}},
	}
	for _, tt := range tests {
		seq := Enumerate(slices.Values(tt.in))
		if seq == nil {
			t.Fatalf("Enumerate(%q) = nil, want a non-nil iterator", tt.in)
		}
		got, p := collect2(seq, -1)
		if p != nil {
			t.Fatalf("range over Enumerate(%q) panicked: %v", tt.in, p)
		}
		if !slices.Equal(got, tt.want) {
			t.Errorf("Enumerate(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestEnumerateBreak(t *testing.T) {
	src := &tracked{vals: []string{"a", "b", "c", "d"}}
	got, p := collect2(Enumerate(src.seq()), 2)
	if p != nil {
		t.Fatalf("break after 2 pairs of Enumerate([a b c d]) panicked: %v (stop when yield returns false)", p)
	}
	if want := []string{"0:a", "1:b"}; !slices.Equal(got, want) {
		t.Errorf("first 2 pairs of Enumerate([a b c d]) = %q, want %q", got, want)
	}
	if src.pulled != 2 {
		t.Errorf("break after 2 pairs of Enumerate([a b c d]): source values produced = %d, want 2", src.pulled)
	}
}

func TestChain(t *testing.T) {
	tests := []struct {
		parts [][]string
		want  []string
	}{
		{[][]string{{"1", "2"}, {"3"}}, []string{"1", "2", "3"}},
		{[][]string{{"a"}, {}, {"b", "c"}, nil}, []string{"a", "b", "c"}},
		{[][]string{{"solo"}}, []string{"solo"}},
		{nil, nil},
	}
	for _, tt := range tests {
		var seqs []iter.Seq[string]
		for _, p := range tt.parts {
			seqs = append(seqs, slices.Values(p))
		}
		seq := Chain(seqs...)
		if seq == nil {
			t.Fatalf("Chain(%q) = nil, want a non-nil iterator", tt.parts)
		}
		got, p := collect(seq, -1)
		if p != nil {
			t.Fatalf("range over Chain(%q) panicked: %v", tt.parts, p)
		}
		if !slices.Equal(got, tt.want) {
			t.Errorf("Chain(%q) = %q, want %q", tt.parts, got, tt.want)
		}
	}
}

func TestChainSkipsNil(t *testing.T) {
	got, p := collect(Chain(nil, slices.Values([]int{1}), nil, slices.Values([]int{2})), -1)
	if p != nil {
		t.Fatalf("range over Chain(nil, [1], nil, [2]) panicked: %v (skip nil sequences)", p)
	}
	if want := []int{1, 2}; !slices.Equal(got, want) {
		t.Errorf("Chain(nil, [1], nil, [2]) = %v, want %v", got, want)
	}
}

func TestChainBreak(t *testing.T) {
	a := &tracked{vals: []string{"a1", "a2", "a3"}}
	b := &tracked{vals: []string{"b1", "b2"}}
	seq := Chain(a.seq(), b.seq())
	if a.started || b.started {
		t.Errorf("Chain(a, b) before ranging: started a = %v, b = %v, want false, false", a.started, b.started)
	}
	got, p := collect(seq, 2)
	if p != nil {
		t.Fatalf("break after 2 values of Chain([a1 a2 a3], [b1 b2]) panicked: %v (return, don't break, when yield returns false)", p)
	}
	if want := []string{"a1", "a2"}; !slices.Equal(got, want) {
		t.Errorf("first 2 values of Chain([a1 a2 a3], [b1 b2]) = %q, want %q", got, want)
	}
	if a.pulled != 2 || b.started {
		t.Errorf("break after 2 values of Chain([a1 a2 a3], [b1 b2]): values pulled from a = %d, b started = %v, want 2, false", a.pulled, b.started)
	}

	a = &tracked{vals: []string{"a1"}}
	b = &tracked{vals: []string{"b1", "b2", "b3"}}
	got, p = collect(Chain(a.seq(), b.seq()), 2)
	if p != nil {
		t.Fatalf("break after 2 values of Chain([a1], [b1 b2 b3]) panicked: %v", p)
	}
	if want := []string{"a1", "b1"}; !slices.Equal(got, want) || b.pulled != 1 {
		t.Errorf("first 2 values of Chain([a1], [b1 b2 b3]) = %q (b pulled %d), want %q (b pulled 1)", got, b.pulled, want)
	}
}

func TestKeysValues(t *testing.T) {
	words := []string{"x", "y", "z"}
	keys, p := collect(Keys(slices.All(words)), -1)
	if p != nil {
		t.Fatalf("range over Keys(slices.All(%q)) panicked: %v", words, p)
	}
	if want := []int{0, 1, 2}; !slices.Equal(keys, want) {
		t.Errorf("Keys(slices.All(%q)) = %v, want %v", words, keys, want)
	}
	vals, p := collect(Values(slices.All(words)), -1)
	if p != nil {
		t.Fatalf("range over Values(slices.All(%q)) panicked: %v", words, p)
	}
	if !slices.Equal(vals, words) {
		t.Errorf("Values(slices.All(%q)) = %q, want %q", words, vals, words)
	}
	none, _ := collect(Keys(slices.All([]int(nil))), -1)
	if len(none) != 0 {
		t.Errorf("Keys(slices.All([])) = %v, want []", none)
	}
	round, _ := collect(Values(Enumerate(slices.Values(words))), -1)
	if !slices.Equal(round, words) {
		t.Errorf("Values(Enumerate(%q)) = %q, want %q", words, round, words)
	}
}

func TestKeysValuesBreak(t *testing.T) {
	src := &tracked{vals: []string{"a", "b", "c"}}
	got, p := collect(Keys(Enumerate(src.seq())), 1)
	if p != nil {
		t.Fatalf("break after 1 value of Keys(Enumerate([a b c])) panicked: %v (stop when yield returns false)", p)
	}
	if !slices.Equal(got, []int{0}) || src.pulled != 1 {
		t.Errorf("break after 1 value of Keys(Enumerate([a b c])): source values produced = %d, want 1", src.pulled)
	}
	src = &tracked{vals: []string{"a", "b", "c"}}
	vals, p := collect(Values(Enumerate(src.seq())), 2)
	if p != nil {
		t.Fatalf("break after 2 values of Values(Enumerate([a b c])) panicked: %v (stop when yield returns false)", p)
	}
	if !slices.Equal(vals, []string{"a", "b"}) || src.pulled != 2 {
		t.Errorf("break after 2 values of Values(Enumerate([a b c])): source values produced = %d, want 2", src.pulled)
	}
}
