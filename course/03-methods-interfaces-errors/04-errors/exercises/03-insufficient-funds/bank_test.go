package bank

import (
	"errors"
	"testing"
)

func TestInsufficientFundsError(t *testing.T) {
	tests := []struct {
		e         InsufficientFundsError
		msg       string
		shortfall int
	}{
		{InsufficientFundsError{Need: 150, Have: 100}, "insufficient funds: need 150, have 100", 50},
		{InsufficientFundsError{Need: 1, Have: 0}, "insufficient funds: need 1, have 0", 1},
		{InsufficientFundsError{Need: 10, Have: -20}, "insufficient funds: need 10, have -20", 30},
	}
	for _, tt := range tests {
		e := tt.e
		if got := e.Error(); got != tt.msg {
			t.Errorf("(%+v).Error() = %q, want %q", e, got, tt.msg)
		}
		if got := e.Shortfall(); got != tt.shortfall {
			t.Errorf("(%+v).Shortfall() = %d, want %d", e, got, tt.shortfall)
		}
	}
}

func TestWithdrawOK(t *testing.T) {
	tests := []struct {
		balance, amount, want int
	}{
		{100, 30, 70},
		{100, 100, 0},
		{100, 1, 99},
	}
	for _, tt := range tests {
		a := &Account{Balance: tt.balance}
		if err := a.Withdraw(tt.amount); err != nil {
			t.Errorf("Withdraw(%d) on balance %d error = %v, want <nil>", tt.amount, tt.balance, err)
		}
		if a.Balance != tt.want {
			t.Errorf("Withdraw(%d) on balance %d: balance = %d, want %d", tt.amount, tt.balance, a.Balance, tt.want)
		}
	}
}

func TestWithdrawInvalidAmount(t *testing.T) {
	tests := []struct {
		amount int
		msg    string
	}{
		{0, "withdraw 0: invalid amount"},
		{-5, "withdraw -5: invalid amount"},
	}
	for _, tt := range tests {
		a := &Account{Balance: 100}
		err := a.Withdraw(tt.amount)
		if !errors.Is(err, ErrInvalidAmount) {
			t.Errorf("Withdraw(%d) error = %v, want %v", tt.amount, err, ErrInvalidAmount)
			continue
		}
		if err.Error() != tt.msg {
			t.Errorf("Withdraw(%d) error = %q, want %q", tt.amount, err.Error(), tt.msg)
		}
		if a.Balance != 100 {
			t.Errorf("Withdraw(%d): balance = %d, want 100 (unchanged)", tt.amount, a.Balance)
		}
	}
}

func TestWithdrawInsufficient(t *testing.T) {
	tests := []struct {
		balance, amount int
	}{
		{100, 150},
		{0, 1},
		{100, 101},
	}
	for _, tt := range tests {
		a := &Account{Balance: tt.balance}
		err := a.Withdraw(tt.amount)
		want := InsufficientFundsError{Need: tt.amount, Have: tt.balance}
		fe, ok := errors.AsType[*InsufficientFundsError](err)
		if !ok {
			t.Errorf("Withdraw(%d) on balance %d error = %v (%T), want %v", tt.amount, tt.balance, err, err, &want)
			continue
		}
		if *fe != want {
			t.Errorf("Withdraw(%d) on balance %d error = %+v, want %+v", tt.amount, tt.balance, *fe, want)
		}
		if errors.Is(err, ErrInvalidAmount) {
			t.Errorf("Withdraw(%d) error = %v, should not match ErrInvalidAmount", tt.amount, err)
		}
		if a.Balance != tt.balance {
			t.Errorf("Withdraw(%d): balance = %d, want %d (unchanged)", tt.amount, a.Balance, tt.balance)
		}
	}
}

func TestTransfer(t *testing.T) {
	a, b := &Account{Balance: 100}, &Account{Balance: 5}
	if err := a.Transfer(b, 40); err != nil {
		t.Fatalf("Transfer(40) error = %v, want <nil>", err)
	}
	if a.Balance != 60 || b.Balance != 45 {
		t.Errorf("Transfer(40): balances = %d, %d, want 60, 45", a.Balance, b.Balance)
	}

	err := a.Transfer(b, 150)
	var fe *InsufficientFundsError
	if !errors.As(err, &fe) {
		t.Fatalf("Transfer(150) error = %v, want a wrapped *InsufficientFundsError", err)
	}
	if got, want := err.Error(), "transfer: insufficient funds: need 150, have 60"; got != want {
		t.Errorf("Transfer(150) error = %q, want %q", got, want)
	}
	if fe.Shortfall() != 90 {
		t.Errorf("Transfer(150) shortfall = %d, want 90", fe.Shortfall())
	}
	if a.Balance != 60 || b.Balance != 45 {
		t.Errorf("Transfer(150): balances = %d, %d, want 60, 45 (unchanged)", a.Balance, b.Balance)
	}

	err = a.Transfer(b, -1)
	if !errors.Is(err, ErrInvalidAmount) {
		t.Errorf("Transfer(-1) error = %v, want %v", err, ErrInvalidAmount)
	} else if got, want := err.Error(), "transfer: withdraw -1: invalid amount"; got != want {
		t.Errorf("Transfer(-1) error = %q, want %q", got, want)
	}
}
