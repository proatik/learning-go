package search_test

import (
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
	"unicode/utf8"

	search "learngo/course/05-concurrency/01-goroutines/exercises/08-first-match"
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

func isEven(n int) bool { return n%2 == 0 }

func TestInts(t *testing.T) {
	tests := []struct {
		items []int
		want  int
	}{
		{[]int{1, 4, 6, 7}, 1},
		{[]int{1, 3, 5}, -1},
		{[]int{2}, 0},
		{[]int{3}, -1},
		{[]int{1, 1, 1, 1, 1, 1, 1, 8}, 7},
		{[]int{-3, -2, 0}, 1},
		{[]int{8, 8, 8}, 0},
		{[]int{}, -1},
		{nil, -1},
	}
	for _, tt := range tests {
		var got int
		within(t, "FirstMatch", func() { got = search.FirstMatch(tt.items, isEven) })
		if got != tt.want {
			t.Errorf("FirstMatch(%v, isEven) = %d, want %d", tt.items, got, tt.want)
		}
	}
}

func TestStrings(t *testing.T) {
	hasLen4 := func(s string) bool { return utf8.RuneCountInString(s) == 4 }
	tests := []struct {
		items []string
		want  int
	}{
		{[]string{"go", "rust", "zig"}, 1},
		{[]string{"héllo", "wörld"}, -1},
		{[]string{"", "Gopher", "café", "java"}, 2},
		{[]string{"ruby"}, 0},
	}
	for _, tt := range tests {
		var got int
		within(t, "FirstMatch", func() { got = search.FirstMatch(tt.items, hasLen4) })
		if got != tt.want {
			t.Errorf("FirstMatch(%q, hasLen4) = %d, want %d", tt.items, got, tt.want)
		}
	}
}

func TestConcurrentAndOrdered(t *testing.T) {
	// Item i's predicate takes durs[i]; later matches finish first, so the
	// answer must come from input order, not finishing order.
	tests := []struct {
		items []int
		durs  []time.Duration
		want  int
		took  time.Duration
	}{
		{[]int{1, 2, 4, 6}, []time.Duration{1, 4, 3, 2}, 1, 4},
		{[]int{2, 4}, []time.Duration{5, 1}, 0, 5},
		{[]int{1, 3, 5, 7, 9}, []time.Duration{1, 1, 1, 1, 1}, -1, 1},
		{[]int{1, 3, 5, 7, 10}, []time.Duration{2, 2, 2, 2, 2}, 4, 2},
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			calls := make([]atomic.Int32, len(tt.items))
			pred := func(n int) bool {
				for i, v := range tt.items {
					if v == n {
						calls[i].Add(1)
						time.Sleep(tt.durs[i] * time.Second)
					}
				}
				return isEven(n)
			}
			start := time.Now()
			got := search.FirstMatch(tt.items, pred)
			elapsed := time.Since(start)
			if got != tt.want {
				t.Errorf("FirstMatch(%v, slow isEven) = %d, want %d", tt.items, got, tt.want)
			}
			for i := range calls {
				if c := calls[i].Load(); c != 1 {
					t.Errorf("FirstMatch(%v): pred(%d) called %d times, want 1", tt.items, tt.items[i], c)
				}
			}
			if want := tt.took * time.Second; elapsed != want {
				t.Errorf("FirstMatch(%v, slow isEven) took %v, want %v", tt.items, elapsed, want)
			}
		})
	}
}

func TestEmptyNeverCalls(t *testing.T) {
	var calls atomic.Int32
	pred := func(int) bool { calls.Add(1); return true }
	var got int
	within(t, "FirstMatch(nil)", func() { got = search.FirstMatch(nil, pred) })
	if got != -1 || calls.Load() != 0 {
		t.Errorf("FirstMatch(nil, pred) = %d with %d calls, want -1 with 0 calls", got, calls.Load())
	}
}

func TestMany(t *testing.T) {
	items := make([]int, 5000)
	for i := range items {
		items[i] = 2*i + 1
	}
	items[3210] = 42
	items[4999] = 44
	var got int
	within(t, "FirstMatch(5000 ints)", func() { got = search.FirstMatch(items, isEven) })
	if got != 3210 {
		t.Errorf("FirstMatch(5000 odd ints with evens at 3210 and 4999, isEven) = %d, want 3210", got)
	}
}
