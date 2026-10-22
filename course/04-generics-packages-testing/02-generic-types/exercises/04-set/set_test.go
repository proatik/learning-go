package set

import (
	"cmp"
	"slices"
	"testing"
)

// of builds a set without relying on New or Add.
func of[T comparable](vals ...T) Set[T] {
	s := Set[T]{}
	for _, v := range vals {
		s[v] = struct{}{}
	}
	return s
}

// keys lists a set in ascending order without relying on Sorted.
func keys[T cmp.Ordered](s Set[T]) []T {
	out := []T{}
	for v := range s {
		out = append(out, v)
	}
	slices.Sort(out)
	return out
}

func TestNew(t *testing.T) {
	tests := []struct {
		vals []int
		want []int
	}{
		{[]int{3, 1, 3, 2}, []int{1, 2, 3}},
		{[]int{0}, []int{0}},
		{[]int{-1, -1, -1}, []int{-1}},
		{nil, []int{}},
	}
	for _, tt := range tests {
		s := New(tt.vals...)
		if s == nil {
			t.Errorf("New(%v) = nil, want a non-nil set", tt.vals)
			continue
		}
		if got := keys(s); !slices.Equal(got, tt.want) {
			t.Errorf("New(%v) = %v, want %v", tt.vals, got, tt.want)
		}
	}
	if s := New[string](); s == nil {
		t.Errorf("New[string]() = nil, want a non-nil set")
	}
}

func TestAddHasLen(t *testing.T) {
	s := Set[string]{}
	s.Add("go", "rust")
	s.Add("go")
	s.Add()
	s.Add("世界", "")
	for _, v := range []string{"go", "rust", "世界", ""} {
		if !s.Has(v) {
			t.Errorf("after Add, Has(%q) = false, want true", v)
		}
	}
	if s.Has("Go") {
		t.Errorf("Has(%q) = true, want false", "Go")
	}
	if got := s.Len(); got != 4 {
		t.Errorf("Len() = %d, want 4", got)
	}
}

func TestRemove(t *testing.T) {
	s := of(1, 2, 3)
	if !s.Remove(2) {
		t.Errorf("{1 2 3}.Remove(2) = false, want true")
	}
	if s.Remove(2) {
		t.Errorf("second Remove(2) = true, want false")
	}
	if s.Remove(9) {
		t.Errorf("{1 3}.Remove(9) = true, want false")
	}
	if got, want := keys(s), []int{1, 3}; !slices.Equal(got, want) {
		t.Errorf("after Remove(2), set = %v, want %v", got, want)
	}
}

func TestNilSetReads(t *testing.T) {
	var n Set[int]
	if n.Has(1) {
		t.Errorf("nil.Has(1) = true, want false")
	}
	if got := n.Len(); got != 0 {
		t.Errorf("nil.Len() = %d, want 0", got)
	}
	if !n.SubsetOf(of(1)) {
		t.Errorf("nil.SubsetOf({1}) = false, want true")
	}
	if !n.Equal(Set[int]{}) {
		t.Errorf("nil.Equal({}) = false, want true")
	}
	if u := n.Union(nil); u == nil || len(u) != 0 {
		t.Errorf("nil.Union(nil) = %v, want an empty non-nil set", u)
	}
	if got, want := keys(n.Union(of(2))), []int{2}; !slices.Equal(got, want) {
		t.Errorf("nil.Union({2}) = %v, want %v", got, want)
	}
}

func TestAlgebra(t *testing.T) {
	tests := []struct {
		a, b                    []int
		union, inter, diffAB    []int
		aSubB, bSubA, equalBoth bool
	}{
		{[]int{1, 2, 3}, []int{3, 4}, []int{1, 2, 3, 4}, []int{3}, []int{1, 2}, false, false, false},
		{[]int{1, 2, 3}, []int{2, 3, 4}, []int{1, 2, 3, 4}, []int{2, 3}, []int{1}, false, false, false},
		{[]int{1, 2}, []int{1, 2, 3}, []int{1, 2, 3}, []int{1, 2}, []int{}, true, false, false},
		{[]int{2, 1}, []int{1, 2}, []int{1, 2}, []int{1, 2}, []int{}, true, true, true},
		{[]int{1}, []int{2}, []int{1, 2}, []int{}, []int{1}, false, false, false},
		{nil, []int{5}, []int{5}, []int{}, []int{}, true, false, false},
		{nil, nil, []int{}, []int{}, []int{}, true, true, true},
	}
	for _, tt := range tests {
		a, b := of(tt.a...), of(tt.b...)
		if got := a.Union(b); got == nil || !slices.Equal(keys(got), tt.union) {
			t.Errorf("%v.Union(%v) = %v, want %v", tt.a, tt.b, keys(got), tt.union)
		}
		if got := a.Intersect(b); got == nil || !slices.Equal(keys(got), tt.inter) {
			t.Errorf("%v.Intersect(%v) = %v, want %v", tt.a, tt.b, keys(got), tt.inter)
		}
		if got := a.Difference(b); got == nil || !slices.Equal(keys(got), tt.diffAB) {
			t.Errorf("%v.Difference(%v) = %v, want %v", tt.a, tt.b, keys(got), tt.diffAB)
		}
		if got := a.SubsetOf(b); got != tt.aSubB {
			t.Errorf("%v.SubsetOf(%v) = %v, want %v", tt.a, tt.b, got, tt.aSubB)
		}
		if got := b.SubsetOf(a); got != tt.bSubA {
			t.Errorf("%v.SubsetOf(%v) = %v, want %v", tt.b, tt.a, got, tt.bSubA)
		}
		if got := a.Equal(b); got != tt.equalBoth {
			t.Errorf("%v.Equal(%v) = %v, want %v", tt.a, tt.b, got, tt.equalBoth)
		}
	}
}

func TestOperationsDontModifyInputs(t *testing.T) {
	a, b := of(1, 2, 3), of(3, 4)
	a.Union(b)
	a.Intersect(b)
	a.Difference(b)
	b.Difference(a)
	if got := keys(a); !slices.Equal(got, []int{1, 2, 3}) {
		t.Errorf("set operations changed a: a = %v, want [1 2 3]", got)
	}
	if got := keys(b); !slices.Equal(got, []int{3, 4}) {
		t.Errorf("set operations changed b: b = %v, want [3 4]", got)
	}
	u := a.Union(b)
	if u != nil {
		u[99] = struct{}{}
	}
	if a.Has(99) || b.Has(99) {
		t.Errorf("Union result shares storage with its inputs")
	}
}

func TestStructElements(t *testing.T) {
	type point struct{ X, Y int }
	s := of(point{1, 2}, point{0, 0})
	if !s.Has(point{1, 2}) || s.Has(point{2, 1}) {
		t.Errorf("Set[point] Has({1 2}), Has({2 1}) = %v, %v, want true, false", s.Has(point{1, 2}), s.Has(point{2, 1}))
	}
}

func TestSorted(t *testing.T) {
	if got, want := Sorted(of("pear", "Apple", "fig")), []string{"Apple", "fig", "pear"}; !slices.Equal(got, want) {
		t.Errorf("Sorted({pear Apple fig}) = %q, want %q", got, want)
	}
	if got, want := Sorted(of(3.5, -1, 0)), []float64{-1, 0, 3.5}; !slices.Equal(got, want) {
		t.Errorf("Sorted({3.5 -1 0}) = %v, want %v", got, want)
	}
	if got := Sorted(Set[int](nil)); len(got) != 0 {
		t.Errorf("Sorted(nil) = %v, want []", got)
	}
}
