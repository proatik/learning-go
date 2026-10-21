package money

import "fmt"

// Money is an amount of money stored as a whole number of cents.
type Money struct {
	Cents int64
}

// New returns the amount dollars*100 + cents.
func New(dollars, cents int64) Money {
	return Money{Cents: dollars*100 + cents}
}

// Add returns a + b.
func Add(a, b Money) Money {
	return Money{Cents: a.Cents + b.Cents}
}

// Sum returns the total of ms (zero for an empty slice).
func Sum(ms []Money) Money {
	var total Money
	for _, m := range ms {
		total = Add(total, m)
	}
	return total
}

// Format returns m as "$D.CC", or "-$D.CC" when m is negative.
func Format(m Money) string {
	c := m.Cents
	sign := ""
	if c < 0 {
		sign = "-"
		c = -c
	}
	return fmt.Sprintf("%s$%d.%02d", sign, c/100, c%100)
}

// Split divides m into n parts that differ by at most one cent and add up
// to m; the first parts receive the leftover cents. It returns nil if n <= 0.
func Split(m Money, n int) []Money {
	if n <= 0 {
		return nil
	}
	base := m.Cents / int64(n)
	extra := m.Cents % int64(n) // the first `extra` parts get one more cent
	parts := make([]Money, n)
	for i := range parts {
		parts[i] = Money{Cents: base}
		if int64(i) < extra {
			parts[i].Cents++
		}
	}
	return parts
}
