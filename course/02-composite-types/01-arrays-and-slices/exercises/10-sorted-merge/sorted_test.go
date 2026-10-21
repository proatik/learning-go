package sorted

import (
	"slices"
	"testing"
)

func TestMergeSorted(t *testing.T) {
	tests := []struct {
		a, b []int
		want []int
	}{
		{[]int{1, 4, 9}, []int{2, 3, 10}, []int{1, 2, 3, 4, 9, 10}},
		{[]int{1, 1}, []int{1}, []int{1, 1, 1}},
		{nil, []int{5}, []int{5}},
		{[]int{5}, nil, []int{5}},
		{[]int{1, 2, 3}, []int{4, 5, 6}, []int{1, 2, 3, 4, 5, 6}},
		{[]int{4, 5, 6}, []int{1, 2, 3}, []int{1, 2, 3, 4, 5, 6}},
		{[]int{-3, 0, 0, 7}, []int{-5, 0, 8, 8}, []int{-5, -3, 0, 0, 0, 7, 8, 8}},
		{[]int{2}, []int{1, 3, 5, 7}, []int{1, 2, 3, 5, 7}},
	}
	for _, tt := range tests {
		a, b := slices.Clone(tt.a), slices.Clone(tt.b)
		got := MergeSorted(a, b)
		if !slices.Equal(got, tt.want) {
			t.Errorf("MergeSorted(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
		if !slices.Equal(a, tt.a) || !slices.Equal(b, tt.b) {
			t.Errorf("MergeSorted(%v, %v) modified its inputs to %v, %v", tt.a, tt.b, a, b)
		}
	}
}

func TestMergeSortedEmpty(t *testing.T) {
	if got := MergeSorted(nil, nil); len(got) != 0 {
		t.Errorf("MergeSorted(nil, nil) = %v, want an empty slice", got)
	}
	if got := MergeSorted([]int{}, []int{}); len(got) != 0 {
		t.Errorf("MergeSorted([], []) = %v, want an empty slice", got)
	}
}

func TestMergeSortedIndependent(t *testing.T) {
	// a has spare capacity: a naive append(a, b...) would write into it.
	backing := []int{1, 3, 100, 100, 100}
	a := backing[:2]
	b := []int{2}
	got := MergeSorted(a, b)
	if !slices.Equal(got, []int{1, 2, 3}) {
		t.Fatalf("MergeSorted([1 3], [2]) = %v, want [1 2 3]", got)
	}
	if !slices.Equal(backing, []int{1, 3, 100, 100, 100}) {
		t.Errorf("MergeSorted wrote into a's backing array: %v", backing)
	}
	got[0] = -1
	if a[0] != 1 {
		t.Errorf("MergeSorted result shares memory with a")
	}
}

func TestDedupe(t *testing.T) {
	tests := []struct {
		in   []int
		want []int
	}{
		{[]int{1, 1, 2, 3, 3, 3}, []int{1, 2, 3}},
		{[]int{4}, []int{4}},
		{[]int{1, 2, 3}, []int{1, 2, 3}},
		{[]int{7, 7, 7, 7}, []int{7}},
		{[]int{-2, -2, 0, 5, 5}, []int{-2, 0, 5}},
		{[]int{0, 0, 1}, []int{0, 1}},
		{[]int{}, []int{}},
		{nil, nil},
	}
	for _, tt := range tests {
		got := Dedupe(slices.Clone(tt.in))
		if !slices.Equal(got, tt.want) {
			t.Errorf("Dedupe(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestDedupeInPlace(t *testing.T) {
	s := []int{1, 1, 2, 2, 3}
	got := Dedupe(s)
	if len(got) != 3 {
		t.Fatalf("Dedupe([1 1 2 2 3]) = %v, want [1 2 3]", got)
	}
	if &got[0] != &s[0] {
		t.Errorf("Dedupe allocated a new array; it must work in place")
	}
	data := make([]int, 200)
	work := make([]int, 200)
	for i := range data {
		data[i] = i / 3
	}
	allocs := testing.AllocsPerRun(100, func() {
		copy(work, data)
		_ = Dedupe(work)
	})
	if allocs != 0 {
		t.Errorf("Dedupe allocated %.0f times per run, want 0", allocs)
	}
}

func TestIntersect(t *testing.T) {
	tests := []struct {
		a, b []int
		want []int
	}{
		{[]int{1, 2, 2, 3, 5}, []int{2, 2, 5, 6}, []int{2, 5}},
		{[]int{1, 2, 3}, []int{1, 2, 3}, []int{1, 2, 3}},
		{[]int{1, 1, 1}, []int{1}, []int{1}},
		{[]int{-4, 0, 4, 8}, []int{-8, -4, 4, 16}, []int{-4, 4}},
		{[]int{5}, []int{1, 2, 3, 4, 5}, []int{5}},
		{[]int{1, 3, 5, 7}, []int{2, 3, 4, 7, 9}, []int{3, 7}},
	}
	for _, tt := range tests {
		a, b := slices.Clone(tt.a), slices.Clone(tt.b)
		got := Intersect(a, b)
		if !slices.Equal(got, tt.want) {
			t.Errorf("Intersect(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
		if !slices.Equal(a, tt.a) || !slices.Equal(b, tt.b) {
			t.Errorf("Intersect(%v, %v) modified its inputs to %v, %v", tt.a, tt.b, a, b)
		}
	}
}

func TestIntersectEmpty(t *testing.T) {
	tests := []struct{ a, b []int }{
		{[]int{1, 2}, []int{3, 4}},
		{nil, []int{1}},
		{[]int{1}, nil},
		{nil, nil},
	}
	for _, tt := range tests {
		if got := Intersect(tt.a, tt.b); len(got) != 0 {
			t.Errorf("Intersect(%v, %v) = %v, want an empty slice", tt.a, tt.b, got)
		}
	}
}

func TestIntersectIndependent(t *testing.T) {
	backing := []int{2, 4, 100, 100}
	a := backing[:2] // spare capacity
	got := Intersect(a, []int{2, 4})
	if !slices.Equal(got, []int{2, 4}) {
		t.Fatalf("Intersect([2 4], [2 4]) = %v, want [2 4]", got)
	}
	got[0] = -1
	if a[0] != 2 || backing[2] != 100 {
		t.Errorf("Intersect result shares memory with its input: backing = %v", backing)
	}
}
