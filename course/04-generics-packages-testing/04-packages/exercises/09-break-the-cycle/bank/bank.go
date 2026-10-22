// Package bank holds bank accounts.
package bank

import (
	"errors"

	// TODO: this import, together with audit importing bank, is the cycle.
	"learngo/course/04-generics-packages-testing/04-packages/exercises/09-break-the-cycle/audit"
)

// Errors returned by Deposit and Withdraw.
var (
	ErrAmount = errors.New("bank: amount must be positive")
	ErrFunds  = errors.New("bank: insufficient funds")
)

// Account is a bank account whose activity is recorded in a log.
type Account struct {
	id      string
	balance int
	log     *audit.Log // TODO: depend on the small event.Recorder interface instead
}

// New returns an empty account. Activity is written to log when it is not nil.
func New(id string, log *audit.Log) *Account {
	return &Account{id: id, log: log}
}

// ID returns the account ID.
func (a *Account) ID() string { return a.id }

// Balance returns the current balance.
func (a *Account) Balance() int { return a.balance }

// Deposit adds n to the balance.
func (a *Account) Deposit(n int) error {
	if n <= 0 {
		return ErrAmount
	}
	a.balance += n
	a.record("deposit", n)
	return nil
}

// Withdraw takes n from the balance. A withdrawal larger than the balance
// fails with ErrFunds and is recorded as "rejected".
func (a *Account) Withdraw(n int) error {
	if n <= 0 {
		return ErrAmount
	}
	if n > a.balance {
		a.record("rejected", n)
		return ErrFunds
	}
	a.balance -= n
	a.record("withdraw", n)
	return nil
}

func (a *Account) record(kind string, amount int) {
	if a.log != nil {
		a.log.Record(a, kind, amount) // TODO: send an event.Event instead
	}
}
