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
	// TODO
	return 0
}

// Contains reports whether p is inside r.
func Contains(r Rect, p Point) bool {
	// TODO
	return false
}

// Intersect returns the overlap of a and b and true,
// or Rect{} and false if the overlap has zero area.
func Intersect(a, b Rect) (Rect, bool) {
	// TODO
	return Rect{}, false
}

// Bounds returns the smallest Rect containing every rect in rs,
// or Rect{} if rs is empty.
func Bounds(rs []Rect) Rect {
	// TODO
	return Rect{}
}
