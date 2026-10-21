package shapes_test

import (
	"fmt"
	"math"
	"sort"
	"testing"

	shapes "learngo/course/03-methods-interfaces-errors/06-module-challenge/exercises/01-shape-catalogue"
)

// square is a Shape the package knows nothing about: Summary counts it as "other".
type square struct{ side float64 }

func (s square) Area() float64      { return s.side * s.side }
func (s square) Perimeter() float64 { return 4 * s.side }

func near(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestAreaPerimeter(t *testing.T) {
	tests := []struct {
		s           shapes.Shape
		area, perim float64
	}{
		{shapes.Circle{R: 1}, math.Pi, 2 * math.Pi},
		{shapes.Circle{R: 2.5}, math.Pi * 6.25, 5 * math.Pi},
		{shapes.Circle{R: 0}, 0, 0},
		{shapes.Rect{W: 2, H: 3}, 6, 10},
		{shapes.Rect{W: 0.5, H: 4}, 2, 9},
		{shapes.Rect{}, 0, 0},
		{shapes.Triangle{A: 3, B: 4, C: 5}, 6, 12},
		{shapes.Triangle{A: 2, B: 2, C: 2}, math.Sqrt(3), 6},
		{shapes.Triangle{A: 1, B: 2, C: 3}, 0, 0},    // degenerate: a+b == c
		{shapes.Triangle{A: 1, B: 1, C: 10}, 0, 0},   // impossible
		{shapes.Triangle{A: -3, B: 4, C: 5}, 0, 0},   // negative side
		{shapes.Triangle{A: 0, B: 4, C: 4}, 0, 0},    // zero side
		{&shapes.Circle{R: 1}, math.Pi, 2 * math.Pi}, // *Circle is a Shape too
	}
	for _, tt := range tests {
		if got := tt.s.Area(); !near(got, tt.area) {
			t.Errorf("%#v.Area() = %v, want %v", tt.s, got, tt.area)
		}
		if got := tt.s.Perimeter(); !near(got, tt.perim) {
			t.Errorf("%#v.Perimeter() = %v, want %v", tt.s, got, tt.perim)
		}
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		s    shapes.Shape
		want string
	}{
		{shapes.Circle{R: 1.5}, "circle(r=1.5)"},
		{shapes.Circle{R: 2}, "circle(r=2)"},
		{shapes.Rect{W: 2, H: 3}, "rect(2x3)"},
		{shapes.Rect{W: 0.25, H: 10}, "rect(0.25x10)"},
		{shapes.Triangle{A: 3, B: 4, C: 5}, "triangle(3,4,5)"},
		{shapes.Triangle{A: 1.5, B: 2, C: 2.5}, "triangle(1.5,2,2.5)"},
		{&shapes.Rect{W: 1, H: 1}, "rect(1x1)"}, // *Rect has String in its method set too
	}
	for _, tt := range tests {
		// fmt uses the String method because the shapes implement fmt.Stringer.
		if got := fmt.Sprint(tt.s); got != tt.want {
			t.Errorf("fmt.Sprint(%#v) = %q, want %q", tt.s, got, tt.want)
		}
	}
	if got := fmt.Sprintf("%v and %s", shapes.Circle{R: 1}, shapes.Rect{W: 4, H: 5}); got != "circle(r=1) and rect(4x5)" {
		t.Errorf("Sprintf(%%v and %%s) = %q, want %q", got, "circle(r=1) and rect(4x5)")
	}
}

func TestByAreaIsSortInterface(t *testing.T) {
	s := shapes.ByArea{shapes.Rect{W: 3, H: 3}, shapes.Circle{R: 1}, shapes.Rect{W: 1, H: 1}}
	if s.Len() != 3 {
		t.Fatalf("ByArea.Len() = %d, want 3", s.Len())
	}
	if !s.Less(1, 0) || s.Less(0, 1) {
		t.Errorf("ByArea.Less: circle(r=1) should be less than rect(3x3)")
	}
	s.Swap(0, 2)
	if got := fmt.Sprint(s[0], " ", s[2]); got != "rect(1x1) rect(3x3)" {
		t.Errorf("after Swap(0, 2): [0], [2] = %q, want %q", got, "rect(1x1) rect(3x3)")
	}
	sort.Sort(s)
	if got := fmt.Sprint([]shapes.Shape(s)); got != "[rect(1x1) circle(r=1) rect(3x3)]" {
		t.Errorf("sort.Sort(ByArea) = %s, want %s", got, "[rect(1x1) circle(r=1) rect(3x3)]")
	}
}

func TestSortByArea(t *testing.T) {
	tests := []struct {
		in   []shapes.Shape
		want string
	}{
		{nil, "[]"},
		{[]shapes.Shape{}, "[]"},
		{[]shapes.Shape{shapes.Circle{R: 1}}, "[circle(r=1)]"},
		{
			[]shapes.Shape{shapes.Rect{W: 10, H: 10}, shapes.Triangle{A: 3, B: 4, C: 5}, shapes.Circle{R: 1}, shapes.Rect{W: 1, H: 2}},
			"[rect(1x2) circle(r=1) triangle(3,4,5) rect(10x10)]",
		},
		{ // equal areas keep their original order (stable)
			[]shapes.Shape{shapes.Rect{W: 2, H: 3}, shapes.Rect{W: 1, H: 1}, shapes.Rect{W: 3, H: 2}, shapes.Triangle{A: 3, B: 4, C: 5}, shapes.Rect{W: 6, H: 1}},
			"[rect(1x1) rect(2x3) rect(3x2) triangle(3,4,5) rect(6x1)]",
		},
		{
			[]shapes.Shape{square{3}, shapes.Rect{W: 2, H: 2}, shapes.Triangle{A: 1, B: 1, C: 5}},
			"[triangle(1,1,5) rect(2x2) {3}]",
		},
	}
	for _, tt := range tests {
		in := fmt.Sprint(tt.in)
		shapes.SortByArea(tt.in)
		if got := fmt.Sprint(tt.in); got != tt.want {
			t.Errorf("SortByArea(%s) = %s, want %s", in, got, tt.want)
		}
	}
}

func TestLargest(t *testing.T) {
	tests := []struct {
		in     []shapes.Shape
		want   string
		wantOK bool
	}{
		{nil, "<nil>", false},
		{[]shapes.Shape{}, "<nil>", false},
		{[]shapes.Shape{nil, nil}, "<nil>", false},
		{[]shapes.Shape{shapes.Rect{W: 1, H: 1}}, "rect(1x1)", true},
		{[]shapes.Shape{shapes.Rect{W: 1, H: 1}, shapes.Circle{R: 2}, shapes.Rect{W: 3, H: 3}}, "circle(r=2)", true},
		{[]shapes.Shape{nil, shapes.Rect{W: 2, H: 3}, nil, shapes.Rect{W: 3, H: 2}}, "rect(2x3)", true}, // first of equal areas
		{[]shapes.Shape{shapes.Triangle{A: 1, B: 1, C: 9}}, "triangle(1,1,9)", true},                    // area 0 is still a shape
		{[]shapes.Shape{shapes.Rect{W: 1, H: 1}, square{5}}, "{5}", true},
	}
	for _, tt := range tests {
		got, ok := shapes.Largest(tt.in)
		if fmt.Sprint(got) != tt.want || ok != tt.wantOK {
			t.Errorf("Largest(%v) = %v, %v, want %s, %v", tt.in, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestSummary(t *testing.T) {
	tests := []struct {
		in   []shapes.Shape
		want string
	}{
		{nil, "circles=0 rects=0 triangles=0 other=0 area=0.00"},
		{[]shapes.Shape{nil}, "circles=0 rects=0 triangles=0 other=0 area=0.00"},
		{[]shapes.Shape{shapes.Circle{R: 1}}, "circles=1 rects=0 triangles=0 other=0 area=3.14"},
		{
			[]shapes.Shape{shapes.Circle{R: 1}, shapes.Rect{W: 2, H: 3}, shapes.Rect{W: 1, H: 1}, shapes.Triangle{A: 3, B: 4, C: 5}, square{2}},
			"circles=1 rects=2 triangles=1 other=1 area=20.14",
		},
		{ // pointers count as their kind too
			[]shapes.Shape{&shapes.Circle{R: 2}, &shapes.Rect{W: 1, H: 2}, &shapes.Triangle{A: 3, B: 4, C: 5}, nil},
			"circles=1 rects=1 triangles=1 other=0 area=20.57",
		},
		{[]shapes.Shape{square{1}, square{2}, shapes.Triangle{A: 1, B: 2, C: 3}}, "circles=0 rects=0 triangles=1 other=2 area=5.00"},
	}
	for _, tt := range tests {
		if got := shapes.Summary(tt.in); got != tt.want {
			t.Errorf("Summary(%v) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
