package pq

import (
	"slices"
	"testing"
)

// checkHeap fails if q.items is not a valid binary heap for q.less.
func checkHeap[T any](t *testing.T, q *Queue[T]) {
	t.Helper()
	for i := 1; i < len(q.items); i++ {
		if parent := (i - 1) / 2; q.less(q.items[i], q.items[parent]) {
			t.Fatalf("items = %v is not a heap: items[%d] comes before its parent items[%d]", q.items, i, parent)
		}
	}
}

func newInts(t *testing.T) *Queue[int] {
	t.Helper()
	q := NewOrdered[int]()
	if q == nil {
		t.Fatalf("NewOrdered[int]() = nil, want an empty queue")
	}
	return q
}

func TestEmpty(t *testing.T) {
	q := newInts(t)
	if v, ok := q.Pop(); v != 0 || ok {
		t.Errorf("Pop() on empty = %d, %v, want 0, false", v, ok)
	}
	if v, ok := q.Peek(); v != 0 || ok {
		t.Errorf("Peek() on empty = %d, %v, want 0, false", v, ok)
	}
	if got := q.Len(); got != 0 {
		t.Errorf("Len() on empty = %d, want 0", got)
	}
	for v := range q.Drain() {
		t.Errorf("Drain() on empty yielded %d, want nothing", v)
	}
}

func TestOrderedPopsSmallestFirst(t *testing.T) {
	tests := [][]int{
		{5},
		{3, 1, 2},
		{9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
		{4, 4, 1, 4, 1, -7, 0, 12, -7},
		{1, 2, 3, 4, 5, 6, 7, 8},
	}
	for _, in := range tests {
		q := newInts(t)
		q.Push(in...)
		checkHeap(t, q)
		if got := q.Len(); got != len(in) {
			t.Fatalf("Len() after Push(%v) = %d, want %d", in, got, len(in))
		}
		want := slices.Sorted(slices.Values(in))
		if top, ok := q.Peek(); !ok || top != want[0] {
			t.Errorf("Peek() after Push(%v) = %d, %v, want %d, true", in, top, ok, want[0])
		}
		var got []int
		for q.Len() > 0 {
			v, ok := q.Pop()
			if !ok {
				t.Fatalf("Pop() = _, false with Len() > 0")
			}
			checkHeap(t, q)
			got = append(got, v)
		}
		if !slices.Equal(got, want) {
			t.Errorf("Push(%v) then Pop all = %v, want %v", in, got, want)
		}
	}
}

func TestCustomLess(t *testing.T) {
	type task struct {
		name     string
		priority int
	}
	q := New(func(a, b task) bool { return a.priority > b.priority })
	if q == nil {
		t.Fatalf("New(less) = nil, want an empty queue")
	}
	q.Push(task{"write", 2}, task{"deploy", 9}, task{"lunch", 5})
	q.Push(task{"review", 7})
	checkHeap(t, q)
	var got []string
	for tk := range q.Drain() {
		got = append(got, tk.name)
	}
	if want := []string{"deploy", "review", "lunch", "write"}; !slices.Equal(got, want) {
		t.Errorf("Drain() order = %v, want %v", got, want)
	}
	if n := q.Len(); n != 0 {
		t.Errorf("Len() after a full Drain() = %d, want 0", n)
	}
}

func TestStrings(t *testing.T) {
	q := NewOrdered[string]()
	if q == nil {
		t.Fatalf("NewOrdered[string]() = nil, want an empty queue")
	}
	q.Push("pear", "apple", "fig", "Banana")
	if got, want := slices.Collect(q.Drain()), []string{"Banana", "apple", "fig", "pear"}; !slices.Equal(got, want) {
		t.Errorf("Drain() = %q, want %q", got, want)
	}
}

func TestDrainEarlyBreak(t *testing.T) {
	q := newInts(t)
	q.Push(50, 10, 40, 20, 30)
	var got []int
	for v := range q.Drain() {
		got = append(got, v)
		if v == 20 {
			break
		}
	}
	if want := []int{10, 20}; !slices.Equal(got, want) {
		t.Errorf("Drain() with break at 20 = %v, want %v", got, want)
	}
	if n := q.Len(); n != 3 {
		t.Errorf("Len() after the break = %d, want 3", n)
	}
	checkHeap(t, q)
	q.Push(5)
	if got, want := slices.Collect(q.Drain()), []int{5, 30, 40, 50}; !slices.Equal(got, want) {
		t.Errorf("Drain() of the rest = %v, want %v", got, want)
	}
}

func TestInterleaved(t *testing.T) {
	q := newInts(t)
	var model []int
	for i := range 300 {
		if i%3 == 2 {
			v, ok := q.Pop()
			if len(model) == 0 {
				if ok {
					t.Fatalf("step %d: Pop() = %d, true, want _, false", i, v)
				}
				continue
			}
			slices.Sort(model)
			if !ok || v != model[0] {
				t.Fatalf("step %d: Pop() = %d, %v, want %d, true", i, v, ok, model[0])
			}
			model = model[1:]
			continue
		}
		v := (i * 37) % 101
		q.Push(v)
		model = append(model, v)
		checkHeap(t, q)
	}
}

func TestTopK(t *testing.T) {
	tests := []struct {
		in   []int
		k    int
		want []int
	}{
		{[]int{5, 1, 9, 3, 7}, 3, []int{9, 7, 5}},
		{[]int{5, 1, 9, 3, 7}, 1, []int{9}},
		{[]int{2, 8}, 5, []int{8, 2}},
		{[]int{4, 4, 4, 1}, 2, []int{4, 4}},
		{[]int{-3, -1, -2}, 2, []int{-1, -2}},
		{[]int{1, 2, 3}, 0, []int{}},
		{[]int{1, 2, 3}, -1, []int{}},
		{nil, 3, []int{}},
	}
	for _, tt := range tests {
		got := TopK(slices.Values(tt.in), tt.k)
		if got == nil || !slices.Equal(got, tt.want) {
			t.Errorf("TopK(%v, %d) = %#v, want %v", tt.in, tt.k, got, tt.want)
		}
	}
	words := TopK(slices.Values([]string{"b", "d", "a", "c"}), 2)
	if want := []string{"d", "c"}; !slices.Equal(words, want) {
		t.Errorf("TopK([b d a c], 2) = %q, want %q", words, want)
	}
}
