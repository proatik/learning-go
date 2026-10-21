package temperature

import (
	"math"
	"testing"
)

func near(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestCToF(t *testing.T) {
	tests := []struct{ c, f float64 }{{0, 32}, {100, 212}, {-40, -40}, {37, 98.6}}
	for _, tt := range tests {
		if got := CToF(tt.c); !near(got, tt.f) {
			t.Errorf("CToF(%v) = %v, want %v", tt.c, got, tt.f)
		}
	}
}

func TestFToC(t *testing.T) {
	tests := []struct{ f, c float64 }{{32, 0}, {212, 100}, {-40, -40}, {50, 10}}
	for _, tt := range tests {
		if got := FToC(tt.f); !near(got, tt.c) {
			t.Errorf("FToC(%v) = %v, want %v", tt.f, got, tt.c)
		}
	}
}

func TestDescribe(t *testing.T) {
	if got, want := Describe(37), "37.0°C = 98.6°F"; got != want {
		t.Errorf("Describe(37) = %q, want %q", got, want)
	}
}
