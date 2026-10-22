package psort

import (
	"cmp"
	"math/rand/v2"
	"slices"
	"strings"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

func TestSortSmall(t *testing.T) {
	inputs := [][]int{
		nil,
		{},
		{7},
		{2, 1},
		{1, 2},
		{5, 2, 8, 1, 9, 3},
		{1, 2, 3, 4, 5, 6, 7},
		{7, 6, 5, 4, 3, 2, 1},
		{3, 3, 1, 3, 1, 1, 2},
		{-5, 0, 10, -20, 3, 0, -5},
	}
	for _, in := range inputs {
		want := slices.Sorted(slices.Values(in))
		for _, depth := range []int{-1, 0, 1, 2, 3, 10} {
			got := slices.Clone(in)
			Sort(got, depth)
			if !slices.Equal(got, want) {
				t.Errorf("Sort(%v, %d) = %v, want %v", in, depth, got, want)
			}
		}
	}
}

func TestSortStrings(t *testing.T) {
	in := strings.Fields("pear apple fig banana cherry apple date grape elderberry kiwi")
	want := slices.Sorted(slices.Values(in))
	got := slices.Clone(in)
	Sort(got, 2)
	if !slices.Equal(got, want) {
		t.Errorf("Sort(%q, 2) = %q, want %q", in, got, want)
	}
}

func TestSortFuncStructs(t *testing.T) {
	type person struct {
		Name string
		Age  int
	}
	byAge := func(a, b person) int {
		return cmp.Or(cmp.Compare(a.Age, b.Age), strings.Compare(a.Name, b.Name))
	}
	in := []person{{"ann", 31}, {"bob", 25}, {"cat", 31}, {"dan", 19}, {"eve", 25}, {"fay", 40}}
	want := slices.Clone(in)
	slices.SortFunc(want, byAge)
	got := slices.Clone(in)
	SortFunc(got, byAge, 2)
	if !slices.Equal(got, want) {
		t.Errorf("SortFunc(%v, byAge, 2) = %v, want %v", in, got, want)
	}
	desc := func(a, b int) int { return cmp.Compare(b, a) }
	nums := []int{4, 1, 3, 5, 2}
	SortFunc(nums, desc, 1)
	if want := []int{5, 4, 3, 2, 1}; !slices.Equal(nums, want) {
		t.Errorf("SortFunc([4 1 3 5 2], descending, 1) = %v, want %v", nums, want)
	}
}

func TestSortLarge(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for _, size := range []int{1000, 4097, 200_000} {
		in := make([]int, size)
		for i := range in {
			in[i] = r.IntN(size) - size/2
		}
		want := slices.Clone(in)
		slices.Sort(want)
		for _, depth := range []int{1, 4, 8} {
			got := slices.Clone(in)
			Sort(got, depth)
			if i := firstDiff(got, want); i >= 0 {
				t.Fatalf("Sort(%d random ints, %d): element %d = %d, want %d", size, depth, i, got[i], want[i])
			}
		}
	}
}

func firstDiff(a, b []int) int {
	for i := range a {
		if a[i] != b[i] {
			return i
		}
	}
	if len(a) != len(b) {
		return min(len(a), len(b))
	}
	return -1
}

// TestParallelism counts the comparisons running at the same moment: each
// one sleeps 1ms of fake time, so pieces sorted concurrently overlap.
func TestParallelism(t *testing.T) {
	for _, tt := range []struct{ depth, want int }{{0, 1}, {1, 2}, {2, 4}, {3, 8}} {
		synctest.Test(t, func(t *testing.T) {
			r := rand.New(rand.NewPCG(3, 4))
			in := make([]int, 128)
			for i := range in {
				in[i] = r.IntN(1000)
			}
			var running, maxRunning, calls atomic.Int64
			slow := func(a, b int) int {
				calls.Add(1)
				n := running.Add(1)
				for {
					m := maxRunning.Load()
					if n <= m || maxRunning.CompareAndSwap(m, n) {
						break
					}
				}
				time.Sleep(time.Millisecond)
				running.Add(-1)
				return cmp.Compare(a, b)
			}
			got := slices.Clone(in)
			SortFunc(got, slow, tt.depth)
			if calls.Load() == 0 {
				t.Fatalf("SortFunc(128 ints, cmp, %d) never called cmp", tt.depth)
			}
			if !slices.IsSorted(got) {
				t.Fatalf("SortFunc(128 ints, cmp, %d) result is not sorted", tt.depth)
			}
			if m := maxRunning.Load(); m != int64(tt.want) {
				t.Errorf("SortFunc(128 ints, cmp, %d): max cmp calls at once = %d, want %d", tt.depth, m, tt.want)
			}
		})
	}
}
