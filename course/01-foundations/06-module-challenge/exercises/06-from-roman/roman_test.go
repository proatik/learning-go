package roman

import (
	"strings"
	"testing"
)

func TestFromRoman(t *testing.T) {
	tests := []struct {
		s    string
		want int
	}{
		{"I", 1},
		{"III", 3},
		{"IV", 4},
		{"V", 5},
		{"IX", 9},
		{"XIV", 14},
		{"XL", 40},
		{"XLII", 42},
		{"XLIX", 49},
		{"XC", 90},
		{"XCIX", 99},
		{"CD", 400},
		{"CDXLIV", 444},
		{"D", 500},
		{"CM", 900},
		{"M", 1000},
		{"MCMXCIV", 1994},
		{"MMXXVII", 2027},
		{"MMMDCCCLXXXVIII", 3888},
		{"MMMCMXCIX", 3999},
	}
	for _, tt := range tests {
		got, err := FromRoman(tt.s)
		if err != nil || got != tt.want {
			t.Errorf("FromRoman(%q) = %d, %v, want %d, nil", tt.s, got, err, tt.want)
		}
	}
}

// refRoman is the test's own reference converter.
func refRoman(n int) string {
	var sb strings.Builder
	for _, p := range []struct {
		v int
		s string
	}{
		{1000, "M"}, {900, "CM"}, {500, "D"}, {400, "CD"}, {100, "C"}, {90, "XC"},
		{50, "L"}, {40, "XL"}, {10, "X"}, {9, "IX"}, {5, "V"}, {4, "IV"}, {1, "I"},
	} {
		for n >= p.v {
			sb.WriteString(p.s)
			n -= p.v
		}
	}
	return sb.String()
}

func TestFromRomanRoundTrip(t *testing.T) {
	for n := 1; n <= 3999; n++ {
		s := refRoman(n)
		if got, err := FromRoman(s); err != nil || got != n {
			t.Fatalf("FromRoman(%q) = %d, %v, want %d, nil", s, got, err, n)
		}
	}
}

func TestFromRomanInvalid(t *testing.T) {
	invalid := []string{
		"",           // empty
		"iv",         // lowercase
		"Iv",         // mixed case
		" IV",        // leading space
		"IV ",        // trailing space
		"I V",        // inner space
		"A",          // not a symbol
		"XIIZ",       // not a symbol
		"IIII",       // should be IV
		"VIIII",      // should be IX
		"XXXX",       // should be XL
		"CCCC",       // should be CD
		"MMMM",       // 4000: out of range
		"VV",         // should be X
		"LL",         // should be C
		"DD",         // should be M
		"IC",         // should be XCIX
		"IL",         // should be XLIX
		"XM",         // should be CMXC
		"VX",         // should be V
		"IIV",        // not a valid subtraction
		"IXI",        // should be X
		"IVI",        // should be V
		"MCMC",       // should be MM
		"CMM",        // should be MCM
		"MCMXCIVI",   // should be MCMXCV
		"MMMCMXCIXI", // 4000: out of range
		"Ⅻ",          // a Unicode Roman numeral, not ASCII letters
	}
	for _, s := range invalid {
		if got, err := FromRoman(s); err == nil || got != 0 {
			t.Errorf("FromRoman(%q) = %d, %v, want 0, error", s, got, err)
		}
	}
}
