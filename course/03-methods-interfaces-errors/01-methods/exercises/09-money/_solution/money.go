package money

import (
	"fmt"
	"strconv"
	"strings"
)

// Money is an amount of money in whole cents.
type Money int64

// New returns dollars*100 + cents.
func New(dollars, cents int64) Money {
	return Money(dollars*100 + cents)
}

// Add returns m + o.
func (m Money) Add(o Money) Money {
	return m + o
}

// Sub returns m - o.
func (m Money) Sub(o Money) Money {
	return m - o
}

// Times returns m multiplied by n.
func (m Money) Times(n int64) Money {
	return m * Money(n)
}

// Abs returns the absolute value of m.
func (m Money) Abs() Money {
	if m < 0 {
		return -m
	}
	return m
}

// String formats m like "$1,234.56" or "-$0.50".
func (m Money) String() string {
	sign := ""
	if m < 0 {
		sign = "-"
	}
	c := int64(m.Abs())
	return fmt.Sprintf("%s$%s.%02d", sign, group(strconv.FormatInt(c/100, 10)), c%100)
}

// group inserts a comma every three digits from the right.
func group(digits string) string {
	var b strings.Builder
	for i, r := range digits {
		if i > 0 && (len(digits)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(r)
	}
	return b.String()
}

// Parse reads an amount in the format produced by String (commas and the
// cents part are optional) and reports whether s was valid.
func Parse(s string) (Money, bool) {
	s, neg := strings.CutPrefix(s, "-")
	s, ok := strings.CutPrefix(s, "$")
	if !ok {
		return 0, false
	}
	whole, frac, hasFrac := strings.Cut(s, ".")
	cents := int64(0)
	if hasFrac {
		if len(frac) != 2 || !allDigits(frac) {
			return 0, false
		}
		cents = int64(frac[0]-'0')*10 + int64(frac[1]-'0')
	}

	groups := strings.Split(whole, ",")
	if n := len(groups[0]); n == 0 || (len(groups) > 1 && n > 3) {
		return 0, false
	}
	for _, g := range groups[1:] {
		if len(g) != 3 {
			return 0, false
		}
	}
	digits := strings.Join(groups, "")
	if !allDigits(digits) {
		return 0, false
	}
	dollars, err := strconv.ParseInt(digits, 10, 64)
	if err != nil {
		return 0, false
	}

	m := New(dollars, cents)
	if neg {
		m = -m
	}
	return m, true
}

// allDigits reports whether s consists only of ASCII digits.
func allDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}
