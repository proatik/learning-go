package money

// Money is an amount of money in whole cents.
type Money int64

// New returns dollars*100 + cents.
func New(dollars, cents int64) Money {
	// TODO
	return 0
}

// Add returns m + o.
func (m Money) Add(o Money) Money {
	// TODO
	return 0
}

// Sub returns m - o.
func (m Money) Sub(o Money) Money {
	// TODO
	return 0
}

// Times returns m multiplied by n.
func (m Money) Times(n int64) Money {
	// TODO
	return 0
}

// Abs returns the absolute value of m.
func (m Money) Abs() Money {
	// TODO
	return 0
}

// String formats m like "$1,234.56" or "-$0.50".
func (m Money) String() string {
	// TODO
	return ""
}

// Parse reads an amount in the format produced by String (commas and the
// cents part are optional) and reports whether s was valid.
func Parse(s string) (Money, bool) {
	// TODO
	return 0, false
}
