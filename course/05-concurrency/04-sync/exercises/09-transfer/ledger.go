package ledger

import (
	"errors"
	"sync"
)

var (
	// ErrInvalidAmount is returned for a transfer of 0 or less.
	ErrInvalidAmount = errors.New("ledger: amount must be positive")
	// ErrSameAccount is returned for a transfer from an account to itself.
	ErrSameAccount = errors.New("ledger: cannot transfer to the same account")
	// ErrInsufficientFunds is returned when from's balance is too small.
	ErrInsufficientFunds = errors.New("ledger: insufficient funds")
)

// Account is a bank account. Every account has a unique ID, which also fixes
// the order in which locks are taken.
type Account struct {
	ID      int
	mu      sync.Mutex // guards balance
	balance int
}

// NewAccount returns an account with the given ID and starting balance.
func NewAccount(id, balance int) *Account {
	// TODO
	return &Account{}
}

// Balance returns a's current balance.
func (a *Account) Balance() int {
	// TODO
	return 0
}

// Transfer moves amount from one account to another, as one step: no other
// goroutine sees the money missing from both or present in both. It locks
// the two accounts in ID order (smaller ID first), so opposite transfers
// can't deadlock.
func Transfer(from, to *Account, amount int) error {
	// TODO
	return nil
}

// Total returns the sum of the balances of the given (distinct) accounts at
// one consistent moment: it holds every account's lock, taken in ID order.
func Total(accounts ...*Account) int {
	// TODO
	return 0
}
