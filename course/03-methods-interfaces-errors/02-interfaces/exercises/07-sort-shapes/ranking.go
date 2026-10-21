package ranking

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

// Area returns Side².
func (s Square) Area() float64 {
	// TODO
	return 0
}

// Name returns "square".
func (s Square) Name() string {
	// TODO
	return ""
}

// Area returns W*H.
func (r Rect) Area() float64 {
	// TODO
	return 0
}

// Name returns "rect".
func (r Rect) Name() string {
	// TODO
	return ""
}

// Area returns πR².
func (c Circle) Area() float64 {
	// TODO
	return 0
}

// Name returns "circle".
func (c Circle) Name() string {
	// TODO
	return ""
}

// SortByArea returns a sorted copy of shapes: area ascending, then name.
func SortByArea(shapes []Shape) []Shape {
	// TODO
	return nil
}

// TopN returns the n largest shapes, largest first (ties by name).
func TopN(shapes []Shape, n int) []Shape {
	// TODO
	return nil
}

// Kind classifies s as "polygon", "round", "none" or "unknown".
func Kind(s Shape) string {
	// TODO
	return ""
}
