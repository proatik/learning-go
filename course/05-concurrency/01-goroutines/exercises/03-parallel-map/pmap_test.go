package pmap_test

import (
	"slices"
	"strconv"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
	"unicode/utf8"

	pmap "learngo/course/05-concurrency/01-goroutines/exercises/03-parallel-map"
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

func TestOrderAndConcurrency(t *testing.T) {
	tests := []struct {
		in      []int
		want    []int
		elapsed time.Duration
	}{
		{[]int{3, 1, 2}, []int{30, 10, 20}, 3 * time.Second},
		{[]int{1, 2, 3, 4, 5}, []int{10, 20, 30, 40, 50}, 5 * time.Second},
		{[]int{5, 4, 3, 2, 1}, []int{50, 40, 30, 20, 10}, 5 * time.Second},
		{[]int{2, 2, 2, 2}, []int{20, 20, 20, 20}, 2 * time.Second},
		{[]int{0, -1, 1}, []int{0, -10, 10}, time.Second},
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			var calls atomic.Int32
			f := func(n int) int {
				calls.Add(1)
				time.Sleep(time.Duration(n) * time.Second) // negative: no sleep
				return n * 10
			}
			start := time.Now()
			got := pmap.Map(tt.in, f)
			elapsed := time.Since(start)
			if !slices.Equal(got, tt.want) {
				t.Errorf("Map(%v, sleep n then n*10) = %v, want %v", tt.in, got, tt.want)
			}
			if c := int(calls.Load()); c != len(tt.in) {
				t.Errorf("Map(%v): f called %d times, want %d", tt.in, c, len(tt.in))
			}
			if elapsed != tt.elapsed {
				t.Errorf("Map(%v, sleep n then n*10) took %v, want %v", tt.in, elapsed, tt.elapsed)
			}
		})
	}
}

func TestOtherTypes(t *testing.T) {
	var got []int
	in := []string{"go", "héllo", "", "世界"}
	within(t, "Map", func() { got = pmap.Map(in, utf8.RuneCountInString) })
	if want := []int{2, 5, 0, 2}; !slices.Equal(got, want) {
		t.Errorf("Map(%q, RuneCount) = %v, want %v", in, got, want)
	}

	var strs []string
	nums := []int{7, -2, 0}
	within(t, "Map", func() { strs = pmap.Map(nums, strconv.Itoa) })
	if want := []string{"7", "-2", "0"}; !slices.Equal(strs, want) {
		t.Errorf("Map(%v, strconv.Itoa) = %q, want %q", nums, strs, want)
	}
}

func TestEmpty(t *testing.T) {
	var calls atomic.Int32
	f := func(n int) int { calls.Add(1); return n }
	for _, in := range [][]int{nil, {}} {
		var got []int
		within(t, "Map(empty)", func() { got = pmap.Map(in, f) })
		if len(got) != 0 {
			t.Errorf("Map(%v, f) = %v, want []", in, got)
		}
	}
	if c := calls.Load(); c != 0 {
		t.Errorf("Map(empty, f): f called %d times, want 0", c)
	}
}

func TestMany(t *testing.T) {
	in := make([]int, 2000)
	want := make([]int, len(in))
	for i := range in {
		in[i] = i
		want[i] = i * i
	}
	var got []int
	within(t, "Map(2000 ints)", func() { got = pmap.Map(in, func(n int) int { return n * n }) })
	if !slices.Equal(got, want) {
		t.Errorf("Map(0..1999, square) = %d results (first wrong index %d), want %d squares in order",
			len(got), firstDiff(got, want), len(want))
	}
}

// firstDiff returns the first index where a and b differ.
func firstDiff(a, b []int) int {
	for i := range min(len(a), len(b)) {
		if a[i] != b[i] {
			return i
		}
	}
	return min(len(a), len(b))
}
