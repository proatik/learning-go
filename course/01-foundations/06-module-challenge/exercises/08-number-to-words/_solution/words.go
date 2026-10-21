package words

import "fmt"

var (
	// small[i] spells i for 0 <= i < 20.
	small = []string{
		"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine",
		"ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen",
		"seventeen", "eighteen", "nineteen",
	}
	// tens[i] spells i*10 for 2 <= i <= 9.
	tens = []string{"", "", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety"}
)

// NumberToWords spells n (0 to 999999) in lowercase English words,
// e.g. 123 is "one hundred twenty-three".
func NumberToWords(n int) (string, error) {
	if n < 0 || n > 999_999 {
		return "", fmt.Errorf("words: %d is out of range 0..999999", n)
	}
	if n == 0 {
		return "zero", nil
	}

	thousands, rest := n/1000, n%1000
	s := ""
	if thousands > 0 {
		s = spell999(thousands) + " thousand"
	}
	if rest > 0 {
		if s != "" {
			s += " "
		}
		s += spell999(rest)
	}
	return s, nil
}

// spell999 spells n for 1 <= n <= 999.
func spell999(n int) string {
	hundreds, rest := n/100, n%100
	s := ""
	if hundreds > 0 {
		s = small[hundreds] + " hundred"
	}
	if rest == 0 {
		return s
	}
	if s != "" {
		s += " "
	}
	switch {
	case rest < 20:
		s += small[rest]
	case rest%10 == 0:
		s += tens[rest/10]
	default:
		s += tens[rest/10] + "-" + small[rest%10]
	}
	return s
}
