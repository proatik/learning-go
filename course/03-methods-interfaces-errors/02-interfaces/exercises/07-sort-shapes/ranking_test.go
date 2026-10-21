package ranking_test

import (
	"math"
	"slices"
	"testing"

	ranking "learngo/course/03-methods-interfaces-errors/02-interfaces/exercises/07-sort-shapes"
)

type (
	sq     = ranking.Square
	rect   = ranking.Rect
	circle = ranking.Circle
)

// hexagon is a Shape type the package doesn't know about.
type hexagon struct{ S float64 }

func (h hexagon) Area() float64 { return 3 * math.Sqrt(3) / 2 * h.S * h.S }
func (h hexagon) Name() string  { return "hexagon" }

func TestMethods(t *testing.T) {
	tests := []struct {
		s    ranking.Shape
		area float64
		name string
	}{
		{sq{Side: 3}, 9, "square"},
		{sq{}, 0, "square"},
		{rect{W: 2, H: 3}, 6, "rect"},
		{circle{R: 1}, math.Pi, "circle"},
		{circle{R: 2}, 4 * math.Pi, "circle"},
	}
	for _, tt := range tests {
		if got := tt.s.Area(); math.Abs(got-tt.area) > 1e-9 {
			t.Errorf("%+v.Area() = %v, want %v", tt.s, got, tt.area)
		}
		if got := tt.s.Name(); got != tt.name {
			t.Errorf("%+v.Name() = %q, want %q", tt.s, got, tt.name)
		}
	}
}

func TestSortByArea(t *testing.T) {
	tests := []struct {
		in, want []ranking.Shape
	}{
		{nil, nil},
		{[]ranking.Shape{}, nil},
		{[]ranking.Shape{sq{Side: 1}}, []ranking.Shape{sq{Side: 1}}},
		{
			[]ranking.Shape{rect{W: 2, H: 3}, circle{R: 1}, sq{Side: 1}},
			[]ranking.Shape{sq{Side: 1}, circle{R: 1}, rect{W: 2, H: 3}},
		},
		{
			[]ranking.Shape{sq{Side: 2}, rect{W: 1, H: 4}},
			[]ranking.Shape{rect{W: 1, H: 4}, sq{Side: 2}},
		},
		{
			[]ranking.Shape{hexagon{S: 1}, sq{Side: 10}, rect{W: 0, H: 9}, circle{R: 1}},
			[]ranking.Shape{rect{W: 0, H: 9}, hexagon{S: 1}, circle{R: 1}, sq{Side: 10}},
		},
		{
			[]ranking.Shape{rect{W: 4, H: 1}, rect{W: 1, H: 4}, rect{W: 2, H: 2}},
			[]ranking.Shape{rect{W: 4, H: 1}, rect{W: 1, H: 4}, rect{W: 2, H: 2}}, // full ties keep input order
		},
	}
	for _, tt := range tests {
		in := slices.Clone(tt.in)
		got := ranking.SortByArea(in)
		if !slices.Equal(got, tt.want) {
			t.Errorf("SortByArea(%v) = %v, want %v", tt.in, got, tt.want)
		}
		if !slices.Equal(in, tt.in) {
			t.Errorf("SortByArea(%v) modified its input to %v", tt.in, in)
		}
	}
}

func TestTopN(t *testing.T) {
	shapes := []ranking.Shape{sq{Side: 1}, sq{Side: 3}, circle{R: 1}, rect{W: 3, H: 3}, rect{W: 0.5, H: 1}}
	tests := []struct {
		in   []ranking.Shape
		n    int
		want []ranking.Shape
	}{
		{shapes, 2, []ranking.Shape{rect{W: 3, H: 3}, sq{Side: 3}}},
		{shapes, 1, []ranking.Shape{rect{W: 3, H: 3}}},
		{shapes, 3, []ranking.Shape{rect{W: 3, H: 3}, sq{Side: 3}, circle{R: 1}}},
		{shapes, 10, []ranking.Shape{rect{W: 3, H: 3}, sq{Side: 3}, circle{R: 1}, sq{Side: 1}, rect{W: 0.5, H: 1}}},
		{shapes, 0, nil},
		{shapes, -1, nil},
		{nil, 3, nil},
		{[]ranking.Shape{sq{Side: 1}, sq{Side: 3}, circle{R: 1}}, 2, []ranking.Shape{sq{Side: 3}, circle{R: 1}}},
	}
	for _, tt := range tests {
		in := slices.Clone(tt.in)
		got := ranking.TopN(in, tt.n)
		if !slices.Equal(got, tt.want) {
			t.Errorf("TopN(%v, %d) = %v, want %v", tt.in, tt.n, got, tt.want)
		}
		if !slices.Equal(in, tt.in) {
			t.Errorf("TopN(%v, %d) modified its input to %v", tt.in, tt.n, in)
		}
	}
}

func TestKind(t *testing.T) {
	tests := []struct {
		s    ranking.Shape
		want string
	}{
		{sq{Side: 1}, "polygon"},
		{rect{W: 1, H: 2}, "polygon"},
		{circle{R: 1}, "round"},
		{nil, "none"},
		{hexagon{S: 1}, "unknown"},
	}
	for _, tt := range tests {
		if got := ranking.Kind(tt.s); got != tt.want {
			t.Errorf("Kind(%#v) = %q, want %q", tt.s, got, tt.want)
		}
	}
}
