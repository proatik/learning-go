package bounded

import (
	"fmt"
	"runtime"
	"slices"
	"strconv"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

// mapWithin runs Map outside a bubble and fails if it takes over 2s.
func mapWithin[T, R any](t *testing.T, name string, limit int, in []T, f func(T) R) []R {
	t.Helper()
	done := make(chan []R, 1)
	go func() { done <- Map(limit, in, f) }()
	select {
	case got := <-done:
		return got
	case <-time.After(2 * time.Second):
		t.Fatalf("%s did not return within 2s", name)
		return nil
	}
}

// peak tracks the highest value a counter reaches.
type peak struct{ cur, max atomic.Int64 }

func (p *peak) enter() {
	now := p.cur.Add(1)
	for {
		m := p.max.Load()
		if now <= m || p.max.CompareAndSwap(m, now) {
			return
		}
	}
}

func (p *peak) leave() { p.cur.Add(-1) }

func TestMapResults(t *testing.T) {
	in := []int{3, -1, 0, 42, 7}
	got := mapWithin(t, "Map", 2, in, strconv.Itoa)
	if want := []string{"3", "-1", "0", "42", "7"}; !slices.Equal(got, want) {
		t.Errorf("Map(2, %v, strconv.Itoa) = %q, want %q", in, got, want)
	}
	if got := mapWithin(t, "Map", 3, []int{}, strconv.Itoa); len(got) != 0 {
		t.Errorf("Map(3, [], strconv.Itoa) = %q, want []", got)
	}
	if got := mapWithin(t, "Map", 3, nil, strconv.Itoa); len(got) != 0 {
		t.Errorf("Map(3, nil, strconv.Itoa) = %q, want []", got)
	}
	words := []string{"go", "héllo", ""}
	lens := mapWithin(t, "Map", 0, words, func(s string) int { return len(s) })
	if want := []int{2, 6, 0}; !slices.Equal(lens, want) {
		t.Errorf("Map(0, %q, len) = %v, want %v", words, lens, want)
	}
}

// TestMapKeepsOrder: later items finish first, results still follow the
// input order.
func TestMapKeepsOrder(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		in := []int{1, 2, 3, 4, 5, 6}
		got := Map(6, in, func(n int) int {
			time.Sleep(time.Duration(10-n) * time.Second)
			return n * n
		})
		if want := []int{1, 4, 9, 16, 25, 36}; !slices.Equal(got, want) {
			t.Errorf("Map(6, %v, square) = %v, want %v", in, got, want)
		}
	})
}

func TestMapLimit(t *testing.T) {
	tests := []struct {
		limit, items int
		wantPeak     int64
		wantElapsed  time.Duration
	}{
		{3, 10, 3, 4 * time.Second},
		{1, 4, 1, 4 * time.Second},
		{4, 4, 4, 1 * time.Second},
		{50, 10, 10, 1 * time.Second},
		{0, 3, 1, 3 * time.Second},
		{-1, 2, 1, 2 * time.Second},
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			name := fmt.Sprintf("Map(%d, %d items of 1s)", tt.limit, tt.items)
			in := make([]int, tt.items)
			for i := range in {
				in[i] = i
			}
			var p peak
			start := time.Now()
			got := Map(tt.limit, in, func(n int) int {
				p.enter()
				defer p.leave()
				time.Sleep(time.Second)
				return n + 100
			})
			elapsed := time.Since(start)
			if len(got) != tt.items {
				t.Fatalf("%s returned %d results, want %d", name, len(got), tt.items)
			}
			for i, v := range got {
				if v != i+100 {
					t.Errorf("%s: result[%d] = %d, want %d", name, i, v, i+100)
				}
			}
			if m := p.max.Load(); m != tt.wantPeak {
				t.Errorf("%s max concurrent calls = %d, want %d", name, m, tt.wantPeak)
			}
			if elapsed != tt.wantElapsed {
				t.Errorf("%s took %v, want %v", name, elapsed, tt.wantElapsed)
			}
		})
	}
}

// TestMapBoundsGoroutines: acquire the semaphore before starting a
// goroutine, so 100 items with limit 3 never have 100 goroutines waiting.
func TestMapBoundsGoroutines(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		base := runtime.NumGoroutine()
		var most atomic.Int64
		in := make([]int, 100)
		Map(3, in, func(int) int {
			if n := int64(runtime.NumGoroutine() - base); n > most.Load() {
				most.Store(n)
			}
			time.Sleep(time.Second)
			return 0
		})
		if m := most.Load(); m > 5 {
			t.Errorf("Map(3, 100 items) had %d extra goroutines alive, want at most 3 (acquire before go)", m)
		}
	})
}
