package pipeline_test

import (
	"iter"
	"slices"
	"strconv"
	"strings"
	"testing"

	pipeline "learngo/course/04-generics-packages-testing/03-iterators/exercises/04-filter-map"
)

// collect ranges over seq and returns its values, breaking after max values
// when max > 0 (and always after 10000, so a runaway sequence can't hang the
// test). A panic during the loop is returned instead of crashing.
func collect[V any](seq iter.Seq[V], max int) (got []V, panicked any) {
	defer func() { panicked = recover() }()
	for v := range seq {
		got = append(got, v)
		if len(got) == max || len(got) == 10000 {
			break
		}
	}
	return got, nil
}

// source yields 1, 2, 3, … forever and counts how many values it produced.
func source(pulled *int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 1; ; i++ {
			*pulled++
			if !yield(i) {
				return
			}
		}
	}
}

func isEven(n int) bool { return n%2 == 0 }

func TestFilter(t *testing.T) {
	tests := []struct {
		in   []int
		want []int
	}{
		{[]int{1, 2, 3, 4}, []int{2, 4}},
		{[]int{1, 3, 5}, nil},
		{[]int{2, 2, 8}, []int{2, 2, 8}},
		{nil, nil},
		{[]int{-4, -3, 0}, []int{-4, 0}},
	}
	for _, tt := range tests {
		seq := pipeline.Filter(slices.Values(tt.in), isEven)
		if seq == nil {
			t.Fatalf("Filter(%v, isEven) = nil, want a non-nil iterator", tt.in)
		}
		got, p := collect(seq, -1)
		if p != nil {
			t.Fatalf("range over Filter(%v, isEven) panicked: %v", tt.in, p)
		}
		if !slices.Equal(got, tt.want) {
			t.Errorf("Filter(%v, isEven) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestMap(t *testing.T) {
	got, p := collect(pipeline.Map(slices.Values([]string{"go", "hi", "ñu"}), strings.ToUpper), -1)
	if p != nil {
		t.Fatalf("range over Map(..., strings.ToUpper) panicked: %v", p)
	}
	if want := []string{"GO", "HI", "ÑU"}; !slices.Equal(got, want) {
		t.Errorf("Map([go hi ñu], strings.ToUpper) = %q, want %q", got, want)
	}

	lens, _ := collect(pipeline.Map(slices.Values([]string{"a", "bcd", ""}), func(s string) int { return len(s) }), -1)
	if want := []int{1, 3, 0}; !slices.Equal(lens, want) {
		t.Errorf("Map([a bcd \"\"], len) = %v, want %v", lens, want)
	}

	empty, _ := collect(pipeline.Map(slices.Values([]int(nil)), strconv.Itoa), -1)
	if len(empty) != 0 {
		t.Errorf("Map([], strconv.Itoa) = %q, want []", empty)
	}
}

func TestIsLazy(t *testing.T) {
	pulled, calls := 0, 0
	f := pipeline.Filter(source(&pulled), func(n int) bool { calls++; return true })
	m := pipeline.Map(f, func(n int) int { calls++; return n })
	if pulled != 0 || calls != 0 {
		t.Errorf("building Map(Filter(source)) before ranging: source values + callback calls = %d, want 0", pulled+calls)
	}
	_ = m
}

func TestCallbacksOncePerValue(t *testing.T) {
	pulled := 0
	keepCalls := map[int]int{}
	mapCalls := map[int]int{}
	evens := pipeline.Filter(source(&pulled), func(n int) bool { keepCalls[n]++; return isEven(n) })
	squares := pipeline.Map(evens, func(n int) int { mapCalls[n]++; return n * n })
	got, p := collect(squares, 3)
	if p != nil {
		t.Fatalf("break after 3 values of Map(Filter(1, 2, 3, …, isEven), square) panicked: %v (stop when yield returns false)", p)
	}
	if want := []int{4, 16, 36}; !slices.Equal(got, want) {
		t.Errorf("first 3 values of Map(Filter(1, 2, 3, …, isEven), square) = %v, want %v", got, want)
	}
	if pulled != 6 {
		t.Errorf("break after 3 values of Map(Filter(1, 2, 3, …, isEven), square): source values produced = %d, want 6", pulled)
	}
	for n := 1; n <= 6; n++ {
		if keepCalls[n] != 1 {
			t.Errorf("keep(%d) calls = %d, want 1", n, keepCalls[n])
		}
	}
	for _, n := range []int{2, 4, 6} {
		if mapCalls[n] != 1 {
			t.Errorf("f(%d) calls = %d, want 1", n, mapCalls[n])
		}
	}
	if len(mapCalls) != 3 {
		t.Errorf("Map's f was called for %d different values, want 3 (only the kept ones)", len(mapCalls))
	}
}

func TestBreak(t *testing.T) {
	pulled := 0
	got, p := collect(pipeline.Filter(source(&pulled), isEven), 1)
	if p != nil {
		t.Fatalf("break after 1 value of Filter(1, 2, 3, …, isEven) panicked: %v (stop when yield returns false)", p)
	}
	if !slices.Equal(got, []int{2}) || pulled != 2 {
		t.Errorf("break after 1 value of Filter(1, 2, 3, …, isEven): source values produced = %d, want 2", pulled)
	}

	pulled = 0
	strs, p := collect(pipeline.Map(source(&pulled), strconv.Itoa), 2)
	if p != nil {
		t.Fatalf("break after 2 values of Map(1, 2, 3, …, strconv.Itoa) panicked: %v (stop when yield returns false)", p)
	}
	if !slices.Equal(strs, []string{"1", "2"}) || pulled != 2 {
		t.Errorf("break after 2 values of Map(1, 2, 3, …, strconv.Itoa): source values produced = %d, want 2", pulled)
	}
}

func TestReduce(t *testing.T) {
	add := func(a, v int) int { return a + v }
	tests := []struct {
		in   []int
		init int
		want int
	}{
		{[]int{1, 2, 3}, 0, 6},
		{[]int{1, 2, 3}, 100, 106},
		{nil, 10, 10},
		{[]int{-5}, 5, 0},
	}
	for _, tt := range tests {
		if got := pipeline.Reduce(slices.Values(tt.in), tt.init, add); got != tt.want {
			t.Errorf("Reduce(%v, %d, add) = %d, want %d", tt.in, tt.init, got, tt.want)
		}
	}

	joined := pipeline.Reduce(slices.Values([]int{1, 2, 3}), ">", func(acc string, v int) string {
		return acc + strconv.Itoa(v)
	})
	if joined != ">123" {
		t.Errorf("Reduce([1 2 3], %q, appendDigit) = %q, want %q", ">", joined, ">123")
	}
}

func TestPipeline(t *testing.T) {
	nums := slices.Values([]int{1, 2, 3, 4, 5, 6})
	square := func(n int) int { return n * n }
	got := pipeline.Reduce(pipeline.Map(pipeline.Filter(nums, isEven), square), 0, func(a, v int) int { return a + v })
	if got != 56 {
		t.Errorf("Reduce(Map(Filter([1..6], isEven), square), 0, add) = %d, want %d", got, 56)
	}
}
