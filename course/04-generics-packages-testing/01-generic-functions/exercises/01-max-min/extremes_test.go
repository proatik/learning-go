package extremes_test

import (
	"testing"

	ext "learngo/course/04-generics-packages-testing/01-generic-functions/exercises/01-max-min"
)

// Celsius is a named type: cmp.Ordered must accept it through ~float64.
type Celsius float64

func TestMaxOfInts(t *testing.T) {
	tests := []struct {
		in     []int
		want   int
		wantOK bool
	}{
		{[]int{3, 9, 2}, 9, true},
		{[]int{9, 3, 2}, 9, true},
		{[]int{2, 3, 9}, 9, true},
		{[]int{-5, -1, -9}, -1, true},
		{[]int{7}, 7, true},
		{[]int{4, 4, 4}, 4, true},
		{[]int{}, 0, false},
		{nil, 0, false},
	}
	for _, tt := range tests {
		got, ok := ext.MaxOf(tt.in)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("MaxOf(%v) = %v, %t, want %v, %t", tt.in, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestMinOfInts(t *testing.T) {
	tests := []struct {
		in     []int
		want   int
		wantOK bool
	}{
		{[]int{3, 9, 2}, 2, true},
		{[]int{2, 9, 3}, 2, true},
		{[]int{-5, -1, -9}, -9, true},
		{[]int{0}, 0, true},
		{[]int{}, 0, false},
		{nil, 0, false},
	}
	for _, tt := range tests {
		got, ok := ext.MinOf(tt.in)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("MinOf(%v) = %v, %t, want %v, %t", tt.in, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestMaxMinOtherTypes(t *testing.T) {
	if got, ok := ext.MinOf([]float64{1.5, -2, 0}); got != -2 || !ok {
		t.Errorf("MinOf(%v) = %v, %t, want %v, %t", []float64{1.5, -2, 0}, got, ok, -2.0, true)
	}
	if got, ok := ext.MaxOf([]float64{-0.5, -0.25}); got != -0.25 || !ok {
		t.Errorf("MaxOf(%v) = %v, %t, want %v, %t", []float64{-0.5, -0.25}, got, ok, -0.25, true)
	}
	words := []string{"pear", "apple", "zoo", "Zebra"}
	if got, ok := ext.MaxOf(words); got != "zoo" || !ok {
		t.Errorf("MaxOf(%q) = %q, %t, want %q, %t", words, got, ok, "zoo", true)
	}
	if got, ok := ext.MinOf(words); got != "Zebra" || !ok {
		t.Errorf("MinOf(%q) = %q, %t, want %q, %t", words, got, ok, "Zebra", true)
	}
	if got, ok := ext.MinOf[string](nil); got != "" || ok {
		t.Errorf("MinOf(nil) = %q, %t, want %q, %t", got, ok, "", false)
	}
	temps := []Celsius{21.5, 19, 23.25}
	if got, ok := ext.MaxOf(temps); got != 23.25 || !ok {
		t.Errorf("MaxOf(%v) = %v, %t, want %v, %t", temps, got, ok, Celsius(23.25), true)
	}
	if got, ok := ext.MinOf(temps); got != 19 || !ok {
		t.Errorf("MinOf(%v) = %v, %t, want %v, %t", temps, got, ok, Celsius(19), true)
	}
}

func TestMaxOfDoesNotModify(t *testing.T) {
	in := []int{3, 1, 2}
	ext.MaxOf(in)
	ext.MinOf(in)
	if in[0] != 3 || in[1] != 1 || in[2] != 2 {
		t.Errorf("MaxOf/MinOf modified the input: %v, want [3 1 2]", in)
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		v, lo, hi, want int
	}{
		{5, 0, 10, 5},
		{15, 0, 10, 10},
		{-3, 0, 10, 0},
		{0, 0, 10, 0},
		{10, 0, 10, 10},
		{7, 7, 7, 7},
		{5, 10, 0, 5},
		{20, 10, 0, 10},
		{-20, 10, 0, 0},
		{-7, -10, -5, -7},
	}
	for _, tt := range tests {
		if got := ext.Clamp(tt.v, tt.lo, tt.hi); got != tt.want {
			t.Errorf("Clamp(%d, %d, %d) = %d, want %d", tt.v, tt.lo, tt.hi, got, tt.want)
		}
	}
}

func TestClampOtherTypes(t *testing.T) {
	if got := ext.Clamp("m", "a", "f"); got != "f" {
		t.Errorf("Clamp(%q, %q, %q) = %q, want %q", "m", "a", "f", got, "f")
	}
	if got := ext.Clamp("c", "a", "f"); got != "c" {
		t.Errorf("Clamp(%q, %q, %q) = %q, want %q", "c", "a", "f", got, "c")
	}
	if got := ext.Clamp(1.5, 0.0, 1.0); got != 1.0 {
		t.Errorf("Clamp(%v, %v, %v) = %v, want %v", 1.5, 0.0, 1.0, got, 1.0)
	}
	if got := ext.Clamp(Celsius(-40), -10, 50); got != -10 {
		t.Errorf("Clamp(%v, %v, %v) = %v, want %v", -40, -10, 50, got, -10)
	}
}
