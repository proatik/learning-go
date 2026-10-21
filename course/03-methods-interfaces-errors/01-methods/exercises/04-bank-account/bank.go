package bank

// Account is a bank account whose balance (in cents) never goes negative.
// Create one with NewAccount.
type Account struct {
	owner   string
	balance int64
}

// NewAccount returns an account for owner holding initial cents
// (a negative initial amount counts as zero).
func NewAccount(owner string, initial int64) *Account {
	// TODO
	return &Account{}
}

// Owner returns the account owner's name.
func (a *Account) Owner() string {
	// TODO
	return ""
}

// Balance returns the balance in cents.
func (a *Account) Balance() int64 {
	// TODO
	return 0
}

// SetOwner renames the owner and reports whether name was accepted
// (it must not be empty).
func (a *Account) SetOwner(name string) bool {
	// TODO
	return false
}

// Deposit adds amount and reports whether it was accepted (amount > 0).
func (a *Account) Deposit(amount int64) bool {
	// TODO
	return false
}

// Withdraw removes amount and reports whether it succeeded
// (0 < amount <= balance).
func (a *Account) Withdraw(amount int64) bool {
	// TODO
	return false
}

// Transfer moves amount from a to another account and reports whether it
// succeeded. Nothing changes on failure.
func (a *Account) Transfer(to *Account, amount int64) bool {
	// TODO
	return false
}
