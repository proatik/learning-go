package lookup_test

import (
	"errors"
	"strings"
	"testing"

	lookup "learngo/course/04-generics-packages-testing/04-packages/exercises/03-no-stutter"
	"learngo/course/04-generics-packages-testing/04-packages/exercises/03-no-stutter/cache"
)

func TestCached(t *testing.T) {
	calls := 0
	load := func(key string) (string, error) {
		calls++
		return strings.ToUpper(key), nil
	}
	var s cache.Store
	for i := range 3 {
		got, err := lookup.Cached(&s, "go", load)
		if got != "GO" || err != nil {
			t.Errorf("Cached(%q) call %d = %q, %v, want %q, <nil>", "go", i+1, got, err, "GO")
		}
	}
	if calls != 1 {
		t.Errorf("load calls after 3 Cached(%q) = %d, want 1", "go", calls)
	}
	if got, ok := s.Get("go"); got != "GO" || !ok {
		t.Errorf("Get(%q) = %q, %v, want %q, %v", "go", got, ok, "GO", true)
	}
}

func TestCachedError(t *testing.T) {
	errDown := errors.New("backend down")
	calls := 0
	load := func(string) (string, error) {
		calls++
		return "", errDown
	}
	s := cache.New()
	if s == nil {
		t.Fatalf("New() = nil, want a store")
	}
	for range 2 {
		if _, err := lookup.Cached(s, "k", load); !errors.Is(err, errDown) {
			t.Errorf("Cached(%q) error = %v, want %v", "k", err, errDown)
		}
	}
	if calls != 2 {
		t.Errorf("load calls = %d, want 2 (errors must not be cached)", calls)
	}
	if got := s.Len(); got != 0 {
		t.Errorf("Len() after failed loads = %d, want 0", got)
	}
}
