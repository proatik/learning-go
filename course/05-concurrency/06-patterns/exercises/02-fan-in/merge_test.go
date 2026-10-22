package fanin

import (
	"slices"
	"testing"
	"time"
)

// drain reads ch until it is closed, failing if that takes over 2s.
func drain[T any](t *testing.T, name string, ch <-chan T) []T {
	t.Helper()
	if ch == nil {
		t.Fatalf("%s returned a nil channel", name)
	}
	var got []T
	timeout := time.After(2 * time.Second)
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return got
			}
			got = append(got, v)
		case <-timeout:
			t.Fatalf("%s: output not closed within 2s (received %d values so far)", name, len(got))
		}
	}
}

// closed returns a closed channel holding vals.
func closed[T any](vals ...T) <-chan T {
	ch := make(chan T, len(vals))
	for _, v := range vals {
		ch <- v
	}
	close(ch)
	return ch
}

func TestMerge(t *testing.T) {
	tests := []struct {
		name   string
		inputs [][]int
		want   []int // sorted
	}{
		{"three inputs", [][]int{{1, 2, 3}, {10, 20}, {100}}, []int{1, 2, 3, 10, 20, 100}},
		{"one input", [][]int{{5, 4, 3}}, []int{3, 4, 5}},
		{"some empty", [][]int{{}, {7}, {}}, []int{7}},
		{"all empty", [][]int{{}, {}}, nil},
		{"duplicates", [][]int{{1, 1}, {1}}, []int{1, 1, 1}},
	}
	for _, tt := range tests {
		var ins []<-chan int
		for _, vals := range tt.inputs {
			ins = append(ins, closed(vals...))
		}
		got := drain(t, "Merge("+tt.name+")", Merge(ins...))
		slices.Sort(got)
		if !slices.Equal(got, tt.want) {
			t.Errorf("Merge(%v) = %v (sorted), want %v", tt.inputs, got, tt.want)
		}
	}
}

func TestMergeNoInputs(t *testing.T) {
	if got := drain(t, "Merge()", Merge[int]()); len(got) != 0 {
		t.Errorf("Merge() = %v, want nothing", got)
	}
}

func TestMergeNilInputs(t *testing.T) {
	got := drain(t, "Merge(nil, [a b], nil)", Merge(nil, closed("a", "b"), nil))
	slices.Sort(got)
	if want := []string{"a", "b"}; !slices.Equal(got, want) {
		t.Errorf("Merge(nil, [a b], nil) = %q, want %q", got, want)
	}
	if got := drain(t, "Merge(nil, nil)", Merge[string](nil, nil)); len(got) != 0 {
		t.Errorf("Merge(nil, nil) = %q, want nothing", got)
	}
}

// TestMergeKeepsPerInputOrder: values from the same input arrive in the order
// that input sent them.
func TestMergeKeepsPerInputOrder(t *testing.T) {
	a := closed(1, 2, 3, 4, 5)
	b := closed(-1, -2, -3, -4, -5)
	got := drain(t, "Merge", Merge(a, b))
	var pos, neg []int
	for _, v := range got {
		if v > 0 {
			pos = append(pos, v)
		} else {
			neg = append(neg, v)
		}
	}
	if !slices.Equal(pos, []int{1, 2, 3, 4, 5}) || !slices.Equal(neg, []int{-1, -2, -3, -4, -5}) {
		t.Errorf("Merge([1..5], [-1..-5]) = %v, want each input's values in their own order", got)
	}
}

// TestMergeReadsInputsConcurrently: input a sends its second value only after
// b's value has come out of Merge, so reading one input at a time stalls.
func TestMergeReadsInputsConcurrently(t *testing.T) {
	a, b := make(chan int), make(chan int)
	gotB := make(chan struct{})
	go func() {
		defer close(a)
		a <- 1
		select {
		case <-gotB:
			a <- 2
		case <-time.After(2 * time.Second):
		}
	}()
	go func() {
		defer close(b)
		b <- 100
	}()
	out := Merge[int](a, b)
	if out == nil {
		t.Fatalf("Merge returned a nil channel")
	}
	var got []int
	timeout := time.After(2 * time.Second)
	for len(got) < 3 {
		select {
		case v, ok := <-out:
			if !ok {
				t.Fatalf("Merge: output closed after %v, want 3 values", got)
			}
			got = append(got, v)
			if v == 100 {
				close(gotB)
			}
		case <-timeout:
			t.Fatalf("Merge: got %v within 2s, want 3 values (read every input in its own goroutine)", got)
		}
	}
	slices.Sort(got)
	if want := []int{1, 2, 100}; !slices.Equal(got, want) {
		t.Errorf("Merge = %v (sorted), want %v", got, want)
	}
}

func TestMergeMany(t *testing.T) {
	const inputs, per = 50, 1000
	var ins []<-chan int
	for i := range inputs {
		vals := make([]int, per)
		for j := range vals {
			vals[j] = i*per + j
		}
		ins = append(ins, closed(vals...))
	}
	got := drain(t, "Merge(50 inputs)", Merge(ins...))
	if len(got) != inputs*per {
		t.Fatalf("Merge(50 inputs × 1000) yielded %d values, want %d", len(got), inputs*per)
	}
	slices.Sort(got)
	for i, v := range got {
		if v != i {
			t.Fatalf("Merge(50 inputs × 1000): value %d missing or duplicated", i)
		}
	}
}
