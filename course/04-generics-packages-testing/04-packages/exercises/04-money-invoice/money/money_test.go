package money_test

import (
	"errors"
	"testing"

	"learngo/course/04-generics-packages-testing/04-packages/exercises/04-money-invoice/money"
)

func TestNew(t *testing.T) {
	tests := []struct {
		units, cents int64
		want         string
		wantCents    int64
		wantErr      error
	}{
		{12, 5, "12.05", 1205, nil},
		{0, 0, "0.00", 0, nil},
		{0, 99, "0.99", 99, nil},
		{100, 0, "100.00", 10000, nil},
		{1, 10, "1.10", 110, nil},
		{-1, 0, "0.00", 0, money.ErrNegative},
		{1, -1, "0.00", 0, money.ErrNegative},
		{-1, 100, "0.00", 0, money.ErrNegative},
		{3, 100, "0.00", 0, money.ErrCents},
		{0, 250, "0.00", 0, money.ErrCents},
	}
	for _, tt := range tests {
		got, err := money.New(tt.units, tt.cents)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("New(%d, %d) error = %v, want %v", tt.units, tt.cents, err, tt.wantErr)
			continue
		}
		if got.String() != tt.want || got.Cents() != tt.wantCents {
			t.Errorf("New(%d, %d) = %q (%d cents), want %q (%d cents)", tt.units, tt.cents, got.String(), got.Cents(), tt.want, tt.wantCents)
		}
	}
}

func TestFromCents(t *testing.T) {
	tests := []struct {
		in      int64
		want    string
		wantErr error
	}{
		{1999, "19.99", nil},
		{5, "0.05", nil},
		{0, "0.00", nil},
		{100, "1.00", nil},
		{-1, "0.00", money.ErrNegative},
	}
	for _, tt := range tests {
		got, err := money.FromCents(tt.in)
		if !errors.Is(err, tt.wantErr) || got.String() != tt.want {
			t.Errorf("FromCents(%d) = %q, %v, want %q, %v", tt.in, got.String(), err, tt.want, tt.wantErr)
		}
	}
}

func cents(t *testing.T, c int64) money.Amount {
	t.Helper()
	a, err := money.FromCents(c)
	if err != nil {
		t.Fatalf("FromCents(%d) error = %v, want <nil>", c, err)
	}
	return a
}

func TestZeroValue(t *testing.T) {
	var a money.Amount
	if got := a.String(); got != "0.00" {
		t.Errorf("Amount{}.String() = %q, want %q", got, "0.00")
	}
	if got := a.Cents(); got != 0 {
		t.Errorf("Amount{}.Cents() = %d, want 0", got)
	}
}

func TestArithmetic(t *testing.T) {
	p := cents(t, 150)
	tests := []struct {
		name string
		got  money.Amount
		want string
	}{
		{"1.50.Add(0.99)", p.Add(cents(t, 99)), "2.49"},
		{"1.50.Add(0.00)", p.Add(money.Amount{}), "1.50"},
		{"0.50.Add(0.50)", cents(t, 50).Add(cents(t, 50)), "1.00"},
		{"1.50.Times(3)", p.Times(3), "4.50"},
		{"1.50.Times(0)", p.Times(0), "0.00"},
		{"1.50.Times(-2)", p.Times(-2), "0.00"},
		{"0.01.Times(100)", cents(t, 1).Times(100), "1.00"},
	}
	for _, tt := range tests {
		if got := tt.got.String(); got != tt.want {
			t.Errorf("%s = %q, want %q", tt.name, got, tt.want)
		}
	}
}
