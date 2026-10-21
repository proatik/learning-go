package bank

import (
	"errors"
	"fmt"
	"strings"
)

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
	sign := ""
	u := uint64(m)
	if m < 0 {
		sign, u = "-", -u
	}
	return fmt.Sprintf("%s%d.%02d", sign, u/100, u%100)
}

// InsufficientFundsError reports a withdrawal larger than the balance.
type InsufficientFundsError struct {
	ID      string
	Balance Money
	Amount  Money
}

// Error returns "bank: ID has BALANCE, needs AMOUNT".
func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("bank: %s has %v, needs %v", e.ID, e.Balance, e.Amount)
}

// Entry is one line of an account statement.
type Entry struct {
	Desc    string
	Amount  Money // positive = in, negative = out
	Balance Money // balance after the entry
}

// String returns "DESC +AMOUNT = BALANCE" (or "-AMOUNT").
func (e Entry) String() string {
	sign := "+"
	if e.Amount < 0 {
		sign = ""
	}
	return fmt.Sprintf("%s %s%v = %v", e.Desc, sign, e.Amount, e.Balance)
}

// Account is a bank account. Create it with Bank.Open.
type Account struct {
	id, owner string
	balance   Money
	frozen    bool
	entries   []Entry
}

// ID returns the account ID.
func (a *Account) ID() string { return a.id }

// Owner returns the account owner.
func (a *Account) Owner() string { return a.owner }

// Balance returns the current balance.
func (a *Account) Balance() Money { return a.balance }

// Frozen reports whether the account is frozen.
func (a *Account) Frozen() bool { return a.frozen }

// Freeze blocks all deposits and withdrawals.
func (a *Account) Freeze() { a.frozen = true }

// Unfreeze allows deposits and withdrawals again.
func (a *Account) Unfreeze() { a.frozen = false }

// Entries returns a copy of the statement entries, oldest first.
func (a *Account) Entries() []Entry { return append([]Entry{}, a.entries...) }

// credit adds amt and records desc.
func (a *Account) credit(amt Money, desc string) error {
	if amt <= 0 {
		return ErrInvalidAmount
	}
	if a.frozen {
		return fmt.Errorf("%s: %w", a.id, ErrAccountFrozen)
	}
	a.balance += amt
	a.entries = append(a.entries, Entry{Desc: desc, Amount: amt, Balance: a.balance})
	return nil
}

// debit removes amt and records desc.
func (a *Account) debit(amt Money, desc string) error {
	if amt <= 0 {
		return ErrInvalidAmount
	}
	if a.frozen {
		return fmt.Errorf("%s: %w", a.id, ErrAccountFrozen)
	}
	if amt > a.balance {
		return &InsufficientFundsError{ID: a.id, Balance: a.balance, Amount: amt}
	}
	a.balance -= amt
	a.entries = append(a.entries, Entry{Desc: desc, Amount: -amt, Balance: a.balance})
	return nil
}

// Deposit adds amt to the account.
func (a *Account) Deposit(amt Money) error { return a.credit(amt, "deposit") }

// Withdraw removes amt from the account.
func (a *Account) Withdraw(amt Money) error { return a.debit(amt, "withdraw") }

// String returns the statement: a header line plus one indented line per entry.
func (a *Account) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s (%s): %v", a.id, a.owner, a.balance)
	for _, e := range a.entries {
		fmt.Fprintf(&b, "\n  %v", e)
	}
	return b.String()
}

// Bank holds accounts by ID. Create it with NewBank.
type Bank struct {
	accounts map[string]*Account
}

// NewBank returns an empty bank.
func NewBank() *Bank {
	return &Bank{accounts: make(map[string]*Account)}
}

// Open creates a new account with a zero balance.
func (b *Bank) Open(id, owner string) (*Account, error) {
	if id == "" {
		return nil, ErrInvalidID
	}
	if _, ok := b.accounts[id]; ok {
		return nil, fmt.Errorf("open %s: %w", id, ErrAccountExists)
	}
	a := &Account{id: id, owner: owner}
	b.accounts[id] = a
	return a, nil
}

// Account returns the account with the given ID.
func (b *Bank) Account(id string) (*Account, error) {
	a, ok := b.accounts[id]
	if !ok {
		return nil, fmt.Errorf("%s: %w", id, ErrNoAccount)
	}
	return a, nil
}

// Transfer moves amt from one account to another. If any step fails, both
// accounts are left exactly as they were.
func (b *Bank) Transfer(from, to string, amt Money) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("transfer %s -> %s: %w", from, to, err)
		}
	}()
	src, err := b.Account(from)
	if err != nil {
		return err
	}
	dst, err := b.Account(to)
	if err != nil {
		return err
	}
	if src == dst {
		return ErrSameAccount
	}

	if err := src.debit(amt, "transfer to "+to); err != nil {
		return err
	}
	// Roll the debit back if anything after this point fails.
	saved := len(src.entries) - 1
	defer func() {
		if err != nil {
			src.balance += amt
			src.entries = src.entries[:saved]
		}
	}()
	return dst.credit(amt, "transfer from "+from)
}
