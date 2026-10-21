package money

// Money is an amount of money stored as a whole number of cents.
type Money struct {
	Cents int64
}

// New returns the amount dollars*100 + cents.
func New(dollars, cents int64) Money {
	// TODO
	return Money{}
}

// Add returns a + b.
func Add(a, b Money) Money {
	// TODO
	return Money{}
}

// Sum returns the total of ms (zero for an empty slice).
func Sum(ms []Money) Money {
	// TODO
	return Money{}
}

// Format returns m as "$D.CC", or "-$D.CC" when m is negative.
func Format(m Money) string {
	// TODO
	return ""
}

// Split divides m into n parts that differ by at most one cent and add up
// to m; the first parts receive the leftover cents. It returns nil if n <= 0.
func Split(m Money, n int) []Money {
	// TODO
	return nil
}
