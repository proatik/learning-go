package cache_test

import (
	"testing"

	"learngo/course/04-generics-packages-testing/04-packages/exercises/03-no-stutter/cache"
)

func TestStore(t *testing.T) {
	s := cache.New()
	if s == nil {
		t.Fatalf("New() = nil, want a store")
	}
	s.Set("a", "1")
	s.Set("b", "2")
	s.Set("a", "one") // overwrite
	tests := []struct {
		key    string
		want   string
		wantOK bool
	}{
		{"a", "one", true},
		{"b", "2", true},
		{"zzz", "", false},
		{"", "", false},
	}
	for _, tt := range tests {
		if got, ok := s.Get(tt.key); got != tt.want || ok != tt.wantOK {
			t.Errorf("Get(%q) = %q, %v, want %q, %v", tt.key, got, ok, tt.want, tt.wantOK)
		}
	}
	if got := s.Len(); got != 2 {
		t.Errorf("Len() = %d, want 2", got)
	}
	if got, want := s.Stats(), (cache.Stats{Hits: 2, Misses: 2}); got != want {
		t.Errorf("Stats() = %+v, want %+v", got, want)
	}
	s.Delete("a")
	s.Delete("missing")
	if got := s.Len(); got != 1 {
		t.Errorf("Len() after Delete = %d, want 1", got)
	}
	if _, ok := s.Get("a"); ok {
		t.Errorf("Get(%q) after Delete ok = %v, want %v", "a", ok, false)
	}
}

func TestZeroValue(t *testing.T) {
	var s cache.Store
	if got := s.Len(); got != 0 {
		t.Errorf("zero Store Len() = %d, want 0", got)
	}
	if _, ok := s.Get("x"); ok {
		t.Errorf("zero Store Get(%q) ok = %v, want %v", "x", ok, false)
	}
	s.Delete("x")
	s.Set("k", "v") // must not panic
	if got, ok := s.Get("k"); got != "v" || !ok {
		t.Errorf("Get(%q) = %q, %v, want %q, %v", "k", got, ok, "v", true)
	}
	if got := s.Len(); got != 1 {
		t.Errorf("Len() = %d, want 1", got)
	}
	if got, want := s.Stats(), (cache.Stats{Hits: 1, Misses: 1}); got != want {
		t.Errorf("Stats() = %+v, want %+v", got, want)
	}
}

func TestStoresAreIndependent(t *testing.T) {
	a, b := cache.New(), cache.New()
	a.Set("k", "v")
	if got := b.Len(); got != 0 {
		t.Errorf("second store Len() = %d, want 0", got)
	}
}
