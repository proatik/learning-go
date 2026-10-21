package geom

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

// Area returns W*H.
func (r Rect) Area() float64 {
	// TODO
	return 0
}

// Perimeter returns 2*(W+H).
func (r Rect) Perimeter() float64 {
	// TODO
	return 0
}

// Area returns πR².
func (c Circle) Area() float64 {
	// TODO
	return 0
}

// Perimeter returns 2πR.
func (c Circle) Perimeter() float64 {
	// TODO
	return 0
}

// Area returns the area computed with Heron's formula.
func (t Triangle) Area() float64 {
	// TODO
	return 0
}

// Perimeter returns A+B+C.
func (t Triangle) Perimeter() float64 {
	// TODO
	return 0
}

// TotalArea returns the sum of the areas of all shapes.
func TotalArea(shapes []Shape) float64 {
	// TODO
	return 0
}

// Largest returns the shape with the biggest area (the first on a tie)
// and true, or nil and false when shapes is empty.
func Largest(shapes []Shape) (Shape, bool) {
	// TODO
	return nil, false
}
