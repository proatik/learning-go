package shapes

import (
	"fmt"
	"sort"
)

// Shape is anything with an area and a perimeter.
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Circle is a circle with radius R.
type Circle struct{ R float64 }

// Rect is an axis-aligned rectangle of width W and height H.
type Rect struct{ W, H float64 }

// Triangle is a triangle given by its three side lengths.
type Triangle struct{ A, B, C float64 }

// Compile-time checks: the value types satisfy Shape.
var (
	_ Shape          = Circle{}
	_ Shape          = Rect{}
	_ Shape          = Triangle{}
	_ fmt.Stringer   = Circle{}
	_ sort.Interface = ByArea(nil)
)

// Area returns πr².
func (c Circle) Area() float64 {
	// TODO
	return 0
}

// Perimeter returns 2πr.
func (c Circle) Perimeter() float64 {
	// TODO
	return 0
}

// String returns "circle(r=R)".
func (c Circle) String() string {
	// TODO
	return ""
}

// Area returns W×H.
func (r Rect) Area() float64 {
	// TODO
	return 0
}

// Perimeter returns 2(W+H).
func (r Rect) Perimeter() float64 {
	// TODO
	return 0
}

// String returns "rect(WxH)".
func (r Rect) String() string {
	// TODO
	return ""
}

// Area uses Heron's formula; it returns 0 for an invalid triangle.
func (t Triangle) Area() float64 {
	// TODO
	return 0
}

// Perimeter returns A+B+C; it returns 0 for an invalid triangle.
func (t Triangle) Perimeter() float64 {
	// TODO
	return 0
}

// String returns "triangle(A,B,C)".
func (t Triangle) String() string {
	// TODO
	return ""
}

// ByArea sorts shapes by ascending area.
type ByArea []Shape

// Len returns the number of shapes.
func (s ByArea) Len() int {
	// TODO
	return 0
}

// Less reports whether shape i has a smaller area than shape j.
func (s ByArea) Less(i, j int) bool {
	// TODO
	return false
}

// Swap exchanges shapes i and j.
func (s ByArea) Swap(i, j int) {
	// TODO
}

// SortByArea sorts shapes in place by ascending area, keeping the original
// order of shapes with equal areas.
func SortByArea(shapes []Shape) {
	// TODO
}

// Largest returns the first shape with the greatest area, skipping nil
// entries. It returns nil, false if there is no shape.
func Largest(shapes []Shape) (Shape, bool) {
	// TODO
	return nil, false
}

// Summary counts the shapes by kind and adds up their areas.
func Summary(shapes []Shape) string {
	// TODO
	return ""
}
