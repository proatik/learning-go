// Package bank holds bank accounts.
package bank

import (
	"errors"

	"learngo/course/04-generics-packages-testing/04-packages/exercises/09-break-the-cycle/event"
)

// Errors returned by Deposit and Withdraw.
var (
	ErrAmount = errors.New("bank: amount must be positive")
	ErrFunds  = errors.New("bank: insufficient funds")
)

// Account is a bank account whose activity is sent to a Recorder.
type Account struct {
	id      string
	balance int
	rec     event.Recorder
}

// New returns an empty account. Activity is sent to rec when it is not nil.
func New(id string, rec event.Recorder) *Account {
	return &Account{id: id, rec: rec}
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
	if a.rec != nil {
		a.rec.Record(event.Event{Account: a.id, Kind: kind, Amount: amount})
	}
}
