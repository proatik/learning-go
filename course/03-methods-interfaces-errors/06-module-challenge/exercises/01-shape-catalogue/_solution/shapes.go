package shapes

import (
	"fmt"
	"math"
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
func (c Circle) Area() float64 { return math.Pi * c.R * c.R }

// Perimeter returns 2πr.
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.R }

// String returns "circle(r=R)".
func (c Circle) String() string { return fmt.Sprintf("circle(r=%g)", c.R) }

// Area returns W×H.
func (r Rect) Area() float64 { return r.W * r.H }

// Perimeter returns 2(W+H).
func (r Rect) Perimeter() float64 { return 2 * (r.W + r.H) }

// String returns "rect(WxH)".
func (r Rect) String() string { return fmt.Sprintf("rect(%gx%g)", r.W, r.H) }

// valid reports whether the sides form a real (non-degenerate) triangle.
func (t Triangle) valid() bool {
	return t.A > 0 && t.B > 0 && t.C > 0 &&
		t.A+t.B > t.C && t.A+t.C > t.B && t.B+t.C > t.A
}

// Area uses Heron's formula; it returns 0 for an invalid triangle.
func (t Triangle) Area() float64 {
	if !t.valid() {
		return 0
	}
	s := (t.A + t.B + t.C) / 2
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}

// Perimeter returns A+B+C; it returns 0 for an invalid triangle.
func (t Triangle) Perimeter() float64 {
	if !t.valid() {
		return 0
	}
	return t.A + t.B + t.C
}

// String returns "triangle(A,B,C)".
func (t Triangle) String() string { return fmt.Sprintf("triangle(%g,%g,%g)", t.A, t.B, t.C) }

// ByArea sorts shapes by ascending area.
type ByArea []Shape

// Len returns the number of shapes.
func (s ByArea) Len() int { return len(s) }

// Less reports whether shape i has a smaller area than shape j.
func (s ByArea) Less(i, j int) bool { return s[i].Area() < s[j].Area() }

// Swap exchanges shapes i and j.
func (s ByArea) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// SortByArea sorts shapes in place by ascending area, keeping the original
// order of shapes with equal areas.
func SortByArea(shapes []Shape) {
	sort.Stable(ByArea(shapes))
}

// Largest returns the first shape with the greatest area, skipping nil
// entries. It returns nil, false if there is no shape.
func Largest(shapes []Shape) (Shape, bool) {
	var best Shape
	for _, s := range shapes {
		if s == nil {
			continue
		}
		if best == nil || s.Area() > best.Area() {
			best = s
		}
	}
	return best, best != nil
}

// Summary counts the shapes by kind and adds up their areas.
func Summary(shapes []Shape) string {
	var circles, rects, triangles, other int
	total := 0.0
	for _, s := range shapes {
		switch s.(type) {
		case nil:
			continue
		case Circle, *Circle:
			circles++
		case Rect, *Rect:
			rects++
		case Triangle, *Triangle:
			triangles++
		default:
			other++
		}
		total += s.Area()
	}
	return fmt.Sprintf("circles=%d rects=%d triangles=%d other=%d area=%.2f",
		circles, rects, triangles, other, total)
}
