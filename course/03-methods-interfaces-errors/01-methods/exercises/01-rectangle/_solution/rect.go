package rect

// Rect is an axis-aligned rectangle with width W and height H.
type Rect struct {
	W, H float64
}

// Area returns W * H.
func (r Rect) Area() float64 {
	return r.W * r.H
}

// Perimeter returns 2 * (W + H).
func (r Rect) Perimeter() float64 {
	return 2 * (r.W + r.H)
}

// IsSquare reports whether W == H.
func (r Rect) IsSquare() bool {
	return r.W == r.H
}

// Scale multiplies both sides by k in place.
func (r *Rect) Scale(k float64) {
	r.W *= k
	r.H *= k
}

// Grow adds dw to W and dh to H in place, clamping negative sides to 0.
func (r *Rect) Grow(dw, dh float64) {
	r.W = max(r.W+dw, 0)
	r.H = max(r.H+dh, 0)
}
