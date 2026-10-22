package ledger

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
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
	return &Account{ID: id, balance: balance}
}

// Balance returns a's current balance.
func (a *Account) Balance() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.balance
}

// Transfer moves amount from one account to another, as one step: no other
// goroutine sees the money missing from both or present in both. It locks
// the two accounts in ID order (smaller ID first), so opposite transfers
// can't deadlock.
func Transfer(from, to *Account, amount int) error {
	if amount <= 0 {
		return fmt.Errorf("transfer %d: %w", amount, ErrInvalidAmount)
	}
	if from == to {
		return fmt.Errorf("transfer to account %d: %w", from.ID, ErrSameAccount)
	}
	first, second := from, to
	if second.ID < first.ID {
		first, second = second, first
	}
	first.mu.Lock()
	defer first.mu.Unlock()
	second.mu.Lock()
	defer second.mu.Unlock()

	if from.balance < amount {
		return fmt.Errorf("transfer %d from account %d (balance %d): %w", amount, from.ID, from.balance, ErrInsufficientFunds)
	}
	from.balance -= amount
	to.balance += amount
	return nil
}

// Total returns the sum of the balances of the given (distinct) accounts at
// one consistent moment: it holds every account's lock, taken in ID order.
func Total(accounts ...*Account) int {
	sorted := slices.Clone(accounts)
	slices.SortFunc(sorted, func(a, b *Account) int { return cmp.Compare(a.ID, b.ID) })
	for _, a := range sorted {
		a.mu.Lock()
		defer a.mu.Unlock()
	}
	sum := 0
	for _, a := range sorted {
		sum += a.balance
	}
	return sum
}
