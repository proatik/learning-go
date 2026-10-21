package cache

import (
	"fmt"
	"slices"
	"strconv"
	"testing"
)

func newLRU(t *testing.T, capacity int) *LRU {
	t.Helper()
	c := NewLRU(capacity)
	if c == nil {
		t.Fatalf("NewLRU(%d) = nil, want a cache", capacity)
	}
	if c.items == nil {
		t.Fatalf("NewLRU(%d).items = nil, want an empty map", capacity)
	}
	return c
}

// check verifies the list/map invariants and that the keys, most recent
// first, equal want.
func check(t *testing.T, c *LRU, want ...string) {
	t.Helper()
	if want == nil {
		want = []string{}
	}
	var forward []string
	var prev *node
	for n := c.head; n != nil; n = n.next {
		if len(forward) > len(c.items) {
			t.Fatalf("list from head is longer than the map (cycle?): %v...", forward)
		}
		if n.prev != prev {
			t.Fatalf("node %q: prev link is inconsistent", n.key)
		}
		if c.items[n.key] != n {
			t.Fatalf("node %q in the list is not the node stored in items", n.key)
		}
		forward = append(forward, n.key)
		prev = n
	}
	if c.tail != prev {
		t.Fatalf("tail does not point to the last node of the list")
	}
	if len(forward) != len(c.items) {
		t.Fatalf("list has %d nodes, items has %d", len(forward), len(c.items))
	}
	if !slices.Equal(forward, want) {
		t.Fatalf("list order = %v, want %v", forward, want)
	}
	if got := c.Keys(); !slices.Equal(got, want) || got == nil {
		t.Fatalf("Keys() = %#v, want %v", got, want)
	}
	if got := c.Len(); got != len(want) {
		t.Fatalf("Len() = %d, want %d", got, len(want))
	}
}

func put(t *testing.T, c *LRU, key string, value int, wantEvicted string, wantOK bool) {
	t.Helper()
	ev, ok := c.Put(key, value)
	if ev != wantEvicted || ok != wantOK {
		t.Fatalf("Put(%q, %d) = %q, %v, want %q, %v", key, value, ev, ok, wantEvicted, wantOK)
	}
}

func get(t *testing.T, c *LRU, key string, want int, wantOK bool) {
	t.Helper()
	v, ok := c.Get(key)
	if v != want || ok != wantOK {
		t.Fatalf("Get(%q) = %d, %v, want %d, %v", key, v, ok, want, wantOK)
	}
}

func TestNewLRU(t *testing.T) {
	for _, n := range []int{0, -1, -100} {
		if c := NewLRU(n); c != nil {
			t.Errorf("NewLRU(%d) = %v, want nil", n, c)
		}
	}
	c := newLRU(t, 3)
	if c.capacity != 3 {
		t.Errorf("capacity = %d, want 3", c.capacity)
	}
	check(t, c)
}

func TestExample(t *testing.T) {
	c := newLRU(t, 2)
	put(t, c, "a", 1, "", false)
	put(t, c, "b", 2, "", false)
	check(t, c, "b", "a")
	get(t, c, "a", 1, true)
	check(t, c, "a", "b")
	put(t, c, "c", 3, "b", true)
	check(t, c, "c", "a")
	get(t, c, "b", 0, false)
	if got := c.String(); got != "[c:3 a:1]" {
		t.Errorf("String() = %q, want %q", got, "[c:3 a:1]")
	}
}

func TestUpdateAndPeek(t *testing.T) {
	c := newLRU(t, 3)
	put(t, c, "a", 1, "", false)
	put(t, c, "b", 2, "", false)
	put(t, c, "c", 3, "", false)
	put(t, c, "a", 10, "", false) // update: no eviction even though full
	check(t, c, "a", "c", "b")
	if v, ok := c.Peek("b"); v != 2 || !ok {
		t.Fatalf("Peek(%q) = %d, %v, want %d, %v", "b", v, ok, 2, true)
	}
	check(t, c, "a", "c", "b") // Peek does not change the order
	if v, ok := c.Peek("zz"); v != 0 || ok {
		t.Fatalf("Peek(%q) = %d, %v, want %d, %v", "zz", v, ok, 0, false)
	}
	put(t, c, "d", 4, "b", true)
	check(t, c, "d", "a", "c")
	put(t, c, "c", 30, "", false) // update the tail
	check(t, c, "c", "d", "a")
	get(t, c, "c", 30, true)
}

func TestCapacityOne(t *testing.T) {
	c := newLRU(t, 1)
	put(t, c, "x", 1, "", false)
	put(t, c, "y", 2, "x", true)
	check(t, c, "y")
	put(t, c, "y", 3, "", false)
	put(t, c, "z", 4, "y", true)
	check(t, c, "z")
}

func TestDelete(t *testing.T) {
	c := newLRU(t, 4)
	for i, k := range []string{"a", "b", "c", "d"} {
		put(t, c, k, i, "", false)
	}
	tests := []struct {
		key  string
		ok   bool
		left []string
	}{
		{"b", true, []string{"d", "c", "a"}},
		{"b", false, []string{"d", "c", "a"}},
		{"d", true, []string{"c", "a"}},
		{"a", true, []string{"c"}},
		{"c", true, []string{}},
		{"c", false, []string{}},
	}
	for _, tt := range tests {
		if got := c.Delete(tt.key); got != tt.ok {
			t.Fatalf("Delete(%q) = %v, want %v", tt.key, got, tt.ok)
		}
		check(t, c, tt.left...)
	}
	put(t, c, "x", 9, "", false)
	check(t, c, "x")
}

// recorder collects evictions; its method value is used as the callback.
type recorder struct{ got []string }

func (r *recorder) record(key string, value int) {
	r.got = append(r.got, fmt.Sprintf("%s=%d", key, value))
}

func TestOnEvict(t *testing.T) {
	c := newLRU(t, 2)
	put(t, c, "a", 1, "", false) // no callback set yet: must not crash
	put(t, c, "b", 2, "", false)
	put(t, c, "c", 3, "a", true)

	rec := &recorder{}
	c.OnEvict = rec.record // method value: bound to rec
	put(t, c, "d", 4, "b", true)
	put(t, c, "d", 40, "", false) // update: no eviction
	c.Delete("c")                 // Delete is not an eviction
	put(t, c, "e", 5, "", false)
	put(t, c, "f", 6, "d", true)
	if want := []string{"b=2", "d=40"}; !slices.Equal(rec.got, want) {
		t.Errorf("OnEvict calls = %v, want %v", rec.got, want)
	}
}

func TestString(t *testing.T) {
	c := newLRU(t, 3)
	if got := c.String(); got != "[]" {
		t.Errorf("String() on empty cache = %q, want %q", got, "[]")
	}
	put(t, c, "x", -1, "", false)
	if got := fmt.Sprint(c); got != "[x:-1]" {
		t.Errorf("fmt.Sprint(cache) = %q, want %q", got, "[x:-1]")
	}
	put(t, c, "héllo", 7, "", false)
	put(t, c, "y", 0, "", false)
	var s fmt.Stringer = c // *LRU implements fmt.Stringer
	if got := s.String(); got != "[y:0 héllo:7 x:-1]" {
		t.Errorf("String() = %q, want %q", got, "[y:0 héllo:7 x:-1]")
	}
}

func TestNilReceiver(t *testing.T) {
	var c *LRU // e.g. the result of NewLRU(0)
	if got := c.Len(); got != 0 {
		t.Errorf("(*LRU)(nil).Len() = %d, want 0", got)
	}
	if got := c.Keys(); got == nil || len(got) != 0 {
		t.Errorf("(*LRU)(nil).Keys() = %#v, want []string{}", got)
	}
	if got := c.String(); got != "[]" {
		t.Errorf("(*LRU)(nil).String() = %q, want %q", got, "[]")
	}
}

func TestManyOperations(t *testing.T) {
	const capacity = 50
	c := newLRU(t, capacity)
	var model []string
	values := map[string]int{}
	for i := range 2000 {
		key := strconv.Itoa(i * 7 % 83)
		if i%3 == 0 {
			v, ok := c.Get(key)
			idx := slices.Index(model, key)
			if ok != (idx >= 0) || (ok && v != values[key]) {
				t.Fatalf("step %d: Get(%q) = %d, %v; model has it: %v", i, key, v, ok, idx >= 0)
			}
			if idx >= 0 {
				model = slices.Insert(slices.Delete(model, idx, idx+1), 0, key)
			}
			continue
		}
		wantEv, wantOK := "", false
		if idx := slices.Index(model, key); idx >= 0 {
			model = slices.Delete(model, idx, idx+1)
		} else if len(model) == capacity {
			wantEv, wantOK = model[len(model)-1], true
			model = model[:len(model)-1]
		}
		model = slices.Insert(model, 0, key)
		values[key] = i
		put(t, c, key, i, wantEv, wantOK)
	}
	check(t, c, model...)
}
