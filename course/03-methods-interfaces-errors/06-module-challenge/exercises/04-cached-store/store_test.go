package kv_test

import (
	"errors"
	"fmt"
	"testing"

	kv "learngo/course/03-methods-interfaces-errors/06-module-challenge/exercises/04-cached-store"
)

// The concrete types must satisfy the interface.
var (
	_ kv.Store = (*kv.MemoryStore)(nil)
	_ kv.Store = (*kv.CachedStore)(nil)
)

var errDown = errors.New("backend down")

// spy is a Store written by the test: CachedStore must work with ANY Store.
type spy struct {
	data       map[string]string
	gets, sets int
	deletes    int
	failSet    bool
	failGet    bool
}

func newSpy() *spy { return &spy{data: map[string]string{}} }

func (s *spy) Get(key string) (string, error) {
	s.gets++
	if s.failGet {
		return "", errDown
	}
	v, ok := s.data[key]
	if !ok {
		return "", kv.ErrNotFound
	}
	return v, nil
}

func (s *spy) Set(key, value string) error {
	s.sets++
	if s.failSet {
		return errDown
	}
	s.data[key] = value
	return nil
}

func (s *spy) Delete(key string) error {
	s.deletes++
	if _, ok := s.data[key]; !ok {
		return kv.ErrNotFound
	}
	delete(s.data, key)
	return nil
}

func newMemory(t *testing.T) *kv.MemoryStore {
	t.Helper()
	m := kv.NewMemoryStore()
	if m == nil {
		t.Fatalf("NewMemoryStore() = nil, want a store")
	}
	return m
}

func newCached(t *testing.T, backend kv.Store) *kv.CachedStore {
	t.Helper()
	c := kv.NewCachedStore(backend)
	if c == nil {
		t.Fatalf("NewCachedStore(...) = nil, want a store")
	}
	return c
}

func wantGet(t *testing.T, s kv.Store, key, want string) {
	t.Helper()
	got, err := s.Get(key)
	if err != nil || got != want {
		t.Fatalf("Get(%q) = %q, %v, want %q, nil", key, got, err, want)
	}
}

func TestMemoryStore(t *testing.T) {
	m := newMemory(t)
	if got := m.Len(); got != 0 {
		t.Fatalf("Len() = %d, want 0", got)
	}
	if err := m.Set("lang", "go"); err != nil {
		t.Fatalf("Set(%q, %q) error = %v, want nil", "lang", "go", err)
	}
	if err := m.Set("empty", ""); err != nil { // empty values are fine
		t.Fatalf("Set(%q, %q) error = %v, want nil", "empty", "", err)
	}
	wantGet(t, m, "lang", "go")
	wantGet(t, m, "empty", "")
	if err := m.Set("lang", "Go 1.27"); err != nil {
		t.Fatalf("Set(%q) overwrite error = %v, want nil", "lang", err)
	}
	wantGet(t, m, "lang", "Go 1.27")
	if got := m.Len(); got != 2 {
		t.Errorf("Len() = %d, want 2", got)
	}

	_, err := m.Get("nope")
	if !errors.Is(err, kv.ErrNotFound) {
		t.Errorf("Get(%q) error = %v, want %v", "nope", err, kv.ErrNotFound)
	} else if want := `get "nope": kv: key not found`; err.Error() != want {
		t.Errorf("Get(%q) error = %q, want %q", "nope", err.Error(), want)
	}

	if err := m.Delete("lang"); err != nil {
		t.Fatalf("Delete(%q) error = %v, want nil", "lang", err)
	}
	err = m.Delete("lang")
	if !errors.Is(err, kv.ErrNotFound) {
		t.Errorf("Delete(%q) twice error = %v, want %v", "lang", err, kv.ErrNotFound)
	} else if want := `delete "lang": kv: key not found`; err.Error() != want {
		t.Errorf("Delete(%q) error = %q, want %q", "lang", err.Error(), want)
	}
	if got := m.Len(); got != 1 {
		t.Errorf("Len() = %d, want 1", got)
	}
}

func TestEmptyKey(t *testing.T) {
	m := newMemory(t)
	s := newSpy()
	c := newCached(t, s)
	for _, st := range []kv.Store{m, c} {
		if _, err := st.Get(""); !errors.Is(err, kv.ErrEmptyKey) {
			t.Errorf("%T.Get(\"\") error = %v, want %v", st, err, kv.ErrEmptyKey)
		}
		if err := st.Set("", "x"); !errors.Is(err, kv.ErrEmptyKey) {
			t.Errorf("%T.Set(\"\", \"x\") error = %v, want %v", st, err, kv.ErrEmptyKey)
		}
		if err := st.Delete(""); !errors.Is(err, kv.ErrEmptyKey) {
			t.Errorf("%T.Delete(\"\") error = %v, want %v", st, err, kv.ErrEmptyKey)
		}
	}
	if s.gets+s.sets+s.deletes != 0 {
		t.Errorf("CachedStore called the backend for an empty key (%d calls), want 0", s.gets+s.sets+s.deletes)
	}
	if h, ms := c.Stats(); h != 0 || ms != 0 {
		t.Errorf("Stats() after empty keys = %d, %d, want 0, 0", h, ms)
	}
}

func TestCachedReadThrough(t *testing.T) {
	s := newSpy()
	s.data["a"] = "1"
	c := newCached(t, s)

	wantGet(t, c, "a", "1") // miss: loads from the backend
	wantGet(t, c, "a", "1") // hit
	wantGet(t, c, "a", "1") // hit
	if s.gets != 1 {
		t.Errorf("backend Get calls = %d, want 1 (later reads come from the cache)", s.gets)
	}
	if h, m := c.Stats(); h != 2 || m != 1 {
		t.Errorf("Stats() = %d, %d, want 2, 1", h, m)
	}

	_, err := c.Get("missing")
	if !errors.Is(err, kv.ErrNotFound) {
		t.Fatalf("Get(%q) error = %v, want %v", "missing", err, kv.ErrNotFound)
	}
	if _, err := c.Get("missing"); !errors.Is(err, kv.ErrNotFound) {
		t.Fatalf("Get(%q) again error = %v, want %v", "missing", err, kv.ErrNotFound)
	}
	if s.gets != 3 {
		t.Errorf("backend Get calls = %d, want 3 (errors are never cached)", s.gets)
	}
	if h, m := c.Stats(); h != 2 || m != 3 {
		t.Errorf("Stats() = %d, %d, want 2, 3", h, m)
	}
	if got := c.String(); got != "cache: 2 hits, 3 misses, 1 cached" {
		t.Errorf("String() = %q, want %q", got, "cache: 2 hits, 3 misses, 1 cached")
	}
	if got := fmt.Sprint(c); got != "cache: 2 hits, 3 misses, 1 cached" {
		t.Errorf("fmt.Sprint(cache) = %q, want %q", got, "cache: 2 hits, 3 misses, 1 cached")
	}
}

func TestCachedWrapsErrors(t *testing.T) {
	m := newMemory(t)
	c := newCached(t, m)
	_, err := c.Get("k")
	if want := `cache: get "k": kv: key not found`; err == nil || err.Error() != want {
		t.Errorf("Get(%q) error = %v, want %q", "k", err, want)
	}
	if !errors.Is(err, kv.ErrNotFound) {
		t.Errorf("errors.Is(%v, ErrNotFound) = false, want true", err)
	}

	s := newSpy()
	s.failGet = true
	c2 := newCached(t, s)
	if _, err := c2.Get("x"); !errors.Is(err, errDown) || err.Error() != "cache: backend down" {
		t.Errorf("Get(%q) on a failing backend error = %v, want %q wrapping errDown", "x", err, "cache: backend down")
	}
}

func TestCachedSet(t *testing.T) {
	s := newSpy()
	c := newCached(t, s)
	if err := c.Set("a", "1"); err != nil {
		t.Fatalf("Set(%q, %q) error = %v, want nil", "a", "1", err)
	}
	if s.data["a"] != "1" {
		t.Fatalf("backend[%q] = %q after Set, want %q (write-through)", "a", s.data["a"], "1")
	}
	wantGet(t, c, "a", "1")
	if s.gets != 0 {
		t.Errorf("backend Get calls = %d, want 0 (Set already cached the value)", s.gets)
	}

	s.failSet = true
	err := c.Set("a", "2")
	if !errors.Is(err, errDown) {
		t.Fatalf("Set on a failing backend error = %v, want %v", err, errDown)
	}
	wantGet(t, c, "a", "1") // the failed write must not reach the cache
	if s.data["a"] != "1" {
		t.Errorf("backend[%q] = %q, want %q", "a", s.data["a"], "1")
	}
}

func TestCachedDelete(t *testing.T) {
	s := newSpy()
	c := newCached(t, s)
	if err := c.Set("a", "1"); err != nil {
		t.Fatalf("Set error = %v", err)
	}
	wantGet(t, c, "a", "1")
	if err := c.Delete("a"); err != nil {
		t.Fatalf("Delete(%q) error = %v, want nil", "a", err)
	}
	if _, ok := s.data["a"]; ok {
		t.Errorf("Delete(%q) did not delete from the backend", "a")
	}
	if _, err := c.Get("a"); !errors.Is(err, kv.ErrNotFound) {
		t.Errorf("Get(%q) after Delete error = %v, want %v (must not serve a stale value)", "a", err, kv.ErrNotFound)
	}
	if err := c.Delete("a"); !errors.Is(err, kv.ErrNotFound) {
		t.Errorf("Delete(%q) twice error = %v, want %v", "a", err, kv.ErrNotFound)
	}

	// A key changed behind the cache's back is dropped by Delete even if the
	// backend no longer has it.
	if err := c.Set("b", "2"); err != nil {
		t.Fatalf("Set error = %v", err)
	}
	delete(s.data, "b")
	if err := c.Delete("b"); !errors.Is(err, kv.ErrNotFound) {
		t.Errorf("Delete(%q) error = %v, want %v", "b", err, kv.ErrNotFound)
	}
	if _, err := c.Get("b"); !errors.Is(err, kv.ErrNotFound) {
		t.Errorf("Get(%q) after Delete error = %v, want %v", "b", err, kv.ErrNotFound)
	}
}

func TestStackedCaches(t *testing.T) {
	m := newMemory(t)
	inner := newCached(t, m)
	outer := newCached(t, inner) // a CachedStore is a Store, so caches stack
	if err := m.Set("k", "v"); err != nil {
		t.Fatalf("Set error = %v", err)
	}
	wantGet(t, outer, "k", "v")
	wantGet(t, outer, "k", "v")
	if got := inner.String(); got != "cache: 0 hits, 1 misses, 1 cached" {
		t.Errorf("inner.String() = %q, want %q", got, "cache: 0 hits, 1 misses, 1 cached")
	}
	if got := outer.String(); got != "cache: 1 hits, 1 misses, 1 cached" {
		t.Errorf("outer.String() = %q, want %q", got, "cache: 1 hits, 1 misses, 1 cached")
	}
	_, err := outer.Get("zz")
	if want := `cache: cache: get "zz": kv: key not found`; err == nil || err.Error() != want {
		t.Errorf("Get(%q) error = %v, want %q", "zz", err, want)
	}
}
