package finder_test

import (
	"testing"

	fnd "learngo/course/04-generics-packages-testing/01-generic-functions/exercises/02-index-of"
)

type Point struct{ X, Y int }

type Status string

func TestIndexOfStrings(t *testing.T) {
	tests := []struct {
		in             []string
		target         string
		first, last, n int
	}{
		{[]string{"a", "b", "a"}, "a", 0, 2, 2},
		{[]string{"a", "b", "a"}, "b", 1, 1, 1},
		{[]string{"a", "b", "a"}, "z", -1, -1, 0},
		{[]string{"", "x", ""}, "", 0, 2, 2},
		{[]string{"héllo", "日本"}, "日本", 1, 1, 1},
		{[]string{"Go"}, "go", -1, -1, 0},
		{[]string{}, "a", -1, -1, 0},
		{nil, "", -1, -1, 0},
	}
	for _, tt := range tests {
		if got := fnd.IndexOf(tt.in, tt.target); got != tt.first {
			t.Errorf("IndexOf(%q, %q) = %d, want %d", tt.in, tt.target, got, tt.first)
		}
		if got := fnd.LastIndexOf(tt.in, tt.target); got != tt.last {
			t.Errorf("LastIndexOf(%q, %q) = %d, want %d", tt.in, tt.target, got, tt.last)
		}
		if got := fnd.Count(tt.in, tt.target); got != tt.n {
			t.Errorf("Count(%q, %q) = %d, want %d", tt.in, tt.target, got, tt.n)
		}
		want := tt.first >= 0
		if got := fnd.Contains(tt.in, tt.target); got != want {
			t.Errorf("Contains(%q, %q) = %t, want %t", tt.in, tt.target, got, want)
		}
	}
}

func TestIndexOfInts(t *testing.T) {
	tests := []struct {
		in             []int
		target         int
		first, last, n int
	}{
		{[]int{1, 2, 3}, 9, -1, -1, 0},
		{[]int{1, 2, 3}, 3, 2, 2, 1},
		{[]int{1, 2, 1, 1}, 1, 0, 3, 3},
		{[]int{0, 0}, 0, 0, 1, 2},
		{[]int{-4, 4}, -4, 0, 0, 1},
		{nil, 0, -1, -1, 0},
	}
	for _, tt := range tests {
		if got := fnd.IndexOf(tt.in, tt.target); got != tt.first {
			t.Errorf("IndexOf(%v, %d) = %d, want %d", tt.in, tt.target, got, tt.first)
		}
		if got := fnd.LastIndexOf(tt.in, tt.target); got != tt.last {
			t.Errorf("LastIndexOf(%v, %d) = %d, want %d", tt.in, tt.target, got, tt.last)
		}
		if got := fnd.Count(tt.in, tt.target); got != tt.n {
			t.Errorf("Count(%v, %d) = %d, want %d", tt.in, tt.target, got, tt.n)
		}
		want := tt.first >= 0
		if got := fnd.Contains(tt.in, tt.target); got != want {
			t.Errorf("Contains(%v, %d) = %t, want %t", tt.in, tt.target, got, want)
		}
	}
}

func TestStructsAndNamedTypes(t *testing.T) {
	pts := []Point{{1, 2}, {3, 4}, {1, 2}}
	if got := fnd.IndexOf(pts, Point{3, 4}); got != 1 {
		t.Errorf("IndexOf(%v, %v) = %d, want %d", pts, Point{3, 4}, got, 1)
	}
	if got := fnd.Count(pts, Point{1, 2}); got != 2 {
		t.Errorf("Count(%v, %v) = %d, want %d", pts, Point{1, 2}, got, 2)
	}
	if got := fnd.LastIndexOf(pts, Point{1, 2}); got != 2 {
		t.Errorf("LastIndexOf(%v, %v) = %d, want %d", pts, Point{1, 2}, got, 2)
	}
	if got := fnd.Contains(pts, Point{2, 1}); got {
		t.Errorf("Contains(%v, %v) = %t, want %t", pts, Point{2, 1}, got, false)
	}

	statuses := []Status{"todo", "done", "todo"}
	if got := fnd.Count(statuses, "todo"); got != 2 {
		t.Errorf("Count(%q, %q) = %d, want %d", statuses, "todo", got, 2)
	}
	if got := fnd.Contains(statuses, "done"); !got {
		t.Errorf("Contains(%q, %q) = %t, want %t", statuses, "done", got, true)
	}

	a, b := new(int), new(int)
	ptrs := []*int{a, b}
	if got := fnd.IndexOf(ptrs, b); got != 1 {
		t.Errorf("IndexOf(ptrs, b) = %d, want %d (pointers compare by address)", got, 1)
	}
	if got := fnd.Contains(ptrs, new(int)); got {
		t.Errorf("Contains(ptrs, new(int)) = %t, want %t (a different address)", got, false)
	}
}
