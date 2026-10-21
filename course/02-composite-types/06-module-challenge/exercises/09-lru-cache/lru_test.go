package lru

import (
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
	var backward []string
	for n := c.tail; n != nil; n = n.prev {
		if len(backward) > len(c.items) {
			t.Fatalf("list from tail is longer than the map (cycle?)")
		}
		backward = append(backward, n.key)
	}
	slices.Reverse(backward)
	if !slices.Equal(forward, backward) {
		t.Fatalf("forward list %v and backward list %v differ", forward, backward)
	}
	if len(forward) != len(c.items) {
		t.Fatalf("list has %d nodes, items has %d", len(forward), len(c.items))
	}
	if !slices.Equal(forward, want) {
		t.Fatalf("list order = %v, want %v", forward, want)
	}
	if got := Keys(c); !slices.Equal(got, want) || got == nil {
		t.Fatalf("Keys() = %#v, want %v", got, want)
	}
	if got := Len(c); got != len(want) {
		t.Fatalf("Len() = %d, want %d", got, len(want))
	}
}

func put(t *testing.T, c *LRU, key string, value int, wantEvicted string, wantOK bool) {
	t.Helper()
	ev, ok := Put(c, key, value)
	if ev != wantEvicted || ok != wantOK {
		t.Fatalf("Put(%q, %d) = %q, %v, want %q, %v", key, value, ev, ok, wantEvicted, wantOK)
	}
}

func get(t *testing.T, c *LRU, key string, want int, wantOK bool) {
	t.Helper()
	v, ok := Get(c, key)
	if v != want || ok != wantOK {
		t.Fatalf("Get(%q) = %d, %v, want %d, %v", key, v, ok, want, wantOK)
	}
}

func TestNewLRU(t *testing.T) {
	for _, n := range []int{0, -1, -100} {
		if c := NewLRU(n); c != nil {
			t.Errorf("NewLRU(%d) = %+v, want nil", n, c)
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
	check(t, c, "c", "a")
	get(t, c, "c", 3, true)
	get(t, c, "a", 1, true)
	check(t, c, "a", "c")
}

func TestUpdateExisting(t *testing.T) {
	c := newLRU(t, 3)
	put(t, c, "a", 1, "", false)
	put(t, c, "b", 2, "", false)
	put(t, c, "c", 3, "", false)
	put(t, c, "a", 10, "", false) // update: no eviction even though full
	check(t, c, "a", "c", "b")
	get(t, c, "a", 10, true)
	put(t, c, "d", 4, "b", true)
	check(t, c, "d", "a", "c")
	put(t, c, "c", 30, "", false) // update the tail
	check(t, c, "c", "d", "a")
	put(t, c, "c", 31, "", false) // update the head
	check(t, c, "c", "d", "a")
	get(t, c, "c", 31, true)
}

func TestCapacityOne(t *testing.T) {
	c := newLRU(t, 1)
	put(t, c, "x", 1, "", false)
	check(t, c, "x")
	get(t, c, "x", 1, true)
	put(t, c, "y", 2, "x", true)
	check(t, c, "y")
	get(t, c, "x", 0, false)
	put(t, c, "y", 3, "", false)
	check(t, c, "y")
	put(t, c, "z", 4, "y", true)
	check(t, c, "z")
}

func TestGetMovesMiddleAndTail(t *testing.T) {
	c := newLRU(t, 4)
	for i, k := range []string{"a", "b", "c", "d"} {
		put(t, c, k, i, "", false)
	}
	check(t, c, "d", "c", "b", "a")
	get(t, c, "b", 1, true) // middle
	check(t, c, "b", "d", "c", "a")
	get(t, c, "a", 0, true) // tail
	check(t, c, "a", "b", "d", "c")
	get(t, c, "a", 0, true) // head stays head
	check(t, c, "a", "b", "d", "c")
	get(t, c, "zz", 0, false)
	check(t, c, "a", "b", "d", "c")
	put(t, c, "e", 4, "c", true)
	check(t, c, "e", "a", "b", "d")
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
		{"b", true, []string{"d", "c", "a"}},  // middle
		{"b", false, []string{"d", "c", "a"}}, // already gone
		{"d", true, []string{"c", "a"}},       // head
		{"a", true, []string{"c"}},            // tail
		{"zz", false, []string{"c"}},          // never existed
		{"c", true, []string{}},               // last one
		{"c", false, []string{}},              // empty cache
	}
	for _, tt := range tests {
		if got := Delete(c, tt.key); got != tt.ok {
			t.Fatalf("Delete(%q) = %v, want %v", tt.key, got, tt.ok)
		}
		check(t, c, tt.left...)
	}
	if c.head != nil || c.tail != nil {
		t.Fatalf("empty cache: head/tail not nil")
	}
	// The cache is still usable after emptying it.
	put(t, c, "x", 9, "", false)
	check(t, c, "x")
	get(t, c, "x", 9, true)
}

func TestDeleteFreesCapacity(t *testing.T) {
	c := newLRU(t, 2)
	put(t, c, "a", 1, "", false)
	put(t, c, "b", 2, "", false)
	Delete(c, "a")
	put(t, c, "c", 3, "", false) // room again: no eviction
	check(t, c, "c", "b")
}

func TestManyOperations(t *testing.T) {
	const capacity = 50
	c := newLRU(t, capacity)
	// Reference model: a slice ordered most recent first.
	var model []string
	values := map[string]int{}
	for i := range 2000 {
		key := strconv.Itoa(i * 7 % 83)
		if i%3 == 0 {
			v, ok := Get(c, key)
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
			delete(values, wantEv)
		}
		model = slices.Insert(model, 0, key)
		values[key] = i
		put(t, c, key, i, wantEv, wantOK)
	}
	check(t, c, model...)
}
