package wordcount

import (
	"maps"
	"strings"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

// within runs f in a goroutine and fails the test if it doesn't return in 2s.
func within(t *testing.T, name string, f func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		f()
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("%s did not return within 2s", name)
	}
}

func TestCount(t *testing.T) {
	tests := []struct {
		docs []string
		want map[string]int
	}{
		{[]string{"the cat", "The dog the end"}, map[string]int{"the": 3, "cat": 1, "dog": 1, "end": 1}},
		{[]string{"Go go GO"}, map[string]int{"go": 3}},
		{[]string{"  Ünïcode\tünïcode\n"}, map[string]int{"ünïcode": 2}},
		{[]string{"a", "b", "a", "c", "a"}, map[string]int{"a": 3, "b": 1, "c": 1}},
		{[]string{"hello, world", "world!"}, map[string]int{"hello,": 1, "world": 1, "world!": 1}},
		{[]string{"", "   "}, map[string]int{}},
		{[]string{}, map[string]int{}},
		{nil, map[string]int{}},
	}
	for _, tt := range tests {
		var got map[string]int
		within(t, "Count", func() { got = Count(tt.docs) })
		if got == nil {
			t.Errorf("Count(%q) = nil, want %v", tt.docs, tt.want)
			continue
		}
		if !maps.Equal(got, tt.want) {
			t.Errorf("Count(%q) = %v, want %v", tt.docs, got, tt.want)
		}
	}
}

func TestDocumentsCountedConcurrently(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var calls atomic.Int32
		orig := words
		t.Cleanup(func() { words = orig })
		words = func(doc string) []string {
			calls.Add(1)
			time.Sleep(time.Second)
			return orig(doc)
		}
		docs := []string{"a b", "b c", "c d", "d e", "e a"}
		start := time.Now()
		got := Count(docs)
		elapsed := time.Since(start)
		want := map[string]int{"a": 2, "b": 2, "c": 2, "d": 2, "e": 2}
		if !maps.Equal(got, want) {
			t.Errorf("Count(%q) = %v, want %v", docs, got, want)
		}
		if c := calls.Load(); c != int32(len(docs)) {
			t.Errorf("Count(%d docs): words called %d times, want %d", len(docs), c, len(docs))
		}
		if elapsed != time.Second {
			t.Errorf("Count(%d docs, 1s each) took %v, want 1s", len(docs), elapsed)
		}
	})
}

func TestManyDocuments(t *testing.T) {
	docs := make([]string, 300)
	for i := range docs {
		docs[i] = strings.Repeat("go Go rust ", 50) + "zig"
	}
	var got map[string]int
	within(t, "Count(300 docs)", func() { got = Count(docs) })
	want := map[string]int{"go": 300 * 100, "rust": 300 * 50, "zig": 300}
	if !maps.Equal(got, want) {
		t.Errorf("Count(300 docs of %q) = %v, want %v", docs[0][:20]+"…", got, want)
	}
}
