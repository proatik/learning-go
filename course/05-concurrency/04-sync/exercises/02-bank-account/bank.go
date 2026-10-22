package bank

import (
	"errors"
	"sync"
)

// ErrInvalidAmount is returned for a deposit or withdrawal of 0 or less.
var ErrInvalidAmount = errors.New("bank: amount must be positive")

// ErrInsufficientFunds is returned when a withdrawal exceeds the balance.
var ErrInsufficientFunds = errors.New("bank: insufficient funds")

// Account is a bank account that is safe for concurrent use. The zero value
// is an empty account.
type Account struct {
	mu      sync.Mutex // guards balance
	balance int
}

// Deposit adds amount to the balance.
func (a *Account) Deposit(amount int) error {
	// TODO
	return nil
}

// Withdraw takes amount from the balance. It fails, changing nothing, if the
// balance is too small.
func (a *Account) Withdraw(amount int) error {
	// TODO
	return nil
}

// Balance returns the current balance.
func (a *Account) Balance() int {
	// TODO
	return 0
}
