package bank_test

import (
	"testing"

	bank "learngo/course/03-methods-interfaces-errors/01-methods/exercises/04-bank-account"
)

func newAcc(t *testing.T, owner string, initial int64) *bank.Account {
	t.Helper()
	a := bank.NewAccount(owner, initial)
	if a == nil {
		t.Fatalf("NewAccount(%q, %d) = nil, want an account", owner, initial)
	}
	return a
}

func TestNewAccount(t *testing.T) {
	tests := []struct {
		owner   string
		initial int64
		want    int64
	}{
		{"ann", 100, 100},
		{"bob", 0, 0},
		{"cy", -5, 0},
		{"Zoë", 123456789, 123456789},
	}
	for _, tt := range tests {
		a := newAcc(t, tt.owner, tt.initial)
		if got := a.Balance(); got != tt.want {
			t.Errorf("NewAccount(%q, %d).Balance() = %d, want %d", tt.owner, tt.initial, got, tt.want)
		}
		if got := a.Owner(); got != tt.owner {
			t.Errorf("NewAccount(%q, %d).Owner() = %q, want %q", tt.owner, tt.initial, got, tt.owner)
		}
	}
}

func TestSetOwner(t *testing.T) {
	a := newAcc(t, "ann", 0)
	if !a.SetOwner("anna") {
		t.Errorf("SetOwner(%q) = false, want true", "anna")
	}
	if got := a.Owner(); got != "anna" {
		t.Errorf("after SetOwner(%q) Owner() = %q, want %q", "anna", got, "anna")
	}
	if a.SetOwner("") {
		t.Errorf("SetOwner(%q) = true, want false", "")
	}
	if got := a.Owner(); got != "anna" {
		t.Errorf("after SetOwner(%q) Owner() = %q, want %q", "", got, "anna")
	}
}

func TestDeposit(t *testing.T) {
	tests := []struct {
		start, amount int64
		ok            bool
		want          int64
	}{
		{100, 50, true, 150},
		{0, 1, true, 1},
		{100, 0, false, 100},
		{100, -20, false, 100},
	}
	for _, tt := range tests {
		a := newAcc(t, "ann", tt.start)
		ok := a.Deposit(tt.amount)
		if ok != tt.ok {
			t.Errorf("Deposit(%d) on balance %d = %v, want %v", tt.amount, tt.start, ok, tt.ok)
		}
		if got := a.Balance(); got != tt.want {
			t.Errorf("Deposit(%d) on balance %d: Balance() = %d, want %d", tt.amount, tt.start, got, tt.want)
		}
	}
}

func TestWithdraw(t *testing.T) {
	tests := []struct {
		start, amount int64
		ok            bool
		want          int64
	}{
		{150, 30, true, 120},
		{100, 100, true, 0},
		{100, 101, false, 100},
		{100, 500, false, 100},
		{0, 1, false, 0},
		{100, 0, false, 100},
		{100, -5, false, 100},
	}
	for _, tt := range tests {
		a := newAcc(t, "ann", tt.start)
		ok := a.Withdraw(tt.amount)
		if ok != tt.ok {
			t.Errorf("Withdraw(%d) on balance %d = %v, want %v", tt.amount, tt.start, ok, tt.ok)
		}
		if got := a.Balance(); got != tt.want {
			t.Errorf("Withdraw(%d) on balance %d: Balance() = %d, want %d", tt.amount, tt.start, got, tt.want)
		}
	}
}

func TestTransfer(t *testing.T) {
	tests := []struct {
		desc             string
		from, to, amount int64
		ok               bool
		wantFrom, wantTo int64
	}{
		{"normal", 100, 10, 20, true, 80, 30},
		{"everything", 100, 0, 100, true, 0, 100},
		{"too much", 100, 10, 101, false, 100, 10},
		{"zero", 100, 10, 0, false, 100, 10},
		{"negative", 100, 10, -50, false, 100, 10},
	}
	for _, tt := range tests {
		a := newAcc(t, "ann", tt.from)
		b := newAcc(t, "bob", tt.to)
		ok := a.Transfer(b, tt.amount)
		if ok != tt.ok {
			t.Errorf("%s: Transfer(%d) from %d to %d = %v, want %v", tt.desc, tt.amount, tt.from, tt.to, ok, tt.ok)
		}
		if a.Balance() != tt.wantFrom || b.Balance() != tt.wantTo {
			t.Errorf("%s: Transfer(%d) balances = (%d, %d), want (%d, %d)",
				tt.desc, tt.amount, a.Balance(), b.Balance(), tt.wantFrom, tt.wantTo)
		}
	}
}

func TestTransferInvalidTarget(t *testing.T) {
	a := newAcc(t, "ann", 100)
	if a.Transfer(nil, 10) {
		t.Errorf("Transfer(nil, 10) = true, want false")
	}
	if a.Transfer(a, 10) {
		t.Errorf("Transfer(self, 10) = true, want false")
	}
	if got := a.Balance(); got != 100 {
		t.Errorf("after invalid transfers Balance() = %d, want 100", got)
	}
}

func TestSharedPointer(t *testing.T) {
	a := newAcc(t, "ann", 0)
	alias := a
	alias.Deposit(70)
	if got := a.Balance(); got != 70 {
		t.Errorf("Deposit(70) through a copied *Account: Balance() = %d, want 70", got)
	}
}
