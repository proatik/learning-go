package shapes_test

import (
	"math"
	"testing"

	shapes "learngo/course/01-foundations/01-getting-started/exercises/05-exported-names"
)

func near(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestSquareArea(t *testing.T) {
	if got := shapes.SquareArea(3); !near(got, 9) {
		t.Errorf("SquareArea(3) = %v, want 9", got)
	}
}

func TestRectangleArea(t *testing.T) {
	if got := shapes.RectangleArea(2, 4.5); !near(got, 9) {
		t.Errorf("RectangleArea(2, 4.5) = %v, want 9", got)
	}
}

func TestCircleArea(t *testing.T) {
	if shapes.Pi != 3.14159 {
		t.Errorf("Pi = %v, want 3.14159", shapes.Pi)
	}
	if got := shapes.CircleArea(2); !near(got, 12.56636) {
		t.Errorf("CircleArea(2) = %v, want 12.56636", got)
	}
}
