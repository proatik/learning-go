package rect

// Point is a position on an integer grid.
type Point struct {
	X, Y int
}

// Rect contains the points p with Min.X <= p.X < Max.X and Min.Y <= p.Y < Max.Y.
type Rect struct {
	Min, Max Point
}

// Area returns the number of grid cells covered by r.
func Area(r Rect) int {
	return (r.Max.X - r.Min.X) * (r.Max.Y - r.Min.Y)
}

// Contains reports whether p is inside r.
func Contains(r Rect, p Point) bool {
	return r.Min.X <= p.X && p.X < r.Max.X &&
		r.Min.Y <= p.Y && p.Y < r.Max.Y
}

// Intersect returns the overlap of a and b and true,
// or Rect{} and false if the overlap has zero area.
func Intersect(a, b Rect) (Rect, bool) {
	r := Rect{
		Min: Point{X: max(a.Min.X, b.Min.X), Y: max(a.Min.Y, b.Min.Y)},
		Max: Point{X: min(a.Max.X, b.Max.X), Y: min(a.Max.Y, b.Max.Y)},
	}
	if r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y {
		return Rect{}, false
	}
	return r, true
}

// Bounds returns the smallest Rect containing every rect in rs,
// or Rect{} if rs is empty.
func Bounds(rs []Rect) Rect {
	if len(rs) == 0 {
		return Rect{}
	}
	b := rs[0] // a copy: updating b never changes rs[0]
	for _, r := range rs[1:] {
		b.Min.X = min(b.Min.X, r.Min.X)
		b.Min.Y = min(b.Min.Y, r.Min.Y)
		b.Max.X = max(b.Max.X, r.Max.X)
		b.Max.Y = max(b.Max.Y, r.Max.Y)
	}
	return b
}
