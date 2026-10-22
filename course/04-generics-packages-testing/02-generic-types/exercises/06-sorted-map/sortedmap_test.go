package sortedmap_test

import (
	"slices"
	"testing"

	sm "learngo/course/04-generics-packages-testing/02-generic-types/exercises/06-sorted-map"
)

type entry = sm.Entry[string, int]

// build puts every key/value pair, in the given order.
func build(keys []string, vals []int) *sm.Map[string, int] {
	m := &sm.Map[string, int]{}
	for i, k := range keys {
		m.Put(k, vals[i])
	}
	return m
}

func TestZeroValue(t *testing.T) {
	var m sm.Map[string, int]
	if got := m.Len(); got != 0 {
		t.Errorf("empty Map Len() = %d, want 0", got)
	}
	if v, ok := m.Get("a"); v != 0 || ok {
		t.Errorf("empty Map Get(%q) = (%d, %v), want (0, false)", "a", v, ok)
	}
	if m.Delete("a") {
		t.Errorf("empty Map Delete(%q) = true, want false", "a")
	}
	if got := m.Keys(); len(got) != 0 {
		t.Errorf("empty Map Keys() = %q, want []", got)
	}
	if got := m.Entries(); len(got) != 0 {
		t.Errorf("empty Map Entries() = %v, want []", got)
	}
	if e, ok := m.Min(); e != (entry{}) || ok {
		t.Errorf("empty Map Min() = (%v, %v), want ({ 0}, false)", e, ok)
	}
}

func TestPutKeepsKeysSorted(t *testing.T) {
	tests := []struct {
		keys []string
		want []string
	}{
		{[]string{"b", "a", "c"}, []string{"a", "b", "c"}},
		{[]string{"c", "b", "a"}, []string{"a", "b", "c"}},
		{[]string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{[]string{"b", "b", "a", "b"}, []string{"a", "b"}},
		{[]string{"pear", "Apple", "", "fig", "été"}, []string{"", "Apple", "fig", "pear", "été"}},
	}
	for _, tt := range tests {
		m := build(tt.keys, make([]int, len(tt.keys)))
		if got := m.Keys(); !slices.Equal(got, tt.want) {
			t.Errorf("Put %q: Keys() = %q, want %q", tt.keys, got, tt.want)
		}
		if got := m.Len(); got != len(tt.want) {
			t.Errorf("Put %q: Len() = %d, want %d", tt.keys, got, len(tt.want))
		}
	}
}

func TestGetAndOverwrite(t *testing.T) {
	m := build([]string{"b", "a", "c"}, []int{2, 1, 3})
	m.Put("b", 20)
	tests := []struct {
		key  string
		want int
		ok   bool
	}{
		{"a", 1, true},
		{"b", 20, true},
		{"c", 3, true},
		{"zz", 0, false},
		{"", 0, false},
	}
	for _, tt := range tests {
		if v, ok := m.Get(tt.key); v != tt.want || ok != tt.ok {
			t.Errorf("Get(%q) = (%d, %v), want (%d, %v)", tt.key, v, ok, tt.want, tt.ok)
		}
	}
	if got := m.Len(); got != 3 {
		t.Errorf("after overwriting \"b\", Len() = %d, want 3", got)
	}
	m.Put("zero", 0)
	if v, ok := m.Get("zero"); v != 0 || !ok {
		t.Errorf("Get(%q) = (%d, %v), want (0, true)", "zero", v, ok)
	}
}

func TestEntries(t *testing.T) {
	m := build([]string{"c", "a", "b"}, []int{3, 1, 2})
	want := []entry{{"a", 1}, {"b", 2}, {"c", 3}}
	if got := m.Entries(); !slices.Equal(got, want) {
		t.Errorf("Entries() = %v, want %v", got, want)
	}
}

func TestDelete(t *testing.T) {
	m := build([]string{"a", "b", "c", "d"}, []int{1, 2, 3, 4})
	tests := []struct {
		key      string
		want     bool
		keysLeft []string
	}{
		{"b", true, []string{"a", "c", "d"}},
		{"b", false, []string{"a", "c", "d"}},
		{"a", true, []string{"c", "d"}},
		{"d", true, []string{"c"}},
		{"x", false, []string{"c"}},
		{"c", true, []string{}},
	}
	for _, tt := range tests {
		if got := m.Delete(tt.key); got != tt.want {
			t.Errorf("Delete(%q) = %v, want %v", tt.key, got, tt.want)
		}
		if got := m.Keys(); !slices.Equal(got, tt.keysLeft) {
			t.Errorf("after Delete(%q), Keys() = %q, want %q", tt.key, got, tt.keysLeft)
		}
		if _, ok := m.Get(tt.key); ok {
			t.Errorf("after Delete(%q), Get(%q) found it, want missing", tt.key, tt.key)
		}
	}
	m.Put("b", 7)
	if v, ok := m.Get("b"); v != 7 || !ok {
		t.Errorf("Put after Delete: Get(%q) = (%d, %v), want (7, true)", "b", v, ok)
	}
}

func TestRange(t *testing.T) {
	m := build([]string{"a", "c", "e", "g"}, []int{1, 3, 5, 7})
	tests := []struct {
		lo, hi string
		want   []entry
	}{
		{"a", "z", []entry{{"a", 1}, {"c", 3}, {"e", 5}, {"g", 7}}},
		{"c", "g", []entry{{"c", 3}, {"e", 5}}},
		{"b", "f", []entry{{"c", 3}, {"e", 5}}},
		{"c", "d", []entry{{"c", 3}}},
		{"", "a", nil},
		{"h", "z", nil},
		{"e", "e", nil},
		{"g", "a", nil},
	}
	for _, tt := range tests {
		got := m.Range(tt.lo, tt.hi)
		if len(got) != len(tt.want) || (len(got) > 0 && !slices.Equal(got, tt.want)) {
			t.Errorf("Range(%q, %q) = %v, want %v", tt.lo, tt.hi, got, tt.want)
		}
	}
}

func TestResultsAreCopies(t *testing.T) {
	m := build([]string{"a", "b"}, []int{1, 2})
	keys := m.Keys()
	if len(keys) != 2 {
		t.Fatalf("Keys() = %q, want [a b]", keys)
	}
	keys[0] = "zzz"
	if got := m.Keys(); !slices.Equal(got, []string{"a", "b"}) {
		t.Errorf("modifying the Keys() result changed the map: Keys() = %q, want [a b]", got)
	}
}

func TestMin(t *testing.T) {
	m := build([]string{"m", "b", "x"}, []int{1, 2, 3})
	if e, ok := m.Min(); e != (entry{"b", 2}) || !ok {
		t.Errorf("Min() = (%v, %v), want ({b 2}, true)", e, ok)
	}
	m.Delete("b")
	if e, ok := m.Min(); e != (entry{"m", 1}) || !ok {
		t.Errorf("after Delete(\"b\"), Min() = (%v, %v), want ({m 1}, true)", e, ok)
	}
}

func TestOtherTypes(t *testing.T) {
	var m sm.Map[float64, []string]
	m.Put(2.5, []string{"x"})
	m.Put(-1, nil)
	m.Put(0, []string{"a", "b"})
	if got, want := m.Keys(), []float64{-1, 0, 2.5}; !slices.Equal(got, want) {
		t.Errorf("Map[float64, []string] Keys() = %v, want %v", got, want)
	}
	if v, ok := m.Get(0); !slices.Equal(v, []string{"a", "b"}) || !ok {
		t.Errorf("Get(0) = (%q, %v), want ([a b], true)", v, ok)
	}
}
