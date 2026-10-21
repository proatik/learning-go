package kv_test

import (
	"slices"
	"testing"

	kv "learngo/course/03-methods-interfaces-errors/02-interfaces/exercises/08-kv-store"
)

// mapGetter is a test fake: any type with Get is a Getter.
type mapGetter map[string]string

func (m mapGetter) Get(k string) (string, bool) { v, ok := m[k]; return v, ok }

// recorder is a test fake Setter that remembers every call.
type recorder struct{ calls []string }

func (r *recorder) Set(k, v string) { r.calls = append(r.calls, k+"="+v) }

func TestMemStoreIsStore(t *testing.T) {
	if _, ok := any(kv.NewMemStore()).(kv.Store); !ok {
		t.Fatalf("*MemStore does not implement Store")
	}
}

func checkGet(t *testing.T, g kv.Getter, key, want string, wantOK bool) {
	t.Helper()
	if got, ok := g.Get(key); got != want || ok != wantOK {
		t.Errorf("Get(%q) = %q, %v, want %q, %v", key, got, ok, want, wantOK)
	}
}

func TestMemStore(t *testing.T) {
	m := kv.NewMemStore()
	if m == nil {
		t.Fatalf("NewMemStore() = nil, want a store")
	}
	checkGet(t, m, "a", "", false)
	m.Set("a", "1")
	checkGet(t, m, "a", "1", true)
	m.Set("a", "one")
	checkGet(t, m, "a", "one", true)
	m.Set("empty", "")
	checkGet(t, m, "empty", "", true)
	m.Set("ключ", "значение")
	checkGet(t, m, "ключ", "значение", true)
	if got := m.Len(); got != 3 {
		t.Errorf("Len() = %d, want 3", got)
	}
	if got, want := m.Keys(), []string{"a", "empty", "ключ"}; !slices.Equal(got, want) {
		t.Errorf("Keys() = %q, want %q", got, want)
	}
	if got := m.Delete("a"); !got {
		t.Errorf("Delete(%q) = %v, want true", "a", got)
	}
	if got := m.Delete("a"); got {
		t.Errorf("Delete(%q) = %v, want false (already deleted)", "a", got)
	}
	if got := m.Delete("never"); got {
		t.Errorf("Delete(%q) = %v, want false", "never", got)
	}
	checkGet(t, m, "a", "", false)
	if got := m.Len(); got != 2 {
		t.Errorf("Len() = %d, want 2", got)
	}
}

func TestZeroMemStore(t *testing.T) {
	var m kv.MemStore
	checkGet(t, &m, "x", "", false)
	if got := m.Len(); got != 0 {
		t.Errorf("zero store: Len() = %d, want 0", got)
	}
	if got := m.Keys(); len(got) != 0 {
		t.Errorf("zero store: Keys() = %q, want []", got)
	}
	if got := m.Delete("x"); got {
		t.Errorf("zero store: Delete(%q) = %v, want false", "x", got)
	}
	m.Set("x", "1")
	checkGet(t, &m, "x", "1", true)
	if got := m.Len(); got != 1 {
		t.Errorf("Len() = %d, want 1", got)
	}
}

func TestGetOr(t *testing.T) {
	g := mapGetter{"a": "1", "blank": ""}
	tests := []struct {
		key, def, want string
	}{
		{"a", "n/a", "1"},
		{"missing", "n/a", "n/a"},
		{"blank", "n/a", ""},
		{"missing", "", ""},
	}
	for _, tt := range tests {
		if got := kv.GetOr(g, tt.key, tt.def); got != tt.want {
			t.Errorf("GetOr(%v, %q, %q) = %q, want %q", g, tt.key, tt.def, got, tt.want)
		}
	}
}

func TestCopyKeys(t *testing.T) {
	src := mapGetter{"a": "1", "b": "2", "blank": ""}
	tests := []struct {
		keys      []string
		want      int
		wantCalls []string
	}{
		{nil, 0, nil},
		{[]string{"a", "x"}, 1, []string{"a=1"}},
		{[]string{"b", "a", "blank"}, 3, []string{"b=2", "a=1", "blank="}},
		{[]string{"x", "y"}, 0, nil},
	}
	for _, tt := range tests {
		dst := &recorder{}
		if got := kv.CopyKeys(dst, src, tt.keys); got != tt.want {
			t.Errorf("CopyKeys(dst, %v, %q) = %d, want %d", src, tt.keys, got, tt.want)
		}
		if !slices.Equal(dst.calls, tt.wantCalls) {
			t.Errorf("CopyKeys(dst, %v, %q): Set calls = %q, want %q", src, tt.keys, dst.calls, tt.wantCalls)
		}
	}
}

func TestCopyKeysBetweenStores(t *testing.T) {
	src, dst := kv.NewMemStore(), kv.NewMemStore()
	src.Set("k", "v")
	if got := kv.CopyKeys(dst, src, []string{"k"}); got != 1 {
		t.Errorf("CopyKeys(dst, src, [k]) = %d, want 1", got)
	}
	checkGet(t, dst, "k", "v", true)
}

func TestReadOnly(t *testing.T) {
	m := kv.NewMemStore()
	m.Set("a", "1")
	ro := kv.ReadOnly(m)
	if ro == nil {
		t.Fatalf("ReadOnly(m) = nil, want a Getter")
	}
	checkGet(t, ro, "a", "1", true)
	checkGet(t, ro, "b", "", false)
	m.Set("b", "2")
	checkGet(t, ro, "b", "2", true) // reads through to the live store
	if _, ok := ro.(kv.Setter); ok {
		t.Errorf("ReadOnly(m).(Setter) ok = true, want false")
	}
	if _, ok := ro.(kv.Store); ok {
		t.Errorf("ReadOnly(m).(Store) ok = true, want false")
	}
	if _, ok := ro.(*kv.MemStore); ok {
		t.Errorf("ReadOnly(m).(*MemStore) ok = true, want false (don't return m itself)")
	}
}
