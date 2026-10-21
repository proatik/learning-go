package bank

import "errors"

// ErrInvalidAmount is returned (wrapped) for zero or negative amounts.
var ErrInvalidAmount = errors.New("invalid amount")

// InsufficientFundsError reports a withdrawal larger than the balance.
type InsufficientFundsError struct {
	Need int // the amount requested
	Have int // the balance at that moment
}

// Error returns "insufficient funds: need <Need>, have <Have>".
func (e *InsufficientFundsError) Error() string {
	// TODO
	return ""
}

// Shortfall returns how much money is missing.
func (e *InsufficientFundsError) Shortfall() int {
	// TODO
	return 0
}

// Account is a bank account.
type Account struct {
	Balance int
}

// Withdraw takes amount out of the account.
func (a *Account) Withdraw(amount int) error {
	// TODO
	return nil
}

// Transfer moves amount from a to to.
func (a *Account) Transfer(to *Account, amount int) error {
	// TODO
	return nil
}
