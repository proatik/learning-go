package omap

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

func newMap(t *testing.T) *OrderedMap[string, int] {
	t.Helper()
	m := New[string, int]()
	if m == nil {
		t.Fatalf("New[string, int]() = nil, want an empty map")
	}
	return m
}

// pairs renders All as "k=v" strings, oldest first.
func pairs[K comparable, V any](m *OrderedMap[K, V]) []string {
	out := []string{}
	for k, v := range m.All() {
		out = append(out, fmt.Sprintf("%v=%v", k, v))
	}
	return out
}

func check(t *testing.T, op string, m *OrderedMap[string, int], want ...string) {
	t.Helper()
	if want == nil {
		want = []string{}
	}
	if got := pairs(m); !slices.Equal(got, want) {
		t.Fatalf("after %s: All() = %v, want %v", op, got, want)
	}
	var wantKeys []string
	for _, p := range want {
		k, _, _ := strings.Cut(p, "=")
		wantKeys = append(wantKeys, k)
	}
	if got := slices.Collect(m.Keys()); !slices.Equal(got, wantKeys) {
		t.Fatalf("after %s: Keys() = %v, want %v", op, got, wantKeys)
	}
	var back []string
	for k, v := range m.Backward() {
		back = append(back, fmt.Sprintf("%v=%v", k, v))
	}
	wantBack := slices.Clone(want)
	slices.Reverse(wantBack)
	if !slices.Equal(back, wantBack) {
		t.Fatalf("after %s: Backward() = %v, want %v", op, back, wantBack)
	}
	if got := m.Len(); got != len(want) {
		t.Fatalf("after %s: Len() = %d, want %d", op, got, len(want))
	}
}

func TestEmpty(t *testing.T) {
	m := newMap(t)
	check(t, "New", m)
	if v, ok := m.Get("x"); v != 0 || ok {
		t.Errorf("Get(%q) = %d, %v, want 0, false", "x", v, ok)
	}
	if m.Delete("x") {
		t.Errorf("Delete(%q) = true, want false", "x")
	}
}

func TestSetKeepsInsertionOrder(t *testing.T) {
	m := newMap(t)
	m.Set("b", 2)
	m.Set("a", 1)
	m.Set("c", 3)
	check(t, "Set b, a, c", m, "b=2", "a=1", "c=3")

	m.Set("a", 10)
	check(t, `Set("a", 10)`, m, "b=2", "a=10", "c=3")
	if v, ok := m.Get("a"); v != 10 || !ok {
		t.Errorf("Get(%q) = %d, %v, want 10, true", "a", v, ok)
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name string
		del  string
		want []string
	}{
		{"head", "a", []string{"b=2", "c=3", "d=4"}},
		{"middle", "c", []string{"a=1", "b=2", "d=4"}},
		{"tail", "d", []string{"a=1", "b=2", "c=3"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := newMap(t)
			for i, k := range []string{"a", "b", "c", "d"} {
				m.Set(k, i+1)
			}
			if !m.Delete(tt.del) {
				t.Fatalf("Delete(%q) = false, want true", tt.del)
			}
			check(t, "Delete("+tt.del+")", m, tt.want...)
			if m.Delete(tt.del) {
				t.Errorf("second Delete(%q) = true, want false", tt.del)
			}
			m.Set(tt.del, 9)
			check(t, "re-adding "+tt.del, m, append(slices.Clone(tt.want), tt.del+"=9")...)
		})
	}
}

func TestDeleteAll(t *testing.T) {
	m := newMap(t)
	m.Set("x", 1)
	m.Set("y", 2)
	m.Delete("x")
	m.Delete("y")
	check(t, "deleting every key", m)
	m.Set("z", 3)
	check(t, `Set("z", 3) on the emptied map`, m, "z=3")
}

func TestDeleteWhileRanging(t *testing.T) {
	m := newMap(t)
	for i := range 6 {
		m.Set(fmt.Sprint("k", i), i)
	}
	var seen []string
	for k, v := range m.All() {
		seen = append(seen, k)
		if v%2 == 0 {
			m.Delete(k)
		}
	}
	want := []string{"k0", "k1", "k2", "k3", "k4", "k5"}
	if !slices.Equal(seen, want) {
		t.Fatalf("keys seen while deleting = %v, want %v", seen, want)
	}
	check(t, "deleting even values while ranging", m, "k1=1", "k3=3", "k5=5")
}

func TestEarlyBreak(t *testing.T) {
	m := newMap(t)
	for i := range 5 {
		m.Set(fmt.Sprint(i), i)
	}
	var got []string
	for k := range m.Keys() {
		if k == "2" {
			break
		}
		got = append(got, k)
	}
	if want := []string{"0", "1"}; !slices.Equal(got, want) {
		t.Errorf("Keys() with break at %q = %v, want %v", "2", got, want)
	}
	n := 0
	for range m.All() {
		n++
		break
	}
	for range m.Backward() {
		n++
		break
	}
	if n != 2 {
		t.Errorf("break after the first pair ran the body %d times, want 2", n)
	}
}

func TestCollect(t *testing.T) {
	m := Collect(slices.All([]string{"x", "y", "z"}))
	if m == nil {
		t.Fatalf("Collect(slices.All(...)) = nil, want a map")
	}
	got := pairs(m)
	if want := []string{"0=x", "1=y", "2=z"}; !slices.Equal(got, want) {
		t.Errorf("Collect(slices.All([x y z])) pairs = %v, want %v", got, want)
	}

	words := New[string, int]()
	if words == nil {
		t.Fatalf("New[string, int]() = nil, want an empty map")
	}
	words.Set("go", 1)
	words.Set("is", 2)
	words.Set("go", 3)
	again := Collect(words.All())
	if again == nil {
		t.Fatalf("Collect(m.All()) = nil, want a map")
	}
	if got, want := pairs(again), []string{"go=3", "is=2"}; !slices.Equal(got, want) {
		t.Errorf("Collect(m.All()) pairs = %v, want %v", got, want)
	}
}

func TestStructKeys(t *testing.T) {
	type point struct{ X, Y int }
	m := New[point, string]()
	if m == nil {
		t.Fatalf("New[point, string]() = nil, want an empty map")
	}
	m.Set(point{1, 2}, "a")
	m.Set(point{0, 0}, "origin")
	m.Set(point{1, 2}, "b")
	if got, want := pairs(m), []string{"{1 2}=b", "{0 0}=origin"}; !slices.Equal(got, want) {
		t.Errorf("pairs = %v, want %v", got, want)
	}
}
