package geometry

import "math"

// Point is a position in the 2-D plane.
type Point struct {
	X, Y float64
}

// Distance returns the straight-line distance between p and q.
func Distance(p, q Point) float64 {
	return math.Hypot(q.X-p.X, q.Y-p.Y)
}

// Midpoint returns the point halfway between p and q.
func Midpoint(p, q Point) Point {
	return Point{X: (p.X + q.X) / 2, Y: (p.Y + q.Y) / 2}
}

// Translate returns p moved by dx horizontally and dy vertically.
func Translate(p Point, dx, dy float64) Point {
	// p is already a copy of the caller's value: modifying it is safe.
	p.X += dx
	p.Y += dy
	return p
}

// Centroid returns the average of points and true,
// or the zero Point and false if points is empty.
func Centroid(points []Point) (Point, bool) {
	if len(points) == 0 {
		return Point{}, false
	}
	var sum Point // the zero value is a ready-to-use accumulator
	for _, p := range points {
		sum.X += p.X
		sum.Y += p.Y
	}
	n := float64(len(points))
	return Point{X: sum.X / n, Y: sum.Y / n}, true
}
