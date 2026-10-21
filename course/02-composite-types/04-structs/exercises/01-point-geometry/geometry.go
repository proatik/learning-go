package geometry

// Point is a position in the 2-D plane.
type Point struct {
	X, Y float64
}

// Distance returns the straight-line distance between p and q.
func Distance(p, q Point) float64 {
	// TODO
	return 0
}

// Midpoint returns the point halfway between p and q.
func Midpoint(p, q Point) Point {
	// TODO
	return Point{}
}

// Translate returns p moved by dx horizontally and dy vertically.
func Translate(p Point, dx, dy float64) Point {
	// TODO
	return Point{}
}

// Centroid returns the average of points and true,
// or the zero Point and false if points is empty.
func Centroid(points []Point) (Point, bool) {
	// TODO
	return Point{}, false
}
