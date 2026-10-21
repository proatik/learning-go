package bank

import "errors"

// Sentinel errors.
var (
	ErrInvalidAmount = errors.New("bank: invalid amount")
	ErrAccountFrozen = errors.New("bank: account frozen")
	ErrNoAccount     = errors.New("bank: no such account")
	ErrAccountExists = errors.New("bank: account already exists")
	ErrSameAccount   = errors.New("bank: same account")
	ErrInvalidID     = errors.New("bank: invalid account id")
)

// Money is an amount in cents.
type Money int64

// String formats cents as "12.34", "-0.05", "0.00".
func (m Money) String() string {
	// TODO
	return ""
}

// InsufficientFundsError reports a withdrawal larger than the balance.
type InsufficientFundsError struct {
	ID      string
	Balance Money
	Amount  Money
}

// Error returns "bank: ID has BALANCE, needs AMOUNT".
func (e *InsufficientFundsError) Error() string {
	// TODO
	return ""
}

// Entry is one line of an account statement.
type Entry struct {
	Desc    string
	Amount  Money // positive = in, negative = out
	Balance Money // balance after the entry
}

// String returns "DESC +AMOUNT = BALANCE" (or "-AMOUNT").
func (e Entry) String() string {
	// TODO
	return ""
}

// Account is a bank account. Create it with Bank.Open.
type Account struct {
	id, owner string
	balance   Money
	frozen    bool
	entries   []Entry
}

// ID returns the account ID.
func (a *Account) ID() string {
	// TODO
	return ""
}

// Owner returns the account owner.
func (a *Account) Owner() string {
	// TODO
	return ""
}

// Balance returns the current balance.
func (a *Account) Balance() Money {
	// TODO
	return 0
}

// Frozen reports whether the account is frozen.
func (a *Account) Frozen() bool {
	// TODO
	return false
}

// Freeze blocks all deposits and withdrawals.
func (a *Account) Freeze() {
	// TODO
}

// Unfreeze allows deposits and withdrawals again.
func (a *Account) Unfreeze() {
	// TODO
}

// Entries returns a copy of the statement entries, oldest first.
func (a *Account) Entries() []Entry {
	// TODO
	return nil
}

// Deposit adds amt to the account.
func (a *Account) Deposit(amt Money) error {
	// TODO
	return nil
}

// Withdraw removes amt from the account.
func (a *Account) Withdraw(amt Money) error {
	// TODO
	return nil
}

// String returns the statement: a header line plus one indented line per entry.
func (a *Account) String() string {
	// TODO
	return ""
}

// Bank holds accounts by ID. Create it with NewBank.
type Bank struct {
	accounts map[string]*Account
}

// NewBank returns an empty bank.
func NewBank() *Bank {
	// TODO
	return nil
}

// Open creates a new account with a zero balance.
func (b *Bank) Open(id, owner string) (*Account, error) {
	// TODO
	return nil, nil
}

// Account returns the account with the given ID.
func (b *Bank) Account(id string) (*Account, error) {
	// TODO
	return nil, nil
}

// Transfer moves amt from one account to another. If any step fails, both
// accounts are left exactly as they were.
func (b *Bank) Transfer(from, to string, amt Money) (err error) {
	// TODO
	return nil
}
