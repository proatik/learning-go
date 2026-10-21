package bank

import "errors"

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
	// TODO
	return nil
}

// Open creates account id with an initial balance in cents.
func Open(l *Ledger, id, owner string, initial int64) error {
	// TODO
	return nil
}

// Deposit adds amount cents to account id.
func Deposit(l *Ledger, id string, amount int64) error {
	// TODO
	return nil
}

// Withdraw removes amount cents from account id.
func Withdraw(l *Ledger, id string, amount int64) error {
	// TODO
	return nil
}

// Transfer moves amount cents from one account to another.
func Transfer(l *Ledger, from, to string, amount int64) error {
	// TODO
	return nil
}

// Balance returns the balance of account id in cents.
func Balance(l *Ledger, id string) (int64, error) {
	// TODO
	return 0, nil
}

// Statement returns the transactions involving account id, oldest first.
func Statement(l *Ledger, id string) ([]Tx, error) {
	// TODO
	return nil, nil
}

// Total returns the sum of all account balances.
func Total(l *Ledger) int64 {
	// TODO
	return 0
}

// FormatCents formats a number of cents as units with two decimals.
func FormatCents(c int64) string {
	// TODO
	return ""
}
