package merge_test

import (
	"cmp"
	"slices"
	"strings"
	"testing"

	mg "learngo/course/04-generics-packages-testing/01-generic-functions/exercises/09-merge-sorted"
)

type Person struct {
	Name string
	Age  int
}

func byAge(x, y Person) int { return cmp.Compare(x.Age, y.Age) }

func TestMergeInts(t *testing.T) {
	tests := []struct {
		a, b, want []int
	}{
		{[]int{1, 4, 9}, []int{2, 3, 10}, []int{1, 2, 3, 4, 9, 10}},
		{[]int{1, 2, 3}, []int{4, 5}, []int{1, 2, 3, 4, 5}},
		{[]int{4, 5}, []int{1, 2, 3}, []int{1, 2, 3, 4, 5}},
		{[]int{1, 1, 2}, []int{1, 2, 2}, []int{1, 1, 1, 2, 2, 2}},
		{[]int{-5, 0}, []int{-7, 0, 7}, []int{-7, -5, 0, 0, 7}},
		{[]int{}, []int{5}, []int{5}},
		{[]int{5}, nil, []int{5}},
		{nil, nil, []int{}},
	}
	for _, tt := range tests {
		if got := mg.Merge(tt.a, tt.b, cmp.Compare[int]); !slices.Equal(got, tt.want) {
			t.Errorf("Merge(%v, %v, cmp.Compare) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
		if got := mg.MergeOrdered(tt.a, tt.b); !slices.Equal(got, tt.want) {
			t.Errorf("MergeOrdered(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestMergeOrderedStrings(t *testing.T) {
	a, b := []string{"ant", "fox"}, []string{"bee"}
	if got, want := mg.MergeOrdered(a, b), []string{"ant", "bee", "fox"}; !slices.Equal(got, want) {
		t.Errorf("MergeOrdered(%q, %q) = %q, want %q", a, b, got, want)
	}
	a, b = []string{"Zed", "apple"}, []string{"Bob", "zoo"}
	if got, want := mg.MergeOrdered(a, b), []string{"Bob", "Zed", "apple", "zoo"}; !slices.Equal(got, want) {
		t.Errorf("MergeOrdered(%q, %q) = %q, want %q", a, b, got, want)
	}
	fa, fb := []float64{-1.5, 2.5}, []float64{0, 2.5}
	if got, want := mg.MergeOrdered(fa, fb), []float64{-1.5, 0, 2.5, 2.5}; !slices.Equal(got, want) {
		t.Errorf("MergeOrdered(%v, %v) = %v, want %v", fa, fb, got, want)
	}
}

func TestMergeStable(t *testing.T) {
	a := []Person{{"ann", 30}, {"cy", 30}, {"eve", 50}}
	b := []Person{{"bob", 20}, {"dan", 30}, {"fay", 50}}
	got := mg.Merge(a, b, byAge)
	want := []Person{{"bob", 20}, {"ann", 30}, {"cy", 30}, {"dan", 30}, {"eve", 50}, {"fay", 50}}
	if !slices.Equal(got, want) {
		t.Errorf("Merge(%v, %v, byAge) = %v, want %v", a, b, got, want)
	}

	got = mg.Merge(b, a, byAge)
	want = []Person{{"bob", 20}, {"dan", 30}, {"ann", 30}, {"cy", 30}, {"fay", 50}, {"eve", 50}}
	if !slices.Equal(got, want) {
		t.Errorf("Merge(%v, %v, byAge) = %v, want %v", b, a, got, want)
	}
}

func TestMergeCustomComparator(t *testing.T) {
	byLen := func(x, y string) int { return cmp.Compare(len(x), len(y)) }
	a, b := []string{"go", "rust"}, []string{"c", "zig", "java"}
	if got, want := mg.Merge(a, b, byLen), []string{"c", "go", "zig", "rust", "java"}; !slices.Equal(got, want) {
		t.Errorf("Merge(%q, %q, byLen) = %q, want %q", a, b, got, want)
	}

	fold := func(x, y string) int { return strings.Compare(strings.ToLower(x), strings.ToLower(y)) }
	a, b = []string{"apple", "Cherry"}, []string{"Banana", "date"}
	if got, want := mg.Merge(a, b, fold), []string{"apple", "Banana", "Cherry", "date"}; !slices.Equal(got, want) {
		t.Errorf("Merge(%q, %q, caseInsensitive) = %q, want %q", a, b, got, want)
	}
}

func TestMergeIsLinear(t *testing.T) {
	a := []int{1, 3, 5, 7, 9, 11, 13, 15}
	b := []int{2, 4, 6, 8, 10, 12, 14, 16}
	calls := 0
	counting := func(x, y int) int { calls++; return cmp.Compare(x, y) }
	mg.Merge(a, b, counting)
	if limit := len(a) + len(b) - 1; calls > limit {
		t.Errorf("Merge(%v, %v, compare) called compare %d times, want at most %d (one pass, no sorting)", a, b, calls, limit)
	}
}

func TestMergeNewSlice(t *testing.T) {
	a := []int{1, 2, 3}
	got := mg.Merge(a, nil, cmp.Compare[int])
	if len(got) != 3 {
		t.Fatalf("Merge(%v, nil, cmp.Compare) = %v, want %v", a, got, a)
	}
	got[0] = 99
	if a[0] != 1 {
		t.Errorf("Merge(%v, nil, cmp.Compare) shares memory with a: changing the result changed a to %v", []int{1, 2, 3}, a)
	}

	x, y := []int{5, 6}, []int{1, 2}
	mg.Merge(x, y, cmp.Compare[int])
	if !slices.Equal(x, []int{5, 6}) || !slices.Equal(y, []int{1, 2}) {
		t.Errorf("Merge modified its inputs: a = %v, b = %v, want [5 6], [1 2]", x, y)
	}
}

func TestDesc(t *testing.T) {
	desc := mg.Desc(cmp.Compare[int])
	checks := []struct {
		x, y int
		sign int
	}{
		{1, 2, 1}, {2, 1, -1}, {3, 3, 0},
	}
	for _, c := range checks {
		got := desc(c.x, c.y)
		if cmp.Compare(got, 0) != c.sign {
			t.Errorf("Desc(cmp.Compare)(%d, %d) = %d, want sign %d", c.x, c.y, got, c.sign)
		}
	}

	a, b := []int{9, 4}, []int{5, 1}
	if got, want := mg.Merge(a, b, mg.Desc(cmp.Compare[int])), []int{9, 5, 4, 1}; !slices.Equal(got, want) {
		t.Errorf("Merge(%v, %v, Desc(cmp.Compare)) = %v, want %v", a, b, got, want)
	}

	old := []Person{{"eve", 50}, {"ann", 30}}
	young := []Person{{"fay", 50}, {"bob", 20}}
	got := mg.Merge(old, young, mg.Desc(byAge))
	want := []Person{{"eve", 50}, {"fay", 50}, {"ann", 30}, {"bob", 20}}
	if !slices.Equal(got, want) {
		t.Errorf("Merge(%v, %v, Desc(byAge)) = %v, want %v", old, young, got, want)
	}
}

func TestMergeAll(t *testing.T) {
	got := mg.MergeAll(cmp.Compare[int], []int{1, 7}, []int{3}, []int{2, 8})
	if want := []int{1, 2, 3, 7, 8}; !slices.Equal(got, want) {
		t.Errorf("MergeAll(cmp.Compare, [1 7], [3], [2 8]) = %v, want %v", got, want)
	}

	got = mg.MergeAll(cmp.Compare[int], []int{5})
	if want := []int{5}; !slices.Equal(got, want) {
		t.Errorf("MergeAll(cmp.Compare, [5]) = %v, want %v", got, want)
	}

	got = mg.MergeAll(cmp.Compare[int], nil, []int{}, []int{4, 4})
	if want := []int{4, 4}; !slices.Equal(got, want) {
		t.Errorf("MergeAll(cmp.Compare, nil, [], [4 4]) = %v, want %v", got, want)
	}

	if got := mg.MergeAll(cmp.Compare[int]); len(got) != 0 {
		t.Errorf("MergeAll(cmp.Compare) = %v, want []", got)
	}

	l1 := []Person{{"a1", 1}, {"a2", 2}}
	l2 := []Person{{"b1", 1}, {"b3", 3}}
	l3 := []Person{{"c1", 1}, {"c2", 2}}
	people := mg.MergeAll(byAge, l1, l2, l3)
	wantP := []Person{{"a1", 1}, {"b1", 1}, {"c1", 1}, {"a2", 2}, {"c2", 2}, {"b3", 3}}
	if !slices.Equal(people, wantP) {
		t.Errorf("MergeAll(byAge, %v, %v, %v) = %v, want %v", l1, l2, l3, people, wantP)
	}
}
