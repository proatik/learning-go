package roman

import (
	"fmt"
	"strings"
)

// Value/symbol pairs from largest to smallest, including the six
// subtractive pairs, so a greedy loop always yields the canonical form.
var (
	values  = []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	symbols = []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
)

// ToRoman converts n (1 to 3999) to an uppercase Roman numeral.
func ToRoman(n int) (string, error) {
	if n < 1 || n > 3999 {
		return "", fmt.Errorf("roman: %d is out of range 1..3999", n)
	}

	var sb strings.Builder
	for i, v := range values {
		// Take the largest value that still fits as many times as it fits.
		for n >= v {
			sb.WriteString(symbols[i])
			n -= v
		}
	}
	return sb.String(), nil
}
