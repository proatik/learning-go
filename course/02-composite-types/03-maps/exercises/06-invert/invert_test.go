package invert

import (
	"maps"
	"reflect"
	"slices"
	"testing"
)

func TestPositions(t *testing.T) {
	tests := []struct {
		in   []string
		want map[string][]int
	}{
		{
			[]string{"to", "be", "or", "not", "to", "be"},
			map[string][]int{"to": {0, 4}, "be": {1, 5}, "or": {2}, "not": {3}},
		},
		{[]string{"a"}, map[string][]int{"a": {0}}},
		{[]string{"x", "x", "x"}, map[string][]int{"x": {0, 1, 2}}},
		{[]string{"", "b", ""}, map[string][]int{"": {0, 2}, "b": {1}}},
		{[]string{}, map[string][]int{}},
		{nil, map[string][]int{}},
	}
	for _, tt := range tests {
		in := slices.Clone(tt.in)
		got := Positions(in)
		if got == nil {
			t.Errorf("Positions(%q) = nil, want an initialized map", tt.in)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Positions(%q) = %v, want %v", tt.in, got, tt.want)
		}
		if !slices.Equal(in, tt.in) {
			t.Errorf("Positions modified its input: %q became %q", tt.in, in)
		}
	}
}

func TestInvert(t *testing.T) {
	tests := []struct {
		in   map[string]int
		want map[int][]string
	}{
		{
			map[string]int{"alice": 90, "bob": 85, "carol": 90},
			map[int][]string{85: {"bob"}, 90: {"alice", "carol"}},
		},
		{map[string]int{"x": 1}, map[int][]string{1: {"x"}}},
		{
			map[string]int{"d": 0, "c": 0, "b": 0, "a": 0, "e": 7},
			map[int][]string{0: {"a", "b", "c", "d"}, 7: {"e"}},
		},
		{
			map[string]int{"zeta": -1, "alpha": -1, "mid": 3, "beta": 3, "omega": 2},
			map[int][]string{-1: {"alpha", "zeta"}, 2: {"omega"}, 3: {"beta", "mid"}},
		},
		{map[string]int{}, map[int][]string{}},
		{nil, map[int][]string{}},
	}
	for _, tt := range tests {
		in := maps.Clone(tt.in)
		// Run several times: unsorted groups would only fail some of the time.
		for range 5 {
			got := Invert(in)
			if got == nil {
				t.Fatalf("Invert(%v) = nil, want an initialized map", tt.in)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Invert(%v) = %v, want %v", tt.in, got, tt.want)
				break
			}
		}
		if !maps.Equal(in, tt.in) {
			t.Errorf("Invert modified its input: %v became %v", tt.in, in)
		}
	}
}
