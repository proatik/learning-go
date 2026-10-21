package counts

import (
	"maps"
	"testing"
)

func TestMergeCounts(t *testing.T) {
	tests := []struct {
		a, b map[string]int
		want map[string]int
	}{
		{
			map[string]int{"go": 2, "rust": 1},
			map[string]int{"go": 3, "zig": 4},
			map[string]int{"go": 5, "rust": 1, "zig": 4},
		},
		{map[string]int{"x": 1}, nil, map[string]int{"x": 1}},
		{nil, map[string]int{"y": 2}, map[string]int{"y": 2}},
		{map[string]int{"a": 1, "b": 2}, map[string]int{"a": 1, "b": 2}, map[string]int{"a": 2, "b": 4}},
		{map[string]int{"neg": -3}, map[string]int{"neg": 3}, map[string]int{"neg": 0}},
		{map[string]int{}, map[string]int{}, map[string]int{}},
		{nil, nil, map[string]int{}},
	}
	for _, tt := range tests {
		a, b := maps.Clone(tt.a), maps.Clone(tt.b)
		got := MergeCounts(a, b)
		if got == nil {
			t.Errorf("MergeCounts(%v, %v) = nil, want an initialized map", tt.a, tt.b)
			continue
		}
		if !maps.Equal(got, tt.want) {
			t.Errorf("MergeCounts(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
		if !maps.Equal(a, tt.a) || !maps.Equal(b, tt.b) {
			t.Errorf("MergeCounts modified its inputs: (%v, %v) became (%v, %v)", tt.a, tt.b, a, b)
		}
	}
}

func TestMergeCountsReturnsNewMap(t *testing.T) {
	a := map[string]int{"go": 1}
	got := MergeCounts(a, nil)
	if got == nil {
		t.Fatal("MergeCounts(a, nil) = nil, want an initialized map")
	}
	got["go"] = 100
	got["new"] = 1
	if a["go"] != 1 || len(a) != 1 {
		t.Errorf("writing to the result changed a: a = %v; return a new map, not a itself", a)
	}
}

func TestPrune(t *testing.T) {
	tests := []struct {
		m        map[string]int
		minCount int
		want     map[string]int
	}{
		{map[string]int{"a": 1, "b": 5, "c": 3}, 3, map[string]int{"b": 5, "c": 3}},
		{map[string]int{"a": 1}, 0, map[string]int{"a": 1}},
		{map[string]int{"a": 1, "b": 2}, 10, map[string]int{}},
		{map[string]int{"a": -1, "b": 0, "c": 1}, 0, map[string]int{"b": 0, "c": 1}},
		{map[string]int{}, 1, map[string]int{}},
	}
	for _, tt := range tests {
		m := maps.Clone(tt.m)
		Prune(m, tt.minCount)
		if !maps.Equal(m, tt.want) {
			t.Errorf("Prune(%v, %d) left %v, want %v", tt.m, tt.minCount, m, tt.want)
		}
	}
}

func TestPruneNil(t *testing.T) {
	Prune(nil, 5) // deleting from a nil map is a no-op: must not panic
}

func ExamplePrintCounts() {
	PrintCounts(map[string]int{"zig": 4, "go": 5, "rust": 1})
	// Output:
	// go: 5
	// rust: 1
	// zig: 4
}

func ExamplePrintCounts_many() {
	m := map[string]int{}
	for _, w := range []string{"e", "d", "c", "b", "a", "a", "c", "e", "e"} {
		m[w]++
	}
	PrintCounts(m)
	// Output:
	// a: 2
	// b: 1
	// c: 2
	// d: 1
	// e: 3
}

func ExamplePrintCounts_empty() {
	PrintCounts(map[string]int{})
	PrintCounts(nil)
	// Output:
	// (none)
	// (none)
}

func ExampleMergeCounts() {
	PrintCounts(MergeCounts(map[string]int{"go": 2, "rust": 1}, map[string]int{"go": 3, "zig": 4}))
	// Output:
	// go: 5
	// rust: 1
	// zig: 4
}
