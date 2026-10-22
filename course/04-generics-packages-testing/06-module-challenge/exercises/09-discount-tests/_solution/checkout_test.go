package checkout

import (
	"errors"
	"testing"
)

// cart builds a cart from (cents, qty) pairs, naming the SKUs a, b, c, ...
func cart(pairs ...[2]int64) []Line {
	lines := make([]Line, len(pairs))
	for i, p := range pairs {
		lines[i] = Line{SKU: string(rune('a' + i)), Cents: p[0], Qty: int(p[1])}
	}
	return lines
}

// checkTotal fails unless Total(lines, code) returns want and no error.
func checkTotal(t *testing.T, lines []Line, code string, want int64) {
	t.Helper()
	got, err := Total(lines, code)
	if err != nil || got != want {
		t.Errorf("Total(%v, %q) = %d, %v, want %d, <nil>", lines, code, got, err, want)
	}
}

// checkErr fails unless Total(lines, code) returns an error matching want.
func checkErr(t *testing.T, lines []Line, code string, want error) {
	t.Helper()
	got, err := Total(lines, code)
	if !errors.Is(err, want) {
		t.Errorf("Total(%v, %q) = %d, %v, want error %v", lines, code, got, err, want)
	}
}

func TestTotal(t *testing.T) {
	tests := []struct {
		name  string
		lines []Line
		code  string
		want  int64
	}{
		{"plain cart", cart([2]int64{250, 2}, [2]int64{100, 1}), "", 600},
		{"free item is fine", cart([2]int64{0, 3}, [2]int64{100, 1}), "", 100},
		{"bulk starts at 10 units", cart([2]int64{30, 10}), "", 270},
		{"bulk rounds down", cart([2]int64{13, 11}), "", 128},
		{"bulk is per line", cart([2]int64{100, 10}, [2]int64{500, 1}), "", 1400},
		{"SAVE5 at exactly 2500", cart([2]int64{2500, 1}), "SAVE5", 2000},
		{"codes ignore case", cart([2]int64{3000, 1}), "save5", 2500},
		{"HALF takes half of one priciest unit", cart([2]int64{300, 2}, [2]int64{1000, 3}), "HALF", 3100},
		{"HALF in mixed case", cart([2]int64{40, 1}), "Half", 20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkTotal(t, tt.lines, tt.code, tt.want)
		})
	}
}

func TestTotalErrors(t *testing.T) {
	tests := []struct {
		name  string
		lines []Line
		code  string
		want  error
	}{
		{"nil cart", nil, "", ErrEmptyCart},
		{"empty cart", []Line{}, "SAVE5", ErrEmptyCart},
		{"zero quantity", cart([2]int64{100, 1}, [2]int64{100, 0}), "", ErrBadLine},
		{"negative price", cart([2]int64{-1, 1}), "", ErrBadLine},
		{"SAVE5 below 2500", cart([2]int64{2499, 1}), "SAVE5", ErrCodeNotAllowed},
		{"SAVE5 after bulk discount", cart([2]int64{270, 10}), "SAVE5", ErrCodeNotAllowed},
		{"unknown code", cart([2]int64{100, 1}), "FREE", ErrUnknownCode},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkErr(t, tt.lines, tt.code, tt.want)
		})
	}
}
