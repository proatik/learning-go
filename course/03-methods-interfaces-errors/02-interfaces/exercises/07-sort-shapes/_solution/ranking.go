package ranking

import (
	"cmp"
	"math"
	"slices"
)

// Shape is a named shape with an area.
type Shape interface {
	Area() float64
	Name() string
}

// Square is a square with side length Side.
type Square struct{ Side float64 }

// Rect is a rectangle of width W and height H.
type Rect struct{ W, H float64 }

// Circle is a circle of radius R.
type Circle struct{ R float64 }

var (
	_ Shape = Square{}
	_ Shape = Rect{}
	_ Shape = Circle{}
)

// Area returns Side².
func (s Square) Area() float64 { return s.Side * s.Side }

// Name returns "square".
func (s Square) Name() string { return "square" }

// Area returns W*H.
func (r Rect) Area() float64 { return r.W * r.H }

// Name returns "rect".
func (r Rect) Name() string { return "rect" }

// Area returns πR².
func (c Circle) Area() float64 { return math.Pi * c.R * c.R }

// Name returns "circle".
func (c Circle) Name() string { return "circle" }

// byArea orders shapes by area ascending, then by name.
func byArea(a, b Shape) int {
	return cmp.Or(
		cmp.Compare(a.Area(), b.Area()),
		cmp.Compare(a.Name(), b.Name()),
	)
}

// SortByArea returns a sorted copy of shapes: area ascending, then name.
func SortByArea(shapes []Shape) []Shape {
	out := slices.Clone(shapes)
	slices.SortStableFunc(out, byArea)
	return out
}

// TopN returns the n largest shapes, largest first (ties by name).
func TopN(shapes []Shape, n int) []Shape {
	if n <= 0 {
		return nil
	}
	out := slices.Clone(shapes)
	slices.SortStableFunc(out, func(a, b Shape) int {
		return cmp.Or(
			cmp.Compare(b.Area(), a.Area()), // descending
			cmp.Compare(a.Name(), b.Name()),
		)
	})
	return out[:min(n, len(out))]
}

// Kind classifies s as "polygon", "round", "none" or "unknown".
func Kind(s Shape) string {
	switch s.(type) {
	case nil:
		return "none"
	case Square, Rect:
		return "polygon"
	case Circle:
		return "round"
	default:
		return "unknown"
	}
}
