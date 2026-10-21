package people

import (
	"slices"
	"sort"
	"testing"
)

var (
	_ sort.Interface = ByAge(nil)
	_ sort.Interface = ByName(nil)
)

func sample() []Person {
	return []Person{{"Cy", 40}, {"Ann", 31}, {"Bob", 25}, {"Al", 31}}
}

func ages(ps []Person) []int {
	var out []int
	for _, p := range ps {
		out = append(out, p.Age)
	}
	return out
}

func TestByAgeMethods(t *testing.T) {
	ps := sample()
	a := ByAge(ps)
	if got := a.Len(); got != 4 {
		t.Errorf("ByAge(%v).Len() = %d, want 4", ps, got)
	}
	if got := ByAge(nil).Len(); got != 0 {
		t.Errorf("ByAge(nil).Len() = %d, want 0", got)
	}
	lessTests := []struct {
		i, j int
		want bool
	}{{2, 1, true}, {1, 2, false}, {1, 3, false}, {3, 1, false}, {1, 0, true}}
	for _, tt := range lessTests {
		if got := a.Less(tt.i, tt.j); got != tt.want {
			t.Errorf("ByAge(%v).Less(%d, %d) = %v, want %v", ps, tt.i, tt.j, got, tt.want)
		}
	}
	a.Swap(0, 2)
	if want := []Person{{"Bob", 25}, {"Ann", 31}, {"Cy", 40}, {"Al", 31}}; !slices.Equal(ps, want) {
		t.Errorf("ByAge.Swap(0, 2) = %v, want %v", ps, want)
	}
}

func TestByNameMethods(t *testing.T) {
	ps := sample()
	n := ByName(ps)
	if got := n.Len(); got != 4 {
		t.Errorf("ByName(%v).Len() = %d, want 4", ps, got)
	}
	lessTests := []struct {
		i, j int
		want bool
	}{{3, 1, true}, {1, 3, false}, {1, 2, true}, {0, 2, false}, {2, 2, false}}
	for _, tt := range lessTests {
		if got := n.Less(tt.i, tt.j); got != tt.want {
			t.Errorf("ByName(%v).Less(%d, %d) = %v, want %v", ps, tt.i, tt.j, got, tt.want)
		}
	}
	n.Swap(1, 3)
	if want := []Person{{"Cy", 40}, {"Al", 31}, {"Bob", 25}, {"Ann", 31}}; !slices.Equal(ps, want) {
		t.Errorf("ByName.Swap(1, 3) = %v, want %v", ps, want)
	}
}

func TestSortByAge(t *testing.T) {
	ps := sample()
	sort.Sort(ByAge(ps))
	if got, want := ages(ps), []int{25, 31, 31, 40}; !slices.Equal(got, want) {
		t.Errorf("sort.Sort(ByAge(%v)) ages = %v, want %v", sample(), got, want)
	}
}

func TestSortByName(t *testing.T) {
	ps := sample()
	sort.Sort(ByName(ps))
	want := []Person{{"Al", 31}, {"Ann", 31}, {"Bob", 25}, {"Cy", 40}}
	if !slices.Equal(ps, want) {
		t.Errorf("sort.Sort(ByName(%v)) = %v, want %v", sample(), ps, want)
	}
}

func TestSortOldestFirst(t *testing.T) {
	tests := []struct {
		in, want []Person
	}{
		{sample(), []Person{{"Cy", 40}, {"Ann", 31}, {"Al", 31}, {"Bob", 25}}},
		{nil, nil},
		{[]Person{{"Solo", 7}}, []Person{{"Solo", 7}}},
		{
			[]Person{{"a", 1}, {"b", 2}, {"c", 1}, {"d", 2}, {"e", 1}, {"f", 3}},
			[]Person{{"f", 3}, {"b", 2}, {"d", 2}, {"a", 1}, {"c", 1}, {"e", 1}},
		},
	}
	for _, tt := range tests {
		got := slices.Clone(tt.in)
		SortOldestFirst(got)
		if !slices.Equal(got, tt.want) {
			t.Errorf("SortOldestFirst(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestSortAgeThenName(t *testing.T) {
	tests := []struct {
		in, want []Person
	}{
		{sample(), []Person{{"Bob", 25}, {"Al", 31}, {"Ann", 31}, {"Cy", 40}}},
		{nil, nil},
		{[]Person{}, []Person{}},
		{
			[]Person{{"zed", 1}, {"amy", 2}, {"bob", 1}, {"Amy", 2}},
			[]Person{{"bob", 1}, {"zed", 1}, {"Amy", 2}, {"amy", 2}},
		},
	}
	for _, tt := range tests {
		got := slices.Clone(tt.in)
		SortAgeThenName(got)
		if !slices.Equal(got, tt.want) {
			t.Errorf("SortAgeThenName(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
}
