package bank

import (
	"errors"
	"fmt"
)

// ErrInvalidAmount is returned (wrapped) for zero or negative amounts.
var ErrInvalidAmount = errors.New("invalid amount")

// InsufficientFundsError reports a withdrawal larger than the balance.
type InsufficientFundsError struct {
	Need int // the amount requested
	Have int // the balance at that moment
}

// Error returns "insufficient funds: need <Need>, have <Have>".
func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("insufficient funds: need %d, have %d", e.Need, e.Have)
}

// Shortfall returns how much money is missing.
func (e *InsufficientFundsError) Shortfall() int {
	return e.Need - e.Have
}

// Account is a bank account.
type Account struct {
	Balance int
}

// Withdraw takes amount out of the account.
func (a *Account) Withdraw(amount int) error {
	if amount <= 0 {
		return fmt.Errorf("withdraw %d: %w", amount, ErrInvalidAmount)
	}
	if amount > a.Balance {
		return &InsufficientFundsError{Need: amount, Have: a.Balance}
	}
	a.Balance -= amount
	return nil
}

// Transfer moves amount from a to to.
func (a *Account) Transfer(to *Account, amount int) error {
	if err := a.Withdraw(amount); err != nil {
		return fmt.Errorf("transfer: %w", err)
	}
	to.Balance += amount
	return nil
}
