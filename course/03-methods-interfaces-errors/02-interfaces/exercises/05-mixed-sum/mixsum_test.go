package mixsum_test

import (
	"math"
	"slices"
	"testing"

	mixsum "learngo/course/03-methods-interfaces-errors/02-interfaces/exercises/05-mixed-sum"
)

type cents int

func (c cents) Value() float64 { return float64(c) / 100 }

type weight struct{ kg float64 }

func (w *weight) Value() float64 { return w.kg }

func TestSum(t *testing.T) {
	n, m := 7, -2
	var nilInt *int
	tests := []struct {
		in          []any
		want        float64
		wantSkipped int
	}{
		{nil, 0, 0},
		{[]any{}, 0, 0},
		{[]any{1, int64(2), 2.5}, 5.5, 0},
		{[]any{1, "4", nil, true}, 1, 3},
		{[]any{float32(0.5), cents(250)}, 3, 0},
		{[]any{&n, nilInt}, 7, 1},
		{[]any{&n, &m, &n}, 12, 0},
		{[]any{-1, -1.5, int64(-10)}, -12.5, 0},
		{[]any{&weight{kg: 1.25}, weight{kg: 9}}, 1.25, 1}, // only *weight is a Valuer
		{[]any{int32(5), uint(5), []int{1}, map[string]int{}}, 0, 4},
		{[]any{"1", "2", "3"}, 0, 3},
		{[]any{nil, nil}, 0, 2},
		{[]any{math.MaxInt32, 1}, math.MaxInt32 + 1, 0},
	}
	for _, tt := range tests {
		got, skipped := mixsum.Sum(tt.in)
		if math.Abs(got-tt.want) > 1e-9 || skipped != tt.wantSkipped {
			t.Errorf("Sum(%v) = %v, %d, want %v, %d", tt.in, got, skipped, tt.want, tt.wantSkipped)
		}
	}
}

func TestStrings(t *testing.T) {
	tests := []struct {
		in   []any
		want []string
	}{
		{nil, nil},
		{[]any{1, 2}, nil},
		{[]any{"a", 1, "b", nil}, []string{"a", "b"}},
		{[]any{"", "ünï", 'x', []string{"no"}}, []string{"", "ünï"}},
		{[]any{"only"}, []string{"only"}},
	}
	for _, tt := range tests {
		if got := mixsum.Strings(tt.in); !slices.Equal(got, tt.want) {
			t.Errorf("Strings(%v) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
