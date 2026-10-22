package numeric_test

import (
	"slices"
	"testing"

	num "learngo/course/04-generics-packages-testing/01-generic-functions/exercises/05-number-constraint"
)

// Named types: Number needs ~ to accept them.
type (
	Celsius  float64
	Score    int
	Readings []Celsius
)

// Compile-time checks: Scale must return the caller's slice type S.
var (
	_ Readings = num.Scale(Readings{}, 2)
	_ []int8   = num.Scale([]int8{}, 2)
)

func TestSum(t *testing.T) {
	ints := []struct {
		in   []int
		want int
	}{
		{[]int{1, 2, 3}, 6},
		{[]int{-5, 5}, 0},
		{[]int{-1, -2}, -3},
		{[]int{42}, 42},
		{[]int{}, 0},
		{nil, 0},
	}
	for _, tt := range ints {
		if got := num.Sum(tt.in); got != tt.want {
			t.Errorf("Sum(%v) = %d, want %d", tt.in, got, tt.want)
		}
	}

	if got := num.Sum([]float64{1.5, 2.25}); got != 3.75 {
		t.Errorf("Sum(%v) = %v, want %v", []float64{1.5, 2.25}, got, 3.75)
	}
	if got := num.Sum([]Celsius{20.5, 21}); got != 41.5 {
		t.Errorf("Sum(%v) = %v, want %v", []Celsius{20.5, 21}, got, Celsius(41.5))
	}
	if got := num.Sum([]Score{10, 20, 30}); got != 60 {
		t.Errorf("Sum(%v) = %v, want %v", []Score{10, 20, 30}, got, Score(60))
	}
	if got := num.Sum([]int8{100, -50, 7}); got != 57 {
		t.Errorf("Sum(%v) = %v, want %v", []int8{100, -50, 7}, got, 57)
	}
	if got := num.Sum([]float32{0.5, 0.25}); got != 0.75 {
		t.Errorf("Sum(%v) = %v, want %v", []float32{0.5, 0.25}, got, 0.75)
	}
	if got := num.Sum([]int64{1 << 40, 1 << 40}); got != 1<<41 {
		t.Errorf("Sum(%v) = %v, want %v", []int64{1 << 40, 1 << 40}, got, int64(1<<41))
	}
	if got := num.Sum(Readings{1, 2}); got != 3 {
		t.Errorf("Sum(%v) = %v, want %v", Readings{1, 2}, got, Celsius(3))
	}
}

func TestAbs(t *testing.T) {
	ints := []struct{ in, want int }{
		{-7, 7}, {7, 7}, {0, 0}, {-1, 1},
	}
	for _, tt := range ints {
		if got := num.Abs(tt.in); got != tt.want {
			t.Errorf("Abs(%d) = %d, want %d", tt.in, got, tt.want)
		}
	}
	if got := num.Abs(Celsius(-3.5)); got != 3.5 {
		t.Errorf("Abs(%v) = %v, want %v", Celsius(-3.5), got, Celsius(3.5))
	}
	if got := num.Abs(-0.25); got != 0.25 {
		t.Errorf("Abs(%v) = %v, want %v", -0.25, got, 0.25)
	}
	if got := num.Abs(Score(-12)); got != 12 {
		t.Errorf("Abs(%v) = %v, want %v", Score(-12), got, Score(12))
	}
	if got := num.Abs(int16(-300)); got != 300 {
		t.Errorf("Abs(%v) = %v, want %v", int16(-300), got, int16(300))
	}
}

func TestScale(t *testing.T) {
	r := Readings{1, 2.5, -4}
	got := num.Scale(r, 2)
	if want := (Readings{2, 5, -8}); !slices.Equal(got, want) {
		t.Errorf("Scale(%v, 2) = %v, want %v", r, got, want)
	}
	if want := (Readings{1, 2.5, -4}); !slices.Equal(r, want) {
		t.Errorf("Scale(%v, 2) modified the input: %v, want %v", want, r, want)
	}

	ints := []int{1, 2, 3}
	if got, want := num.Scale(ints, -3), []int{-3, -6, -9}; !slices.Equal(got, want) {
		t.Errorf("Scale(%v, -3) = %v, want %v", ints, got, want)
	}
	if got, want := num.Scale(ints, 0), []int{0, 0, 0}; !slices.Equal(got, want) {
		t.Errorf("Scale(%v, 0) = %v, want %v", ints, got, want)
	}
	scores := []Score{5, 10}
	if got, want := num.Scale(scores, 10), []Score{50, 100}; !slices.Equal(got, want) {
		t.Errorf("Scale(%v, 10) = %v, want %v", scores, got, want)
	}
	if got := num.Scale(Readings{}, 2); len(got) != 0 {
		t.Errorf("Scale(Readings{}, 2) = %v, want []", got)
	}
	if got := num.Scale([]int(nil), 2); len(got) != 0 {
		t.Errorf("Scale(nil, 2) = %v, want []", got)
	}
}

func TestScaleReturnsNewSlice(t *testing.T) {
	in := []int{1, 2}
	got := num.Scale(in, 1)
	if len(got) != len(in) {
		t.Fatalf("Scale(%v, 1) = %v, want %v", in, got, in)
	}
	if &got[0] == &in[0] {
		t.Errorf("Scale(%v, 1) shares memory with the input, want a new slice", in)
	}
}
