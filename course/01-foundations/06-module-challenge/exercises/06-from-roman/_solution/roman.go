package roman

import (
	"errors"
	"fmt"
	"strings"
)

var (
	values  = []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	symbols = []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
)

// FromRoman parses a canonical uppercase Roman numeral (I to MMMCMXCIX).
func FromRoman(s string) (int, error) {
	if s == "" {
		return 0, errors.New("roman: empty numeral")
	}

	total := 0
	for i := range len(s) {
		v := symbolValue(s[i])
		if v == 0 {
			return 0, fmt.Errorf("roman: invalid character %q in %q", s[i], s)
		}
		// Subtract a symbol that is smaller than its right neighbour.
		if i+1 < len(s) && v < symbolValue(s[i+1]) {
			total -= v
		} else {
			total += v
		}
	}

	// The rule above also "accepts" forms like IIII or IC. Only the
	// canonical spelling of total round-trips.
	if total < 1 || total > 3999 || toRoman(total) != s {
		return 0, fmt.Errorf("roman: %q is not a canonical numeral", s)
	}
	return total, nil
}

// symbolValue returns the value of one Roman symbol, or 0 if c is not one.
func symbolValue(c byte) int {
	switch c {
	case 'I':
		return 1
	case 'V':
		return 5
	case 'X':
		return 10
	case 'L':
		return 50
	case 'C':
		return 100
	case 'D':
		return 500
	case 'M':
		return 1000
	}
	return 0
}

// toRoman converts n (assumed to be 1..3999) to its canonical numeral.
func toRoman(n int) string {
	var sb strings.Builder
	for i, v := range values {
		for n >= v {
			sb.WriteString(symbols[i])
			n -= v
		}
	}
	return sb.String()
}
