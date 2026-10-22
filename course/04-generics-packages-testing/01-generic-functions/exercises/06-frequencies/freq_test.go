package freq_test

import (
	"maps"
	"slices"
	"testing"

	fq "learngo/course/04-generics-packages-testing/01-generic-functions/exercises/06-frequencies"
)

type Point struct{ X, Y int }

type Grade string

func TestFrequencies(t *testing.T) {
	tests := []struct {
		in   []string
		want map[string]int
	}{
		{[]string{"a", "b", "a"}, map[string]int{"a": 2, "b": 1}},
		{[]string{"x"}, map[string]int{"x": 1}},
		{[]string{"", "", "é", "e"}, map[string]int{"": 2, "é": 1, "e": 1}},
		{[]string{"Go", "go"}, map[string]int{"Go": 1, "go": 1}},
		{[]string{}, map[string]int{}},
	}
	for _, tt := range tests {
		got := fq.Frequencies(tt.in)
		if !maps.Equal(got, tt.want) {
			t.Errorf("Frequencies(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}

	pts := []Point{{1, 2}, {1, 2}, {0, 0}}
	if got, want := fq.Frequencies(pts), map[Point]int{{1, 2}: 2, {0, 0}: 1}; !maps.Equal(got, want) {
		t.Errorf("Frequencies(%v) = %v, want %v", pts, got, want)
	}
}

func TestFrequenciesNeverNil(t *testing.T) {
	for _, in := range [][]int{nil, {}} {
		got := fq.Frequencies(in)
		if got == nil {
			t.Errorf("Frequencies(%v) = nil map, want an empty non-nil map", in)
		}
	}
}

func TestSortedKeys(t *testing.T) {
	m := map[string]int{"b": 1, "a": 2, "c": 0, "B": 9}
	if got, want := fq.SortedKeys(m), []string{"B", "a", "b", "c"}; !slices.Equal(got, want) {
		t.Errorf("SortedKeys(%v) = %q, want %q", m, got, want)
	}

	ids := map[int]bool{30: true, -1: false, 7: true, 0: true}
	if got, want := fq.SortedKeys(ids), []int{-1, 0, 7, 30}; !slices.Equal(got, want) {
		t.Errorf("SortedKeys(%v) = %v, want %v", ids, got, want)
	}

	grades := map[Grade][]string{"B": {"bob"}, "A": {"ann", "cy"}}
	if got, want := fq.SortedKeys(grades), []Grade{"A", "B"}; !slices.Equal(got, want) {
		t.Errorf("SortedKeys(%v) = %q, want %q", grades, got, want)
	}

	floats := map[float64]struct{}{2.5: {}, -1.5: {}, 0: {}}
	if got, want := fq.SortedKeys(floats), []float64{-1.5, 0, 2.5}; !slices.Equal(got, want) {
		t.Errorf("SortedKeys(%v) = %v, want %v", floats, got, want)
	}

	if got := fq.SortedKeys(map[int]bool{}); len(got) != 0 {
		t.Errorf("SortedKeys(map[]) = %v, want []", got)
	}
	if got := fq.SortedKeys[string, int](nil); len(got) != 0 {
		t.Errorf("SortedKeys(nil) = %v, want []", got)
	}
}

func TestSortedKeysStable(t *testing.T) {
	m := map[int]string{}
	for i := 100; i > 0; i-- {
		m[i] = "x"
	}
	for range 5 {
		got := fq.SortedKeys(m)
		if len(got) != 100 || !slices.IsSorted(got) {
			t.Fatalf("SortedKeys(1..100) = %v, want 100 keys in ascending order", got)
		}
	}
}

func TestMostCommon(t *testing.T) {
	tests := []struct {
		in        []string
		want      string
		wantCount int
		wantOK    bool
	}{
		{[]string{"go", "rust", "go"}, "go", 2, true},
		{[]string{"zig"}, "zig", 1, true},
		{[]string{"b", "a"}, "a", 1, true},
		{[]string{"c", "b", "c", "b", "a"}, "b", 2, true},
		{[]string{"", "", "x"}, "", 2, true},
		{[]string{}, "", 0, false},
		{nil, "", 0, false},
	}
	for _, tt := range tests {
		got, n, ok := fq.MostCommon(tt.in)
		if got != tt.want || n != tt.wantCount || ok != tt.wantOK {
			t.Errorf("MostCommon(%q) = %q, %d, %t, want %q, %d, %t", tt.in, got, n, ok, tt.want, tt.wantCount, tt.wantOK)
		}
	}

	ints := []struct {
		in        []int
		want      int
		wantCount int
		wantOK    bool
	}{
		{[]int{3, 1, 3, 1, 2}, 1, 2, true},
		{[]int{5, 5, 5, -2, -2}, 5, 3, true},
		{[]int{0, 0, 7}, 0, 2, true},
		{[]int{9, -9}, -9, 1, true},
		{[]int{}, 0, 0, false},
	}
	for _, tt := range ints {
		got, n, ok := fq.MostCommon(tt.in)
		if got != tt.want || n != tt.wantCount || ok != tt.wantOK {
			t.Errorf("MostCommon(%v) = %d, %d, %t, want %d, %d, %t", tt.in, got, n, ok, tt.want, tt.wantCount, tt.wantOK)
		}
	}
}
