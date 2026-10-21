package geom

import "math"

// Shape is anything with an area and a perimeter.
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rect is an axis-aligned rectangle of width W and height H.
type Rect struct{ W, H float64 }

// Circle is a circle of radius R.
type Circle struct{ R float64 }

// Triangle is a triangle given by its three side lengths.
type Triangle struct{ A, B, C float64 }

// Compile-time checks that every type satisfies Shape.
var (
	_ Shape = Rect{}
	_ Shape = Circle{}
	_ Shape = Triangle{}
)

// Area returns W*H.
func (r Rect) Area() float64 { return r.W * r.H }

// Perimeter returns 2*(W+H).
func (r Rect) Perimeter() float64 { return 2 * (r.W + r.H) }

// Area returns πR².
func (c Circle) Area() float64 { return math.Pi * c.R * c.R }

// Perimeter returns 2πR.
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.R }

// Area returns the area computed with Heron's formula.
func (t Triangle) Area() float64 {
	s := t.Perimeter() / 2
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}

// Perimeter returns A+B+C.
func (t Triangle) Perimeter() float64 { return t.A + t.B + t.C }

// TotalArea returns the sum of the areas of all shapes.
func TotalArea(shapes []Shape) float64 {
	total := 0.0
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

// Largest returns the shape with the biggest area (the first on a tie)
// and true, or nil and false when shapes is empty.
func Largest(shapes []Shape) (Shape, bool) {
	if len(shapes) == 0 {
		return nil, false
	}
	best := shapes[0]
	for _, s := range shapes[1:] {
		if s.Area() > best.Area() {
			best = s
		}
	}
	return best, true
}
