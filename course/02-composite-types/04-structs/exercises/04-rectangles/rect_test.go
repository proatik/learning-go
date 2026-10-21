package rect

import (
	"slices"
	"testing"
)

func R(x0, y0, x1, y1 int) Rect {
	return Rect{Min: Point{X: x0, Y: y0}, Max: Point{X: x1, Y: y1}}
}

func TestArea(t *testing.T) {
	tests := []struct {
		r    Rect
		want int
	}{
		{R(0, 0, 4, 3), 12},
		{R(-2, -2, 2, 2), 16},
		{R(1, 1, 2, 2), 1},
		{R(0, 0, 0, 5), 0},
		{R(3, 3, 7, 3), 0},
		{Rect{}, 0},
	}
	for _, tt := range tests {
		if got := Area(tt.r); got != tt.want {
			t.Errorf("Area(%+v) = %d, want %d", tt.r, got, tt.want)
		}
	}
}

func TestContains(t *testing.T) {
	r := R(0, 0, 4, 3)
	tests := []struct {
		r    Rect
		p    Point
		want bool
	}{
		{r, Point{0, 0}, true},
		{r, Point{3, 2}, true},
		{r, Point{2, 1}, true},
		{r, Point{4, 1}, false},
		{r, Point{1, 3}, false},
		{r, Point{4, 3}, false},
		{r, Point{-1, 0}, false},
		{r, Point{0, -1}, false},
		{R(-5, -5, -1, -1), Point{-3, -2}, true},
		{R(2, 2, 2, 6), Point{2, 3}, false},
		{Rect{}, Point{}, false},
	}
	for _, tt := range tests {
		if got := Contains(tt.r, tt.p); got != tt.want {
			t.Errorf("Contains(%+v, %+v) = %v, want %v", tt.r, tt.p, got, tt.want)
		}
	}
}

func TestIntersect(t *testing.T) {
	r := R(0, 0, 4, 3)
	tests := []struct {
		a, b   Rect
		want   Rect
		wantOK bool
	}{
		{r, R(2, 1, 6, 5), R(2, 1, 4, 3), true},
		{R(2, 1, 6, 5), r, R(2, 1, 4, 3), true},
		{r, r, r, true},
		{r, R(1, 1, 2, 2), R(1, 1, 2, 2), true},
		{R(-3, -3, 3, 3), R(-1, -10, 1, 10), R(-1, -3, 1, 3), true},
		{r, R(4, 0, 6, 3), Rect{}, false},
		{r, R(0, 3, 4, 6), Rect{}, false},
		{r, R(10, 10, 12, 12), Rect{}, false},
		{r, R(1, 1, 1, 2), Rect{}, false},
	}
	for _, tt := range tests {
		got, ok := Intersect(tt.a, tt.b)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("Intersect(%+v, %+v) = %+v, %v; want %+v, %v", tt.a, tt.b, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestBounds(t *testing.T) {
	tests := []struct {
		rs   []Rect
		want Rect
	}{
		{[]Rect{R(0, 0, 4, 3), R(-1, 2, 1, 5)}, R(-1, 0, 4, 5)},
		{[]Rect{R(1, 1, 2, 2)}, R(1, 1, 2, 2)},
		{[]Rect{R(5, 5, 6, 6), R(-3, 0, -2, 1), R(0, -4, 1, 0)}, R(-3, -4, 6, 6)},
		{[]Rect{R(2, 2, 3, 3), R(1, 1, 4, 4)}, R(1, 1, 4, 4)},
		{nil, Rect{}},
		{[]Rect{}, Rect{}},
	}
	for _, tt := range tests {
		in := slices.Clone(tt.rs)
		if got := Bounds(in); got != tt.want {
			t.Errorf("Bounds(%+v) = %+v, want %+v", tt.rs, got, tt.want)
		}
		if !slices.Equal(in, tt.rs) {
			t.Errorf("Bounds modified its input: got %+v, want %+v", in, tt.rs)
		}
	}
}
