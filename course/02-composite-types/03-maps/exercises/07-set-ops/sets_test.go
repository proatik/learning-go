package sets

import (
	"maps"
	"slices"
	"testing"
)

// mk builds a Set without relying on NewSet.
func mk(items ...string) Set {
	s := Set{}
	for _, it := range items {
		s[it] = struct{}{}
	}
	return s
}

func TestNewSet(t *testing.T) {
	tests := []struct {
		items []string
		want  Set
	}{
		{[]string{"go", "rust", "zig"}, mk("go", "rust", "zig")},
		{[]string{"go", "c", "zig", "c"}, mk("go", "c", "zig")},
		{[]string{"x", "x", "x"}, mk("x")},
		{nil, mk()},
	}
	for _, tt := range tests {
		got := NewSet(tt.items...)
		if got == nil {
			t.Errorf("NewSet(%q...) = nil, want an initialized set", tt.items)
			continue
		}
		if !maps.Equal(got, tt.want) {
			t.Errorf("NewSet(%q...) = %v, want %v", tt.items, got, tt.want)
		}
	}
}

type binCase struct {
	a, b Set
	want []string
}

func checkBinary(t *testing.T, name string, f func(a, b Set) []string, tests []binCase) {
	t.Helper()
	for _, tt := range tests {
		a, b := maps.Clone(tt.a), maps.Clone(tt.b)
		got := f(a, b)
		if !slices.Equal(got, tt.want) {
			t.Errorf("%s(%v, %v) = %q, want %q", name, tt.a, tt.b, got, tt.want)
		}
		if !maps.Equal(a, tt.a) || !maps.Equal(b, tt.b) {
			t.Errorf("%s modified its inputs: (%v, %v) became (%v, %v)", name, tt.a, tt.b, a, b)
		}
	}
}

var (
	langsA = mk("go", "rust", "zig")
	langsB = mk("go", "c", "zig")
	many   = mk("m", "e", "x", "a", "q", "b", "z", "k")
)

func TestUnion(t *testing.T) {
	checkBinary(t, "Union", Union, []binCase{
		{langsA, langsB, []string{"c", "go", "rust", "zig"}},
		{langsB, langsA, []string{"c", "go", "rust", "zig"}},
		{langsA, langsA, []string{"go", "rust", "zig"}},
		{many, mk("a", "y"), []string{"a", "b", "e", "k", "m", "q", "x", "y", "z"}},
		{langsA, nil, []string{"go", "rust", "zig"}},
		{nil, langsB, []string{"c", "go", "zig"}},
		{nil, nil, nil},
		{mk(), mk(), nil},
	})
}

func TestIntersection(t *testing.T) {
	checkBinary(t, "Intersection", Intersection, []binCase{
		{langsA, langsB, []string{"go", "zig"}},
		{langsB, langsA, []string{"go", "zig"}},
		{many, mk("z", "a", "nope", "m"), []string{"a", "m", "z"}},
		{langsA, mk("c"), nil},
		{langsA, nil, nil},
		{nil, langsA, nil},
	})
}

func TestDifference(t *testing.T) {
	checkBinary(t, "Difference", Difference, []binCase{
		{langsA, langsB, []string{"rust"}},
		{langsB, langsA, []string{"c"}},
		{many, mk("a", "b"), []string{"e", "k", "m", "q", "x", "z"}},
		{langsA, langsA, nil},
		{langsA, nil, []string{"go", "rust", "zig"}},
		{nil, langsA, nil},
	})
}

func TestIsSubset(t *testing.T) {
	tests := []struct {
		a, b Set
		want bool
	}{
		{mk("go"), langsA, true},
		{langsA, langsA, true},
		{langsA, langsB, false},
		{mk("go", "zig"), langsB, true},
		{mk("go", "zig", "java"), langsB, false},
		{mk(), langsA, true},
		{nil, langsA, true},
		{nil, nil, true},
		{mk("go"), nil, false},
	}
	for _, tt := range tests {
		if got := IsSubset(tt.a, tt.b); got != tt.want {
			t.Errorf("IsSubset(%v, %v) = %t, want %t", tt.a, tt.b, got, tt.want)
		}
	}
}
