package sortedset

import (
	"slices"
	"testing"
)

func TestHas(t *testing.T) {
	tests := []struct {
		set  []int
		v    int
		want bool
	}{
		{[]int{1, 3, 5}, 3, true},
		{[]int{1, 3, 5}, 1, true},
		{[]int{1, 3, 5}, 5, true},
		{[]int{1, 3, 5}, 4, false},
		{[]int{1, 3, 5}, 0, false},
		{[]int{1, 3, 5}, 6, false},
		{[]int{-2, 0, 2}, -2, true},
		{nil, 0, false},
	}
	for _, tt := range tests {
		if got := Has(tt.set, tt.v); got != tt.want {
			t.Errorf("Has(%v, %d) = %v, want %v", tt.set, tt.v, got, tt.want)
		}
	}
}

func TestInsertSorted(t *testing.T) {
	tests := []struct {
		s    []int
		v    int
		want []int
	}{
		{[]int{1, 3, 5}, 4, []int{1, 3, 4, 5}},
		{[]int{1, 3, 5}, 3, []int{1, 3, 3, 5}},
		{[]int{1, 3, 5}, 0, []int{0, 1, 3, 5}},
		{[]int{1, 3, 5}, 9, []int{1, 3, 5, 9}},
		{[]int{2, 2}, 2, []int{2, 2, 2}},
		{nil, 7, []int{7}},
	}
	for _, tt := range tests {
		orig := slices.Clone(tt.s)
		if got := InsertSorted(tt.s, tt.v); !slices.Equal(got, tt.want) {
			t.Errorf("InsertSorted(%v, %d) = %v, want %v", orig, tt.v, got, tt.want)
		}
	}
}

func TestInsertSortedBuildsSortedSlice(t *testing.T) {
	var s []int
	for _, v := range []int{5, 3, 8, 1, 3, 9, 0, 5, -2} {
		s = InsertSorted(s, v)
	}
	want := []int{-2, 0, 1, 3, 3, 5, 5, 8, 9}
	if !slices.Equal(s, want) {
		t.Errorf("inserting 5 3 8 1 3 9 0 5 -2 one by one gave %v, want %v", s, want)
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		set       []int
		v         int
		want      []int
		wantAdded bool
	}{
		{[]int{1, 3, 5}, 0, []int{0, 1, 3, 5}, true},
		{[]int{1, 3, 5}, 4, []int{1, 3, 4, 5}, true},
		{[]int{1, 3, 5}, 6, []int{1, 3, 5, 6}, true},
		{[]int{1, 3, 5}, 5, []int{1, 3, 5}, false},
		{[]int{1, 3, 5}, 1, []int{1, 3, 5}, false},
		{nil, 2, []int{2}, true},
	}
	for _, tt := range tests {
		orig := slices.Clone(tt.set)
		got, added := Add(tt.set, tt.v)
		if !slices.Equal(got, tt.want) || added != tt.wantAdded {
			t.Errorf("Add(%v, %d) = %v, %v; want %v, %v", orig, tt.v, got, added, tt.want, tt.wantAdded)
		}
	}
}

func TestAddBuildsSet(t *testing.T) {
	var set []int
	count := 0
	for _, v := range []int{4, 2, 4, 9, 2, 1, 9, 4} {
		var added bool
		set, added = Add(set, v)
		if added {
			count++
		}
	}
	if want := []int{1, 2, 4, 9}; !slices.Equal(set, want) {
		t.Errorf("adding 4 2 4 9 2 1 9 4 gave %v, want %v", set, want)
	}
	if count != 4 {
		t.Errorf("Add reported %d successful adds, want 4", count)
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		set       []int
		v         int
		want      []int
		wantFound bool
	}{
		{[]int{1, 3, 5}, 3, []int{1, 5}, true},
		{[]int{1, 3, 5}, 1, []int{3, 5}, true},
		{[]int{1, 3, 5}, 5, []int{1, 3}, true},
		{[]int{1, 3, 5}, 2, []int{1, 3, 5}, false},
		{[]int{1, 3, 3, 5}, 3, []int{1, 3, 5}, true},
		{[]int{7}, 7, []int{}, true},
		{nil, 7, []int{}, false},
	}
	for _, tt := range tests {
		orig := slices.Clone(tt.set)
		got, found := Remove(tt.set, tt.v)
		if !slices.Equal(got, tt.want) || found != tt.wantFound {
			t.Errorf("Remove(%v, %d) = %v, %v; want %v, %v", orig, tt.v, got, found, tt.want, tt.wantFound)
		}
	}
}

func TestAddThenRemove(t *testing.T) {
	var set []int
	for _, v := range []int{10, 20, 30, 40} {
		set, _ = Add(set, v)
	}
	set, _ = Remove(set, 20)
	set, _ = Remove(set, 99)
	set, _ = Add(set, 25)
	set, _ = Remove(set, 10)
	if want := []int{25, 30, 40}; !slices.Equal(set, want) {
		t.Errorf("after adds and removes got %v, want %v", set, want)
	}
	if !Has(set, 30) || Has(set, 20) {
		t.Errorf("Has on %v: Has(30) = %v, Has(20) = %v; want true, false", set, Has(set, 30), Has(set, 20))
	}
}
