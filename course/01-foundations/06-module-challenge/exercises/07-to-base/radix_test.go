package radix

import (
	"math"
	"strconv"
	"testing"
)

func TestToBase(t *testing.T) {
	tests := []struct {
		n, base int
		want    string
	}{
		{0, 2, "0"},
		{0, 36, "0"},
		{1, 2, "1"},
		{10, 2, "1010"},
		{10, 10, "10"},
		{255, 16, "ff"},
		{256, 16, "100"},
		{35, 36, "z"},
		{36, 36, "10"},
		{8, 8, "10"},
		{-1, 2, "-1"},
		{-255, 16, "-ff"},
		{-10, 3, "-101"},
		{1295, 36, "zz"},
	}
	for _, tt := range tests {
		got, err := ToBase(tt.n, tt.base)
		if err != nil || got != tt.want {
			t.Errorf("ToBase(%d, %d) = %q, %v, want %q, nil", tt.n, tt.base, got, err, tt.want)
		}
	}
}

func TestToBaseMatchesStrconv(t *testing.T) {
	values := []int{
		0, 1, -1, 7, -7, 42, 1000, -1000, 123456789, -987654321,
		math.MaxInt32, math.MinInt32, math.MaxInt, math.MinInt, math.MaxInt - 1, math.MinInt + 1,
	}
	for base := 2; base <= 36; base++ {
		for _, n := range values {
			want := strconv.FormatInt(int64(n), base)
			if got, err := ToBase(n, base); err != nil || got != want {
				t.Errorf("ToBase(%d, %d) = %q, %v, want %q, nil", n, base, got, err, want)
			}
		}
	}
}

func TestToBaseInvalid(t *testing.T) {
	for _, base := range []int{-10, -2, 0, 1, 37, 64, 100} {
		if got, err := ToBase(10, base); err == nil || got != "" {
			t.Errorf("ToBase(10, %d) = %q, %v, want \"\", error", base, got, err)
		}
	}
}
