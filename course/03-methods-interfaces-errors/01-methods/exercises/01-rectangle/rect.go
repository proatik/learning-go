package rect

// Rect is an axis-aligned rectangle with width W and height H.
type Rect struct {
	W, H float64
}

// Area returns W * H.
func (r Rect) Area() float64 {
	// TODO
	return 0
}

// Perimeter returns 2 * (W + H).
func (r Rect) Perimeter() float64 {
	// TODO
	return 0
}

// IsSquare reports whether W == H.
func (r Rect) IsSquare() bool {
	// TODO
	return false
}

// Scale multiplies both sides by k in place.
func (r *Rect) Scale(k float64) {
	// TODO
}

// Grow adds dw to W and dh to H in place, clamping negative sides to 0.
func (r *Rect) Grow(dw, dh float64) {
	// TODO
}
