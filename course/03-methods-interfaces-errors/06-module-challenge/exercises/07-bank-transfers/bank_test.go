package bank_test

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"testing"

	bank "learngo/course/03-methods-interfaces-errors/06-module-challenge/exercises/07-bank-transfers"
)

func newBank(t *testing.T) *bank.Bank {
	t.Helper()
	b := bank.NewBank()
	if b == nil {
		t.Fatalf("NewBank() = nil, want a bank")
	}
	return b
}

func open(t *testing.T, b *bank.Bank, id, owner string, deposit bank.Money) *bank.Account {
	t.Helper()
	a, err := b.Open(id, owner)
	if err != nil || a == nil {
		t.Fatalf("Open(%q, %q) = %v, %v, want an account, nil", id, owner, a, err)
	}
	if deposit > 0 {
		if err := a.Deposit(deposit); err != nil {
			t.Fatalf("Deposit(%v) error = %v, want nil", deposit, err)
		}
	}
	return a
}

// setup opens a1 (Ada, 100.00) and b2 (Bob, 0).
func setup(t *testing.T) (*bank.Bank, *bank.Account, *bank.Account) {
	t.Helper()
	b := newBank(t)
	return b, open(t, b, "a1", "Ada", 10_000), open(t, b, "b2", "Bob", 0)
}

func TestMoneyString(t *testing.T) {
	tests := []struct {
		m    bank.Money
		want string
	}{
		{0, "0.00"},
		{5, "0.05"},
		{99, "0.99"},
		{100, "1.00"},
		{1234, "12.34"},
		{-5, "-0.05"},
		{-150, "-1.50"},
		{math.MaxInt64, "92233720368547758.07"},
		{math.MinInt64, "-92233720368547758.08"},
	}
	for _, tt := range tests {
		if got := tt.m.String(); got != tt.want {
			t.Errorf("Money(%d).String() = %q, want %q", int64(tt.m), got, tt.want)
		}
	}
	if got := fmt.Sprintf("%v|%s", bank.Money(250), bank.Money(-1)); got != "2.50|-0.01" {
		t.Errorf("Sprintf(%%v|%%s) = %q, want %q", got, "2.50|-0.01")
	}
}

func TestOpen(t *testing.T) {
	b := newBank(t)
	a := open(t, b, "a1", "Ada", 0)
	if a.ID() != "a1" || a.Owner() != "Ada" || a.Balance() != 0 || a.Frozen() {
		t.Errorf("new account = %q %q %v frozen=%v, want a1 Ada 0.00 frozen=false", a.ID(), a.Owner(), a.Balance(), a.Frozen())
	}
	if got, err := b.Account("a1"); got != a || err != nil {
		t.Errorf("Account(%q) = %p, %v, want the same *Account, nil", "a1", got, err)
	}
	if _, err := b.Open("a1", "Eve"); !errors.Is(err, bank.ErrAccountExists) {
		t.Errorf("Open(%q) twice error = %v, want %v", "a1", err, bank.ErrAccountExists)
	}
	if acc, err := b.Open("", "Nobody"); !errors.Is(err, bank.ErrInvalidID) || acc != nil {
		t.Errorf("Open(\"\") = %v, %v, want nil, %v", acc, err, bank.ErrInvalidID)
	}
	if acc, err := b.Account("zz"); !errors.Is(err, bank.ErrNoAccount) || acc != nil {
		t.Errorf("Account(%q) = %v, %v, want nil, %v", "zz", acc, err, bank.ErrNoAccount)
	}
}

func TestDepositWithdraw(t *testing.T) {
	b := newBank(t)
	a := open(t, b, "a1", "Ada", 0)
	steps := []struct {
		name string
		do   func(bank.Money) error
		amt  bank.Money
		want error // nil, or a sentinel for errors.Is
		bal  bank.Money
	}{
		{"Deposit", a.Deposit, 10_000, nil, 10_000},
		{"Withdraw", a.Withdraw, 2_550, nil, 7_450},
		{"Deposit", a.Deposit, 0, bank.ErrInvalidAmount, 7_450},
		{"Deposit", a.Deposit, -5, bank.ErrInvalidAmount, 7_450},
		{"Withdraw", a.Withdraw, 0, bank.ErrInvalidAmount, 7_450},
		{"Withdraw", a.Withdraw, -1, bank.ErrInvalidAmount, 7_450},
		{"Withdraw", a.Withdraw, 7_450, nil, 0}, // exactly the balance is fine
	}
	for _, s := range steps {
		err := s.do(s.amt)
		if s.want == nil && err != nil || s.want != nil && !errors.Is(err, s.want) {
			t.Fatalf("%s(%v) error = %v, want %v", s.name, s.amt, err, s.want)
		}
		if a.Balance() != s.bal {
			t.Fatalf("after %s(%v): Balance() = %v, want %v", s.name, s.amt, a.Balance(), s.bal)
		}
	}
	if got := len(a.Entries()); got != 3 {
		t.Errorf("len(Entries()) = %d, want 3 (failed operations record nothing)", got)
	}
}

func TestInsufficientFunds(t *testing.T) {
	_, a, _ := setup(t)
	err := a.Withdraw(10_001)
	ife, ok := errors.AsType[*bank.InsufficientFundsError](err)
	if !ok {
		t.Fatalf("Withdraw(100.01) error = %v (%T), want *InsufficientFundsError", err, err)
	}
	if *ife != (bank.InsufficientFundsError{ID: "a1", Balance: 10_000, Amount: 10_001}) {
		t.Errorf("Withdraw(100.01) error = %+v, want {ID:a1 Balance:10000 Amount:10001}", *ife)
	}
	if want := "bank: a1 has 100.00, needs 100.01"; err.Error() != want {
		t.Errorf("Withdraw(100.01) error = %q, want %q", err.Error(), want)
	}
	if a.Balance() != 10_000 {
		t.Errorf("Balance() = %v after a failed Withdraw, want 100.00", a.Balance())
	}
}

func TestFrozen(t *testing.T) {
	_, a, _ := setup(t)
	a.Freeze()
	if !a.Frozen() {
		t.Fatalf("Frozen() = false after Freeze()")
	}
	if err := a.Deposit(1); !errors.Is(err, bank.ErrAccountFrozen) {
		t.Errorf("Deposit on a frozen account error = %v, want %v", err, bank.ErrAccountFrozen)
	}
	err := a.Withdraw(1)
	if !errors.Is(err, bank.ErrAccountFrozen) {
		t.Errorf("Withdraw on a frozen account error = %v, want %v", err, bank.ErrAccountFrozen)
	} else if err.Error() != "a1: bank: account frozen" {
		t.Errorf("Withdraw on a frozen account error = %q, want %q", err.Error(), "a1: bank: account frozen")
	}
	// Amount is checked before frozen; frozen before funds.
	if err := a.Withdraw(0); !errors.Is(err, bank.ErrInvalidAmount) {
		t.Errorf("Withdraw(0) on a frozen account error = %v, want %v", err, bank.ErrInvalidAmount)
	}
	if err := a.Withdraw(1_000_000); !errors.Is(err, bank.ErrAccountFrozen) {
		t.Errorf("Withdraw(10000.00) on a frozen account error = %v, want %v", err, bank.ErrAccountFrozen)
	}
	a.Unfreeze()
	if err := a.Withdraw(1); err != nil || a.Balance() != 9_999 {
		t.Errorf("Withdraw(0.01) after Unfreeze = %v, balance %v, want nil, 99.99", err, a.Balance())
	}
}

func TestTransfer(t *testing.T) {
	b, a, c := setup(t)
	if err := b.Transfer("a1", "b2", 2_550); err != nil {
		t.Fatalf("Transfer(a1, b2, 25.50) error = %v, want nil", err)
	}
	if a.Balance() != 7_450 || c.Balance() != 2_550 {
		t.Errorf("balances = %v, %v, want 74.50, 25.50", a.Balance(), c.Balance())
	}
	wantA := []bank.Entry{{"deposit", 10_000, 10_000}, {"transfer to b2", -2_550, 7_450}}
	if got := a.Entries(); !slices.Equal(got, wantA) {
		t.Errorf("a1 Entries() = %v, want %v", got, wantA)
	}
	wantB := []bank.Entry{{"transfer from a1", 2_550, 2_550}}
	if got := c.Entries(); !slices.Equal(got, wantB) {
		t.Errorf("b2 Entries() = %v, want %v", got, wantB)
	}
}

func TestTransferErrors(t *testing.T) {
	tests := []struct {
		from, to string
		amt      bank.Money
		want     error
	}{
		{"zz", "b2", 100, bank.ErrNoAccount},
		{"a1", "zz", 100, bank.ErrNoAccount},
		{"a1", "a1", 100, bank.ErrSameAccount},
		{"a1", "b2", 0, bank.ErrInvalidAmount},
		{"a1", "b2", -100, bank.ErrInvalidAmount},
	}
	for _, tt := range tests {
		b, a, c := setup(t)
		err := b.Transfer(tt.from, tt.to, tt.amt)
		if !errors.Is(err, tt.want) {
			t.Errorf("Transfer(%q, %q, %v) error = %v, want %v", tt.from, tt.to, tt.amt, err, tt.want)
		}
		if a.Balance() != 10_000 || c.Balance() != 0 || len(a.Entries()) != 1 || len(c.Entries()) != 0 {
			t.Errorf("Transfer(%q, %q, %v) failed but changed the accounts", tt.from, tt.to, tt.amt)
		}
	}

	b, _, _ := setup(t)
	err := b.Transfer("b2", "a1", 1)
	if want := "transfer b2 -> a1: bank: b2 has 0.00, needs 0.01"; err == nil || err.Error() != want {
		t.Errorf("Transfer(b2, a1, 0.01) error = %v, want %q", err, want)
	}
	if ife, ok := errors.AsType[*bank.InsufficientFundsError](err); !ok || ife.ID != "b2" {
		t.Errorf("errors.AsType[*InsufficientFundsError](%v) = %v, %v, want b2's error, true", err, ife, ok)
	}
}

func TestTransferRollsBack(t *testing.T) {
	b, a, c := setup(t)
	c.Freeze() // the debit succeeds, the credit fails → undo the debit
	err := b.Transfer("a1", "b2", 3_000)
	if !errors.Is(err, bank.ErrAccountFrozen) {
		t.Fatalf("Transfer to a frozen account error = %v, want %v", err, bank.ErrAccountFrozen)
	}
	if want := "transfer a1 -> b2: b2: bank: account frozen"; err.Error() != want {
		t.Errorf("Transfer to a frozen account error = %q, want %q", err.Error(), want)
	}
	if a.Balance() != 10_000 {
		t.Errorf("a1 Balance() = %v after a failed transfer, want 100.00 (roll back the debit)", a.Balance())
	}
	if want := []bank.Entry{{"deposit", 10_000, 10_000}}; !slices.Equal(a.Entries(), want) {
		t.Errorf("a1 Entries() = %v, want %v (remove the rolled-back entry)", a.Entries(), want)
	}
	if c.Balance() != 0 || len(c.Entries()) != 0 {
		t.Errorf("b2 changed during a failed transfer: %v, %v", c.Balance(), c.Entries())
	}

	// A frozen source fails before anything happens.
	c.Unfreeze()
	a.Freeze()
	if err := b.Transfer("a1", "b2", 1); !errors.Is(err, bank.ErrAccountFrozen) {
		t.Errorf("Transfer from a frozen account error = %v, want %v", err, bank.ErrAccountFrozen)
	}
	a.Unfreeze()
	// And the accounts still work normally afterwards.
	if err := b.Transfer("a1", "b2", 1_000); err != nil || a.Balance() != 9_000 || c.Balance() != 1_000 {
		t.Errorf("Transfer after rollback = %v, balances %v %v, want nil, 90.00 10.00", err, a.Balance(), c.Balance())
	}
}

func TestEntriesIsACopy(t *testing.T) {
	_, a, _ := setup(t)
	e := a.Entries()
	if len(e) != 1 {
		t.Fatalf("Entries() = %v, want one deposit", e)
	}
	e[0].Amount = 1
	if a.Entries()[0].Amount != 10_000 {
		t.Errorf("changing the Entries() result changed the account")
	}
}

func TestStatement(t *testing.T) {
	b, a, c := setup(t)
	if err := b.Transfer("a1", "b2", 2_550); err != nil {
		t.Fatalf("Transfer error = %v", err)
	}
	if err := c.Withdraw(50); err != nil {
		t.Fatalf("Withdraw error = %v", err)
	}
	wantA := "a1 (Ada): 74.50\n  deposit +100.00 = 100.00\n  transfer to b2 -25.50 = 74.50"
	if got := a.String(); got != wantA {
		t.Errorf("a1.String() =\n%s\nwant\n%s", got, wantA)
	}
	wantB := "b2 (Bob): 25.00\n  transfer from a1 +25.50 = 25.50\n  withdraw -0.50 = 25.00"
	if got := fmt.Sprint(c); got != wantB {
		t.Errorf("fmt.Sprint(b2) =\n%s\nwant\n%s", got, wantB)
	}
	empty := open(t, b, "e", "Eve", 0)
	if got := empty.String(); got != "e (Eve): 0.00" {
		t.Errorf("empty.String() = %q, want %q", got, "e (Eve): 0.00")
	}
	if got := fmt.Sprint(bank.Entry{Desc: "deposit", Amount: 5, Balance: 5}); got != "deposit +0.05 = 0.05" {
		t.Errorf("Entry.String() = %q, want %q", got, "deposit +0.05 = 0.05")
	}
}
