// Package money represents non-negative amounts of money as integer cents.
package money

import (
	"errors"
	"fmt"
)

// Errors returned by the constructors.
var (
	ErrNegative = errors.New("money: negative amount")
	ErrCents    = errors.New("money: cents must be between 0 and 99")
)

// Amount is a non-negative amount of money. The zero value is 0.00.
//
// The field is unexported, so other packages can only create an Amount
// through New or FromCents, which validate it.
type Amount struct {
	cents int64
}

// New returns units.cents as an Amount. It fails with ErrNegative when either
// part is negative and with ErrCents when cents is above 99.
func New(units, cents int64) (Amount, error) {
	if units < 0 || cents < 0 {
		return Amount{}, ErrNegative
	}
	if cents > 99 {
		return Amount{}, ErrCents
	}
	return Amount{cents: units*100 + cents}, nil
}

// FromCents returns an Amount of the given cents, or ErrNegative.
func FromCents(cents int64) (Amount, error) {
	if cents < 0 {
		return Amount{}, ErrNegative
	}
	return Amount{cents: cents}, nil
}

// Cents returns the amount in cents.
func (a Amount) Cents() int64 { return a.cents }

// Add returns a + b.
func (a Amount) Add(b Amount) Amount { return Amount{cents: a.cents + b.cents} }

// Times returns a * n, or the zero Amount when n is negative.
func (a Amount) Times(n int) Amount {
	if n < 0 {
		return Amount{}
	}
	return Amount{cents: a.cents * int64(n)}
}

// String formats the amount as "units.cc", e.g. "12.05".
func (a Amount) String() string {
	return fmt.Sprintf("%d.%02d", a.cents/100, a.cents%100)
}
