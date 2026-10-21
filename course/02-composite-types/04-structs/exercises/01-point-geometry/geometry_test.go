package geometry

import (
	"math"
	"testing"
)

func near(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func nearPoint(a, b Point) bool { return near(a.X, b.X) && near(a.Y, b.Y) }

func TestDistance(t *testing.T) {
	tests := []struct {
		p, q Point
		want float64
	}{
		{Point{0, 0}, Point{3, 4}, 5},
		{Point{3, 4}, Point{0, 0}, 5},
		{Point{1, 1}, Point{1, 1}, 0},
		{Point{-1, -1}, Point{2, 3}, 5},
		{Point{0, 0}, Point{0, -7}, 7},
		{Point{0, 0}, Point{1, 1}, math.Sqrt2},
		{Point{1.5, 2.5}, Point{1.5, -2.5}, 5},
	}
	for _, tt := range tests {
		if got := Distance(tt.p, tt.q); !near(got, tt.want) {
			t.Errorf("Distance(%+v, %+v) = %v, want %v", tt.p, tt.q, got, tt.want)
		}
	}
}

func TestMidpoint(t *testing.T) {
	tests := []struct {
		p, q Point
		want Point
	}{
		{Point{0, 0}, Point{4, -2}, Point{X: 2, Y: -1}},
		{Point{1, 1}, Point{1, 1}, Point{X: 1, Y: 1}},
		{Point{-3, 5}, Point{3, -5}, Point{}},
		{Point{0, 0}, Point{1, 3}, Point{X: 0.5, Y: 1.5}},
	}
	for _, tt := range tests {
		if got := Midpoint(tt.p, tt.q); !nearPoint(got, tt.want) {
			t.Errorf("Midpoint(%+v, %+v) = %+v, want %+v", tt.p, tt.q, got, tt.want)
		}
	}
}

func TestTranslate(t *testing.T) {
	tests := []struct {
		p      Point
		dx, dy float64
		want   Point
	}{
		{Point{1, 2}, 10, -5, Point{X: 11, Y: -3}},
		{Point{}, 0, 0, Point{}},
		{Point{-1, -1}, 1, 1, Point{}},
		{Point{2.5, 0}, 0, 0.5, Point{X: 2.5, Y: 0.5}},
	}
	for _, tt := range tests {
		orig := tt.p
		if got := Translate(tt.p, tt.dx, tt.dy); !nearPoint(got, tt.want) {
			t.Errorf("Translate(%+v, %v, %v) = %+v, want %+v", tt.p, tt.dx, tt.dy, got, tt.want)
		}
		if tt.p != orig {
			t.Errorf("Translate changed its argument: %+v, want %+v", tt.p, orig)
		}
	}
}

func TestCentroid(t *testing.T) {
	tests := []struct {
		points []Point
		want   Point
		wantOK bool
	}{
		{[]Point{{0, 0}, {4, 0}, {2, 6}}, Point{X: 2, Y: 2}, true},
		{[]Point{{5, -5}}, Point{X: 5, Y: -5}, true},
		{[]Point{{0, 0}, {1, 0}, {1, 1}, {0, 1}}, Point{X: 0.5, Y: 0.5}, true},
		{[]Point{{-2, 3}, {2, -3}}, Point{}, true},
		{[]Point{}, Point{}, false},
		{nil, Point{}, false},
	}
	for _, tt := range tests {
		got, ok := Centroid(tt.points)
		if ok != tt.wantOK || !nearPoint(got, tt.want) {
			t.Errorf("Centroid(%+v) = %+v, %v; want %+v, %v", tt.points, got, ok, tt.want, tt.wantOK)
		}
	}
}
