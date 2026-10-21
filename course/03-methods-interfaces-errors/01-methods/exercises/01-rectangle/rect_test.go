package rect

import "testing"

func TestArea(t *testing.T) {
	tests := []struct {
		r    Rect
		want float64
	}{
		{Rect{3, 4}, 12},
		{Rect{}, 0},
		{Rect{1.5, 2}, 3},
		{Rect{10, 0}, 0},
	}
	for _, tt := range tests {
		r := tt.r
		if got := r.Area(); got != tt.want {
			t.Errorf("%+v.Area() = %v, want %v", tt.r, got, tt.want)
		}
		if r != tt.r {
			t.Errorf("Area changed the receiver: %+v, want %+v", r, tt.r)
		}
	}
}

func TestPerimeter(t *testing.T) {
	tests := []struct {
		r    Rect
		want float64
	}{
		{Rect{3, 4}, 14},
		{Rect{}, 0},
		{Rect{1.5, 2}, 7},
		{Rect{5, 5}, 20},
	}
	for _, tt := range tests {
		if got := tt.r.Perimeter(); got != tt.want {
			t.Errorf("%+v.Perimeter() = %v, want %v", tt.r, got, tt.want)
		}
	}
}

func TestIsSquare(t *testing.T) {
	tests := []struct {
		r    Rect
		want bool
	}{
		{Rect{2, 2}, true},
		{Rect{}, true},
		{Rect{2, 3}, false},
		{Rect{0.5, 0.5}, true},
	}
	for _, tt := range tests {
		if got := tt.r.IsSquare(); got != tt.want {
			t.Errorf("%+v.IsSquare() = %v, want %v", tt.r, got, tt.want)
		}
	}
}

func TestScale(t *testing.T) {
	tests := []struct {
		r    Rect
		k    float64
		want Rect
	}{
		{Rect{2, 3}, 10, Rect{20, 30}},
		{Rect{2, 3}, 1, Rect{2, 3}},
		{Rect{2, 3}, 0, Rect{0, 0}},
		{Rect{4, 8}, 0.5, Rect{2, 4}},
	}
	for _, tt := range tests {
		r := tt.r
		r.Scale(tt.k)
		if r != tt.want {
			t.Errorf("%+v.Scale(%v) = %+v, want %+v", tt.r, tt.k, r, tt.want)
		}
	}
}

func TestScaleTwice(t *testing.T) {
	r := Rect{1, 2}
	r.Scale(2)
	r.Scale(3)
	if want := (Rect{6, 12}); r != want {
		t.Errorf("Scale(2) then Scale(3) on {1 2} = %+v, want %+v", r, want)
	}
}

func TestGrow(t *testing.T) {
	tests := []struct {
		r      Rect
		dw, dh float64
		want   Rect
	}{
		{Rect{2, 3}, 1, -1, Rect{3, 2}},
		{Rect{2, 3}, 0, 0, Rect{2, 3}},
		{Rect{2, 3}, -5, 0, Rect{0, 3}},
		{Rect{2, 3}, 0, -3, Rect{2, 0}},
		{Rect{}, 4, 5, Rect{4, 5}},
		{Rect{1, 1}, -1, -10, Rect{0, 0}},
	}
	for _, tt := range tests {
		r := tt.r
		r.Grow(tt.dw, tt.dh)
		if r != tt.want {
			t.Errorf("%+v.Grow(%v, %v) = %+v, want %+v", tt.r, tt.dw, tt.dh, r, tt.want)
		}
	}
}
