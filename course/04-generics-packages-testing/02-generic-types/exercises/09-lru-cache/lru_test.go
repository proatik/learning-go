package lru

import (
	"fmt"
	"slices"
	"testing"
)

// walk reads the list head→tail and tail→head (stopping after 1000 nodes in
// case of a cycle) and reports any broken link.
func walk[K comparable, V any](c *Cache[K, V]) (keys []K, problem string) {
	var back []K
	for n := c.head; n != nil && len(keys) < 1000; n = n.next {
		keys = append(keys, n.Key)
		if n.next != nil && n.next.prev != n {
			return keys, fmt.Sprintf("node %v: next.prev doesn't point back", n.Key)
		}
	}
	for n := c.tail; n != nil && len(back) < 1000; n = n.prev {
		back = append(back, n.Key)
	}
	slices.Reverse(back)
	switch {
	case c.head != nil && c.head.prev != nil:
		return keys, "head.prev is not nil"
	case c.tail != nil && c.tail.next != nil:
		return keys, "tail.next is not nil"
	case !slices.Equal(keys, back):
		return keys, fmt.Sprintf("walking from head gives %v, from tail gives %v", keys, back)
	case len(keys) != len(c.items):
		return keys, fmt.Sprintf("list has %d nodes, map has %d", len(keys), len(c.items))
	}
	for _, k := range keys {
		if _, ok := c.items[k]; !ok {
			return keys, fmt.Sprintf("key %v is in the list but not the map", k)
		}
	}
	return keys, ""
}

// check verifies the internal list and Entries against want (most recent
// first).
func check(t *testing.T, step string, c *Cache[string, int], want []Entry[string, int]) {
	t.Helper()
	keys, problem := walk(c)
	if problem != "" {
		t.Fatalf("%s: broken list: %s", step, problem)
	}
	var wantKeys []string
	for _, e := range want {
		wantKeys = append(wantKeys, e.Key)
	}
	if !slices.Equal(keys, wantKeys) {
		t.Fatalf("%s: list order = %v, want %v", step, keys, wantKeys)
	}
	if got := c.Entries(); !slices.Equal(got, want) {
		t.Fatalf("%s: Entries() = %v, want %v", step, got, want)
	}
	if got := c.Len(); got != len(want) {
		t.Fatalf("%s: Len() = %d, want %d", step, got, len(want))
	}
}

type e = Entry[string, int]

func TestNewPanics(t *testing.T) {
	for _, capacity := range []int{0, -1} {
		func() {
			defer func() {
				r := recover()
				if r != "lru: capacity must be at least 1" {
					t.Errorf("New(%d) panic = %v, want %q", capacity, r, "lru: capacity must be at least 1")
				}
			}()
			New[string, int](capacity)
		}()
	}
}

func TestReadmeScenario(t *testing.T) {
	c := New[string, int](2)
	check(t, "New(2)", c, nil)
	c.Put("a", 1)
	c.Put("b", 2)
	check(t, `Put("a", 1), Put("b", 2)`, c, []e{{"b", 2}, {"a", 1}})
	if v, ok := c.Get("a"); v != 1 || !ok {
		t.Fatalf(`Get("a") = (%d, %v), want (1, true)`, v, ok)
	}
	check(t, `Get("a")`, c, []e{{"a", 1}, {"b", 2}})
	if ev, ok := c.Put("c", 3); ev != (e{"b", 2}) || !ok {
		t.Fatalf(`Put("c", 3) = (%v, %v), want ({b 2}, true)`, ev, ok)
	}
	check(t, `Put("c", 3)`, c, []e{{"c", 3}, {"a", 1}})
	if v, ok := c.Peek("a"); v != 1 || !ok {
		t.Fatalf(`Peek("a") = (%d, %v), want (1, true)`, v, ok)
	}
	check(t, `Peek("a")`, c, []e{{"c", 3}, {"a", 1}})
	if ev, ok := c.Put("a", 10); ev != (e{}) || ok {
		t.Fatalf(`Put("a", 10) = (%v, %v), want ({ 0}, false)`, ev, ok)
	}
	check(t, `Put("a", 10)`, c, []e{{"a", 10}, {"c", 3}})
	if v, ok := c.Get("zz"); v != 0 || ok {
		t.Fatalf(`Get("zz") = (%d, %v), want (0, false)`, v, ok)
	}
	if v, ok := c.Peek("b"); v != 0 || ok {
		t.Fatalf(`Peek("b") after eviction = (%d, %v), want (0, false)`, v, ok)
	}
	if !c.Remove("a") {
		t.Fatalf(`Remove("a") = false, want true`)
	}
	if c.Remove("a") {
		t.Fatalf(`second Remove("a") = true, want false`)
	}
	check(t, `Remove("a")`, c, []e{{"c", 3}})
}

func TestRemovePositions(t *testing.T) {
	tests := []struct {
		remove string
		want   []e
	}{
		{"x", []e{{"z", 3}, {"y", 2}}}, // tail (least recent)
		{"y", []e{{"z", 3}, {"x", 1}}}, // middle
		{"z", []e{{"y", 2}, {"x", 1}}}, // head (most recent)
	}
	for _, tt := range tests {
		c := New[string, int](3)
		c.Put("x", 1)
		c.Put("y", 2)
		c.Put("z", 3)
		check(t, "Put x y z", c, []e{{"z", 3}, {"y", 2}, {"x", 1}})
		if !c.Remove(tt.remove) {
			t.Fatalf("Remove(%q) = false, want true", tt.remove)
		}
		check(t, fmt.Sprintf("Remove(%q)", tt.remove), c, tt.want)
		c.Put("w", 4)
		check(t, fmt.Sprintf("Remove(%q), Put(\"w\", 4)", tt.remove), c, append([]e{{"w", 4}}, tt.want...))
	}
}

func TestRemoveOnlyEntry(t *testing.T) {
	c := New[string, int](1)
	c.Put("solo", 1)
	c.Remove("solo")
	check(t, `Put("solo"), Remove("solo")`, c, nil)
	if c.head != nil || c.tail != nil {
		t.Fatalf("after removing the only entry, head/tail should be nil")
	}
	c.Put("again", 2)
	check(t, `Put("again", 2)`, c, []e{{"again", 2}})
}

func TestCapacityOne(t *testing.T) {
	c := New[string, int](1)
	c.Put("a", 1)
	if ev, ok := c.Put("b", 2); ev != (e{"a", 1}) || !ok {
		t.Fatalf(`capacity 1: Put("b", 2) = (%v, %v), want ({a 1}, true)`, ev, ok)
	}
	if ev, ok := c.Put("b", 3); ok {
		t.Fatalf(`capacity 1: Put("b", 3) = (%v, %v), want ({ 0}, false)`, ev, ok)
	}
	check(t, "capacity 1", c, []e{{"b", 3}})
}

// model is a slow, obviously-correct LRU used to cross-check the cache.
type model struct {
	capacity int
	entries  []e // most recent first
}

func (m *model) find(k string) int {
	return slices.IndexFunc(m.entries, func(x e) bool { return x.Key == k })
}

func (m *model) touch(i int) {
	x := m.entries[i]
	m.entries = slices.Insert(slices.Delete(m.entries, i, i+1), 0, x)
}

func TestAgainstModel(t *testing.T) {
	c := New[string, int](4)
	m := &model{capacity: 4}
	seed := uint32(7)
	next := func(n int) int {
		seed = seed*1664525 + 1013904223
		return int(seed>>16) % n
	}
	for step := range 2000 {
		k := string(rune('a' + next(7)))
		switch op := next(4); op {
		case 0, 1:
			v := next(100)
			ev, ok := c.Put(k, v)
			var wantEv e
			wantOK := false
			if i := m.find(k); i >= 0 {
				m.entries[i].Val = v
				m.touch(i)
			} else {
				if len(m.entries) == m.capacity {
					wantEv, wantOK = m.entries[len(m.entries)-1], true
					m.entries = m.entries[:len(m.entries)-1]
				}
				m.entries = slices.Insert(m.entries, 0, e{k, v})
			}
			if ev != wantEv || ok != wantOK {
				t.Fatalf("step %d: Put(%q, %d) = (%v, %v), want (%v, %v)", step, k, v, ev, ok, wantEv, wantOK)
			}
		case 2:
			v, ok := c.Get(k)
			want, wantOK := 0, false
			if i := m.find(k); i >= 0 {
				want, wantOK = m.entries[i].Val, true
				m.touch(i)
			}
			if v != want || ok != wantOK {
				t.Fatalf("step %d: Get(%q) = (%d, %v), want (%d, %v)", step, k, v, ok, want, wantOK)
			}
		case 3:
			got := c.Remove(k)
			i := m.find(k)
			if i >= 0 {
				m.entries = slices.Delete(m.entries, i, i+1)
			}
			if got != (i >= 0) {
				t.Fatalf("step %d: Remove(%q) = %v, want %v", step, k, got, i >= 0)
			}
		}
		check(t, fmt.Sprintf("step %d", step), c, m.entries)
	}
}

func TestOtherTypes(t *testing.T) {
	type point struct{ X, Y int }
	c := New[point, *string](2)
	s := "origin"
	c.Put(point{0, 0}, &s)
	c.Put(point{1, 2}, nil)
	if v, ok := c.Get(point{0, 0}); v != &s || !ok {
		t.Errorf("Get({0 0}) = (%v, %v), want (%p, true)", v, ok, &s)
	}
	if v, ok := c.Get(point{1, 2}); v != nil || !ok {
		t.Errorf("Get({1 2}) = (%v, %v), want (<nil>, true)", v, ok)
	}
	if ev, ok := c.Put(point{3, 3}, nil); ev.Key != (point{0, 0}) || ev.Val != &s || !ok {
		t.Errorf("Put({3 3}) evicted (%v, %v), want ({{0 0} %p}, true)", ev, ok, &s)
	}
}
