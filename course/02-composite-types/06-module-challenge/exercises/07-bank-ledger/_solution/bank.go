package bank

import (
	"errors"
	"fmt"
)

// Account is a bank account. Balance is in cents.
type Account struct {
	ID      string
	Owner   string
	Balance int64 // cents
}

// Tx is one money movement. From is "" for a deposit, To is "" for a
// withdrawal.
type Tx struct {
	From   string
	To     string
	Amount int64
}

// Ledger holds every account by ID and the history of successful
// transactions.
type Ledger struct {
	Accounts map[string]*Account
	History  []Tx
}

// Errors returned by ledger operations.
var (
	ErrAccountExists = errors.New("bank: account already exists")
	ErrNoAccount     = errors.New("bank: no such account")
	ErrSameAccount   = errors.New("bank: transfer to the same account")
	ErrInvalidAmount = errors.New("bank: invalid amount")
	ErrInsufficient  = errors.New("bank: insufficient funds")
)

// NewLedger returns an empty, ready-to-use ledger.
func NewLedger() *Ledger {
	return &Ledger{Accounts: make(map[string]*Account)}
}

// account looks up id, wrapping ErrNoAccount with the ID when missing.
func account(l *Ledger, id string) (*Account, error) {
	a, ok := l.Accounts[id]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrNoAccount, id)
	}
	return a, nil
}

// Open creates account id with an initial balance in cents.
func Open(l *Ledger, id, owner string, initial int64) error {
	switch {
	case id == "":
		return errors.New("bank: empty account id")
	case initial < 0:
		return fmt.Errorf("%w: initial balance %d", ErrInvalidAmount, initial)
	}
	if _, ok := l.Accounts[id]; ok {
		return fmt.Errorf("%w: %q", ErrAccountExists, id)
	}
	l.Accounts[id] = &Account{ID: id, Owner: owner, Balance: initial}
	if initial > 0 {
		l.History = append(l.History, Tx{To: id, Amount: initial})
	}
	return nil
}

// Deposit adds amount cents to account id.
func Deposit(l *Ledger, id string, amount int64) error {
	a, err := account(l, id)
	if err != nil {
		return err
	}
	if amount <= 0 {
		return fmt.Errorf("%w: %d", ErrInvalidAmount, amount)
	}
	a.Balance += amount // a is a pointer: this updates the stored account
	l.History = append(l.History, Tx{To: id, Amount: amount})
	return nil
}

// Withdraw removes amount cents from account id.
func Withdraw(l *Ledger, id string, amount int64) error {
	a, err := account(l, id)
	if err != nil {
		return err
	}
	if amount <= 0 {
		return fmt.Errorf("%w: %d", ErrInvalidAmount, amount)
	}
	if amount > a.Balance {
		return fmt.Errorf("%w: %q has %d, need %d", ErrInsufficient, id, a.Balance, amount)
	}
	a.Balance -= amount
	l.History = append(l.History, Tx{From: id, Amount: amount})
	return nil
}

// Transfer moves amount cents from one account to another.
func Transfer(l *Ledger, from, to string, amount int64) error {
	src, err := account(l, from)
	if err != nil {
		return err
	}
	dst, err := account(l, to)
	if err != nil {
		return err
	}
	switch {
	case from == to:
		return fmt.Errorf("%w: %q", ErrSameAccount, from)
	case amount <= 0:
		return fmt.Errorf("%w: %d", ErrInvalidAmount, amount)
	case amount > src.Balance:
		return fmt.Errorf("%w: %q has %d, need %d", ErrInsufficient, from, src.Balance, amount)
	}
	// All checks passed: apply both sides together.
	src.Balance -= amount
	dst.Balance += amount
	l.History = append(l.History, Tx{From: from, To: to, Amount: amount})
	return nil
}

// Balance returns the balance of account id in cents.
func Balance(l *Ledger, id string) (int64, error) {
	a, err := account(l, id)
	if err != nil {
		return 0, err
	}
	return a.Balance, nil
}

// Statement returns the transactions involving account id, oldest first.
func Statement(l *Ledger, id string) ([]Tx, error) {
	if _, err := account(l, id); err != nil {
		return nil, err
	}
	out := []Tx{} // a fresh slice, never a sub-slice of l.History
	for _, tx := range l.History {
		if tx.From == id || tx.To == id {
			out = append(out, tx)
		}
	}
	return out, nil
}

// Total returns the sum of all account balances.
func Total(l *Ledger) int64 {
	var sum int64
	for _, a := range l.Accounts {
		sum += a.Balance
	}
	return sum
}

// FormatCents formats a number of cents as units with two decimals.
func FormatCents(c int64) string {
	sign := ""
	u := uint64(c)
	if c < 0 {
		sign = "-"
		u = -u // unsigned negation is exact even for math.MinInt64
	}
	return fmt.Sprintf("%s%d.%02d", sign, u/100, u%100)
}
