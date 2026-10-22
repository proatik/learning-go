package bank

import (
	"errors"
	"fmt"
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
	if amount <= 0 {
		return fmt.Errorf("deposit %d: %w", amount, ErrInvalidAmount)
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.balance += amount
	return nil
}

// Withdraw takes amount from the balance. It fails, changing nothing, if the
// balance is too small.
func (a *Account) Withdraw(amount int) error {
	if amount <= 0 {
		return fmt.Errorf("withdraw %d: %w", amount, ErrInvalidAmount)
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if amount > a.balance {
		return fmt.Errorf("withdraw %d from %d: %w", amount, a.balance, ErrInsufficientFunds)
	}
	a.balance -= amount
	return nil
}

// Balance returns the current balance.
func (a *Account) Balance() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.balance
}
