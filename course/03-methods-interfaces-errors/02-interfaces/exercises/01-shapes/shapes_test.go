package geom_test

import (
	"math"
	"testing"

	geom "learngo/course/03-methods-interfaces-errors/02-interfaces/exercises/01-shapes"
)

func near(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestAreaAndPerimeter(t *testing.T) {
	tests := []struct {
		s          geom.Shape
		area, peri float64
	}{
		{geom.Rect{W: 3, H: 4}, 12, 14},
		{geom.Rect{W: 0, H: 5}, 0, 10},
		{geom.Rect{W: 2.5, H: 2}, 5, 9},
		{geom.Circle{R: 1}, math.Pi, 2 * math.Pi},
		{geom.Circle{R: 2}, 4 * math.Pi, 4 * math.Pi},
		{geom.Circle{R: 0}, 0, 0},
		{geom.Triangle{A: 3, B: 4, C: 5}, 6, 12},
		{geom.Triangle{A: 2, B: 2, C: 2}, math.Sqrt(3), 6},
		{geom.Triangle{A: 5, B: 5, C: 6}, 12, 16},
	}
	for _, tt := range tests {
		if got := tt.s.Area(); !near(got, tt.area) {
			t.Errorf("%+v.Area() = %v, want %v", tt.s, got, tt.area)
		}
		if got := tt.s.Perimeter(); !near(got, tt.peri) {
			t.Errorf("%+v.Perimeter() = %v, want %v", tt.s, got, tt.peri)
		}
	}
}

func TestTotalArea(t *testing.T) {
	tests := []struct {
		in   []geom.Shape
		want float64
	}{
		{nil, 0},
		{[]geom.Shape{}, 0},
		{[]geom.Shape{geom.Rect{W: 3, H: 4}}, 12},
		{[]geom.Shape{geom.Rect{W: 3, H: 4}, geom.Triangle{A: 3, B: 4, C: 5}}, 18},
		{[]geom.Shape{geom.Circle{R: 1}, geom.Circle{R: 1}, geom.Rect{W: 1, H: 1}}, 2*math.Pi + 1},
	}
	for _, tt := range tests {
		if got := geom.TotalArea(tt.in); !near(got, tt.want) {
			t.Errorf("TotalArea(%+v) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestLargest(t *testing.T) {
	tests := []struct {
		in     []geom.Shape
		want   geom.Shape
		wantOK bool
	}{
		{nil, nil, false},
		{[]geom.Shape{}, nil, false},
		{[]geom.Shape{geom.Circle{R: 1}}, geom.Circle{R: 1}, true},
		{[]geom.Shape{geom.Rect{W: 1, H: 1}, geom.Circle{R: 1}, geom.Rect{W: 1, H: 4}}, geom.Rect{W: 1, H: 4}, true},
		{[]geom.Shape{geom.Triangle{A: 3, B: 4, C: 5}, geom.Rect{W: 2, H: 2}}, geom.Triangle{A: 3, B: 4, C: 5}, true},
		{[]geom.Shape{geom.Rect{W: 2, H: 3}, geom.Rect{W: 3, H: 2}, geom.Rect{W: 6, H: 1}}, geom.Rect{W: 2, H: 3}, true},
		{[]geom.Shape{geom.Rect{W: 0, H: 0}, geom.Circle{R: 0}}, geom.Rect{W: 0, H: 0}, true},
	}
	for _, tt := range tests {
		got, ok := geom.Largest(tt.in)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("Largest(%+v) = %+v, %v, want %+v, %v", tt.in, got, ok, tt.want, tt.wantOK)
		}
	}
}
