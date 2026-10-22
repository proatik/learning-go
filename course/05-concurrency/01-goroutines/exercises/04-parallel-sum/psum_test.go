package psum

import (
	"slices"
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

func TestSum(t *testing.T) {
	tests := []struct {
		nums     []int
		k        int
		partials []int
		total    int
	}{
		{[]int{1, 2, 3, 4, 5, 6}, 3, []int{3, 7, 11}, 21},
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 3, []int{6, 15, 34}, 55},
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 4, []int{3, 12, 13, 27}, 55},
		{[]int{5, -5, 7}, 1, []int{7}, 7},
		{[]int{1, 2}, 10, []int{1, 2}, 3},
		{[]int{1, 2, 3}, 3, []int{1, 2, 3}, 6},
		{[]int{4, 4}, 0, []int{8}, 8},
		{[]int{4, 4}, -2, []int{8}, 8},
		{[]int{-3, -4, -5, -6}, 2, []int{-7, -11}, -18},
		{[]int{9}, 5, []int{9}, 9},
	}
	for _, tt := range tests {
		var total int
		var partials []int
		within(t, "Sum", func() { total, partials = Sum(tt.nums, tt.k) })
		if !slices.Equal(partials, tt.partials) {
			t.Errorf("Sum(%v, %d) partials = %v, want %v", tt.nums, tt.k, partials, tt.partials)
		}
		if total != tt.total {
			t.Errorf("Sum(%v, %d) total = %d, want %d", tt.nums, tt.k, total, tt.total)
		}
	}
}

func TestEmpty(t *testing.T) {
	for _, nums := range [][]int{nil, {}} {
		var total int
		var partials []int
		within(t, "Sum(empty)", func() { total, partials = Sum(nums, 3) })
		if total != 0 || len(partials) != 0 {
			t.Errorf("Sum(%v, 3) = (%d, %v), want (0, [])", nums, total, partials)
		}
	}
}

// slowChunks swaps in a sumChunk that takes 1s per call and counts calls.
func slowChunks(t *testing.T) *atomic.Int32 {
	var calls atomic.Int32
	orig := sumChunk
	t.Cleanup(func() { sumChunk = orig })
	sumChunk = func(part []int) int {
		calls.Add(1)
		time.Sleep(time.Second)
		return orig(part)
	}
	return &calls
}

func TestChunksRunConcurrently(t *testing.T) {
	for _, k := range []int{1, 2, 4, 7} {
		synctest.Test(t, func(t *testing.T) {
			calls := slowChunks(t)
			nums := make([]int, 70)
			for i := range nums {
				nums[i] = i + 1
			}
			start := time.Now()
			total, _ := Sum(nums, k)
			elapsed := time.Since(start)
			if total != 70*71/2 {
				t.Errorf("Sum(1..70, %d) total = %d, want %d", k, total, 70*71/2)
			}
			if c := int(calls.Load()); c != k {
				t.Errorf("Sum(1..70, %d): sumChunk called %d times, want %d", k, c, k)
			}
			if elapsed != time.Second {
				t.Errorf("Sum(1..70, %d) with 1s chunks took %v, want 1s", k, elapsed)
			}
		})
	}
}

func TestLarge(t *testing.T) {
	nums := make([]int, 100_000)
	want := 0
	for i := range nums {
		nums[i] = i%1000 - 500
		want += nums[i]
	}
	for _, k := range []int{1, 8, 1000} {
		var total int
		var partials []int
		within(t, "Sum(100000 ints)", func() { total, partials = Sum(nums, k) })
		if total != want || len(partials) != k {
			t.Errorf("Sum(100000 ints, %d) = (%d, %d partials), want (%d, %d partials)", k, total, len(partials), want, k)
		}
	}
}
