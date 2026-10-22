package pipe

import (
	"iter"
	"slices"
	"strconv"
	"testing"
)

// counted yields vals and counts how many were pulled in *pulled.
func counted[T any](vals []T, pulled *int) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range vals {
			*pulled++
			if !yield(v) {
				return
			}
		}
	}
}

func nums(n int) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = i + 1
	}
	return out
}

func TestMap(t *testing.T) {
	got := slices.Collect(Map(slices.Values([]int{1, 2, 3}), strconv.Itoa))
	if want := []string{"1", "2", "3"}; !slices.Equal(got, want) {
		t.Errorf("Map([1 2 3], Itoa) = %q, want %q", got, want)
	}
	if got := slices.Collect(Map(slices.Values([]int{}), strconv.Itoa)); len(got) != 0 {
		t.Errorf("Map([], Itoa) = %q, want []", got)
	}
}

func TestFilter(t *testing.T) {
	even := func(n int) bool { return n%2 == 0 }
	tests := []struct {
		in, want []int
	}{
		{[]int{1, 2, 3, 4, 5, 6}, []int{2, 4, 6}},
		{[]int{1, 3}, nil},
		{nil, nil},
	}
	for _, tt := range tests {
		if got := slices.Collect(Filter(slices.Values(tt.in), even)); !slices.Equal(got, tt.want) {
			t.Errorf("Filter(%v, even) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestTake(t *testing.T) {
	tests := []struct {
		n          int
		want       []int
		wantPulled int
	}{
		{3, []int{1, 2, 3}, 3},
		{1, []int{1}, 1},
		{0, nil, 0},
		{-2, nil, 0},
		{10, []int{1, 2, 3, 4, 5}, 5},
	}
	for _, tt := range tests {
		pulled := 0
		got := slices.Collect(Take(counted(nums(5), &pulled), tt.n))
		if !slices.Equal(got, tt.want) {
			t.Errorf("Take([1..5], %d) = %v, want %v", tt.n, got, tt.want)
		}
		if pulled != tt.wantPulled {
			t.Errorf("Take([1..5], %d) pulled %d values, want %d", tt.n, pulled, tt.wantPulled)
		}
	}
}

func TestBatch(t *testing.T) {
	tests := []struct {
		in   []int
		size int
		want [][]int
	}{
		{nums(7), 3, [][]int{{1, 2, 3}, {4, 5, 6}, {7}}},
		{nums(4), 2, [][]int{{1, 2}, {3, 4}}},
		{nums(2), 5, [][]int{{1, 2}}},
		{nums(3), 1, [][]int{{1}, {2}, {3}}},
		{nil, 3, nil},
		{nums(3), 0, nil},
	}
	for _, tt := range tests {
		got := slices.Collect(Batch(slices.Values(tt.in), tt.size))
		if !slices.EqualFunc(got, tt.want, slices.Equal) {
			t.Errorf("Batch(%v, %d) = %v, want %v", tt.in, tt.size, got, tt.want)
		}
	}
}

func TestLaziness(t *testing.T) {
	pulled, mapped := 0, 0
	square := func(n int) int { mapped++; return n * n }
	big := func(n int) bool { return n > 10 }
	seq := Take(Filter(Map(counted(nums(100), &pulled), square), big), 2)
	if pulled != 0 || mapped != 0 {
		t.Fatalf("building the pipeline pulled %d and mapped %d values, want 0 and 0 (adapters must be lazy)", pulled, mapped)
	}
	got := slices.Collect(seq)
	if want := []int{16, 25}; !slices.Equal(got, want) {
		t.Errorf("Take(Filter(Map(1..100, square), >10), 2) = %v, want %v", got, want)
	}
	if pulled != 5 || mapped != 5 {
		t.Errorf("pipeline pulled %d and mapped %d values, want 5 and 5", pulled, mapped)
	}
}

func TestEarlyBreak(t *testing.T) {
	pulled := 0
	for range Map(counted(nums(10), &pulled), strconv.Itoa) {
		break
	}
	for v := range Filter(counted(nums(10), &pulled), func(int) bool { return true }) {
		if v == 2 {
			break
		}
	}
	for range Take(counted(nums(10), &pulled), 5) {
		break
	}
	for range Batch(counted(nums(10), &pulled), 2) {
		break
	}
	if want := 1 + 2 + 1 + 2; pulled != want {
		t.Errorf("breaking early from Map, Filter, Take and Batch pulled %d values, want %d", pulled, want)
	}
}

func TestFold(t *testing.T) {
	sum := Fold(slices.Values(nums(4)), 0, func(a, n int) int { return a + n })
	if sum != 10 {
		t.Errorf("Fold([1 2 3 4], 0, +) = %d, want 10", sum)
	}
	s := Fold(slices.Values([]int{1, 2, 3}), "x", func(a string, n int) string { return a + strconv.Itoa(n) })
	if s != "x123" {
		t.Errorf("Fold([1 2 3], %q, concat) = %q, want %q", "x", s, "x123")
	}
	if got := Fold(slices.Values([]int{}), 7, func(a, n int) int { return a + n }); got != 7 {
		t.Errorf("Fold([], 7, +) = %d, want 7", got)
	}
}
