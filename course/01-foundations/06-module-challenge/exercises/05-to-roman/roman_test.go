package roman

import (
	"math"
	"strings"
	"testing"
)

func TestToRoman(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{1, "I"},
		{2, "II"},
		{3, "III"},
		{4, "IV"},
		{5, "V"},
		{6, "VI"},
		{9, "IX"},
		{10, "X"},
		{14, "XIV"},
		{27, "XXVII"},
		{40, "XL"},
		{48, "XLVIII"},
		{49, "XLIX"},
		{59, "LIX"},
		{90, "XC"},
		{93, "XCIII"},
		{141, "CXLI"},
		{163, "CLXIII"},
		{400, "CD"},
		{402, "CDII"},
		{444, "CDXLIV"},
		{500, "D"},
		{575, "DLXXV"},
		{900, "CM"},
		{911, "CMXI"},
		{1000, "M"},
		{1024, "MXXIV"},
		{1994, "MCMXCIV"},
		{2027, "MMXXVII"},
		{3000, "MMM"},
		{3888, "MMMDCCCLXXXVIII"},
		{3999, "MMMCMXCIX"},
	}
	for _, tt := range tests {
		got, err := ToRoman(tt.n)
		if err != nil || got != tt.want {
			t.Errorf("ToRoman(%d) = %q, %v, want %q, nil", tt.n, got, err, tt.want)
		}
	}
}

func TestToRomanCanonical(t *testing.T) {
	bad := []string{"IIII", "VV", "XXXX", "LL", "CCCC", "DD", "MMMM", "IC", "IL", "VX"}
	for n := 1; n <= 3999; n++ {
		got, err := ToRoman(n)
		if err != nil || got == "" {
			t.Fatalf("ToRoman(%d) = %q, %v, want a numeral, nil", n, got, err)
		}
		for _, b := range bad {
			if strings.Contains(got, b) {
				t.Fatalf("ToRoman(%d) = %q, contains non-canonical %q", n, got, b)
			}
		}
	}
}

func TestToRomanInvalid(t *testing.T) {
	for _, n := range []int{0, -1, -3999, 4000, 5000, math.MaxInt, math.MinInt} {
		if got, err := ToRoman(n); err == nil || got != "" {
			t.Errorf("ToRoman(%d) = %q, %v, want \"\", error", n, got, err)
		}
	}
}
