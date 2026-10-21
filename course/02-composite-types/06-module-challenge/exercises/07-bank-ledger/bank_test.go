package bank

import (
	"errors"
	"math"
	"slices"
	"testing"
)

func newLedger(t *testing.T) *Ledger {
	t.Helper()
	l := NewLedger()
	if l == nil {
		t.Fatalf("NewLedger() = nil, want a ledger")
	}
	if l.Accounts == nil {
		t.Fatalf("NewLedger().Accounts = nil, want an empty map")
	}
	if len(l.History) != 0 {
		t.Fatalf("NewLedger().History = %v, want empty", l.History)
	}
	return l
}

// setup opens a1 (Ada, 100.00) and b2 (Bob, 0).
func setup(t *testing.T) *Ledger {
	t.Helper()
	l := newLedger(t)
	if err := Open(l, "a1", "Ada", 10_000); err != nil {
		t.Fatalf("Open(a1) = %v, want nil", err)
	}
	if err := Open(l, "b2", "Bob", 0); err != nil {
		t.Fatalf("Open(b2) = %v, want nil", err)
	}
	return l
}

func wantBalance(t *testing.T, l *Ledger, id string, want int64) {
	t.Helper()
	got, err := Balance(l, id)
	if err != nil || got != want {
		t.Errorf("Balance(%q) = %d, %v, want %d, nil", id, got, err, want)
	}
}

func TestOpen(t *testing.T) {
	l := setup(t)
	a, ok := l.Accounts["a1"]
	if !ok || a == nil {
		t.Fatalf("Accounts[\"a1\"] missing after Open")
	}
	if *a != (Account{ID: "a1", Owner: "Ada", Balance: 10_000}) {
		t.Errorf("Accounts[\"a1\"] = %+v, want {a1 Ada 10000}", *a)
	}
	if want := []Tx{{"", "a1", 10_000}}; !slices.Equal(l.History, want) {
		t.Errorf("History = %v, want %v (no entry for a zero opening balance)", l.History, want)
	}

	tests := []struct {
		id      string
		initial int64
		want    error
	}{
		{"a1", 5, ErrAccountExists},
		{"c3", -1, ErrInvalidAmount},
		{"c3", math.MinInt64, ErrInvalidAmount},
	}
	for _, tt := range tests {
		if err := Open(l, tt.id, "X", tt.initial); !errors.Is(err, tt.want) {
			t.Errorf("Open(%q, %d) = %v, want %v", tt.id, tt.initial, err, tt.want)
		}
	}
	if err := Open(l, "", "Nobody", 0); err == nil {
		t.Errorf("Open(\"\") = nil, want error")
	}
	if len(l.Accounts) != 2 || len(l.History) != 1 {
		t.Errorf("failed Open calls changed the ledger: %d accounts, history %v", len(l.Accounts), l.History)
	}
	if l.Accounts["a1"].Owner != "Ada" || l.Accounts["a1"].Balance != 10_000 {
		t.Errorf("duplicate Open overwrote a1: %+v", *l.Accounts["a1"])
	}
}

func TestDepositWithdraw(t *testing.T) {
	l := setup(t)
	if err := Deposit(l, "b2", 1_999); err != nil {
		t.Fatalf("Deposit(b2, 1999) = %v", err)
	}
	if err := Withdraw(l, "a1", 2_500); err != nil {
		t.Fatalf("Withdraw(a1, 2500) = %v", err)
	}
	if err := Withdraw(l, "b2", 1_999); err != nil { // exactly the balance
		t.Fatalf("Withdraw(b2, 1999) = %v", err)
	}
	wantBalance(t, l, "a1", 7_500)
	wantBalance(t, l, "b2", 0)
	if l.Accounts["a1"].Balance != 7_500 {
		t.Errorf("Accounts[a1].Balance = %d, want 7500 (update the stored *Account)", l.Accounts["a1"].Balance)
	}
	want := []Tx{{"", "a1", 10_000}, {"", "b2", 1_999}, {"a1", "", 2_500}, {"b2", "", 1_999}}
	if !slices.Equal(l.History, want) {
		t.Fatalf("History = %v, want %v", l.History, want)
	}

	bad := []struct {
		name   string
		op     func(*Ledger, string, int64) error
		id     string
		amount int64
		want   error
	}{
		{"deposit unknown", Deposit, "zz", 100, ErrNoAccount},
		{"deposit unknown+invalid", Deposit, "zz", -5, ErrNoAccount},
		{"deposit zero", Deposit, "a1", 0, ErrInvalidAmount},
		{"deposit negative", Deposit, "a1", -100, ErrInvalidAmount},
		{"withdraw unknown", Withdraw, "", 100, ErrNoAccount},
		{"withdraw zero", Withdraw, "a1", 0, ErrInvalidAmount},
		{"withdraw negative", Withdraw, "a1", -1, ErrInvalidAmount},
		{"withdraw too much", Withdraw, "a1", 7_501, ErrInsufficient},
		{"withdraw from empty", Withdraw, "b2", 1, ErrInsufficient},
	}
	for _, tt := range bad {
		if err := tt.op(l, tt.id, tt.amount); !errors.Is(err, tt.want) {
			t.Errorf("%s: error = %v, want %v", tt.name, err, tt.want)
		}
	}
	wantBalance(t, l, "a1", 7_500)
	wantBalance(t, l, "b2", 0)
	if len(l.History) != len(want) {
		t.Errorf("failed operations added history: %v", l.History)
	}
}

func TestTransfer(t *testing.T) {
	l := setup(t)
	if err := Transfer(l, "a1", "b2", 2_550); err != nil {
		t.Fatalf("Transfer(a1, b2, 2550) = %v", err)
	}
	wantBalance(t, l, "a1", 7_450)
	wantBalance(t, l, "b2", 2_550)
	if err := Transfer(l, "b2", "a1", 2_550); err != nil { // whole balance
		t.Fatalf("Transfer(b2, a1, 2550) = %v", err)
	}
	wantBalance(t, l, "a1", 10_000)
	wantBalance(t, l, "b2", 0)

	bad := []struct {
		from, to string
		amount   int64
		want     error
	}{
		{"zz", "a1", 100, ErrNoAccount},
		{"a1", "zz", 100, ErrNoAccount},
		{"zz", "zz", -1, ErrNoAccount},
		{"a1", "a1", 100, ErrSameAccount},
		{"a1", "a1", -100, ErrSameAccount},
		{"a1", "b2", 0, ErrInvalidAmount},
		{"a1", "b2", -50, ErrInvalidAmount},
		{"b2", "a1", 1, ErrInsufficient},
		{"a1", "b2", 10_001, ErrInsufficient},
		{"a1", "b2", math.MaxInt64, ErrInsufficient},
	}
	for _, tt := range bad {
		if err := Transfer(l, tt.from, tt.to, tt.amount); !errors.Is(err, tt.want) {
			t.Errorf("Transfer(%q, %q, %d) = %v, want %v", tt.from, tt.to, tt.amount, err, tt.want)
		}
	}
	wantBalance(t, l, "a1", 10_000)
	wantBalance(t, l, "b2", 0)
	want := []Tx{{"", "a1", 10_000}, {"a1", "b2", 2_550}, {"b2", "a1", 2_550}}
	if !slices.Equal(l.History, want) {
		t.Errorf("History = %v, want %v", l.History, want)
	}
}

func TestBalanceUnknown(t *testing.T) {
	l := setup(t)
	if got, err := Balance(l, "nope"); !errors.Is(err, ErrNoAccount) || got != 0 {
		t.Errorf("Balance(nope) = %d, %v, want 0, ErrNoAccount", got, err)
	}
}

func TestStatement(t *testing.T) {
	l := setup(t)
	if err := Open(l, "c3", "Cy", 500); err != nil {
		t.Fatalf("Open(c3) = %v", err)
	}
	steps := []error{
		Transfer(l, "a1", "b2", 1_000),
		Deposit(l, "c3", 250),
		Transfer(l, "b2", "c3", 400),
		Withdraw(l, "a1", 100),
	}
	for i, err := range steps {
		if err != nil {
			t.Fatalf("step %d: %v", i, err)
		}
	}

	tests := []struct {
		id   string
		want []Tx
	}{
		{"a1", []Tx{{"", "a1", 10_000}, {"a1", "b2", 1_000}, {"a1", "", 100}}},
		{"b2", []Tx{{"a1", "b2", 1_000}, {"b2", "c3", 400}}},
		{"c3", []Tx{{"", "c3", 500}, {"", "c3", 250}, {"b2", "c3", 400}}},
	}
	for _, tt := range tests {
		got, err := Statement(l, tt.id)
		if err != nil || !slices.Equal(got, tt.want) {
			t.Errorf("Statement(%q) = %v, %v, want %v, nil", tt.id, got, err, tt.want)
		}
	}

	// The statement must be an independent copy.
	history := slices.Clone(l.History)
	got, _ := Statement(l, "a1")
	if len(got) == 0 {
		t.Fatalf("Statement(a1) is empty")
	}
	got[0].Amount = -1
	_ = append(got[:1], Tx{"x", "y", 1})
	if !slices.Equal(l.History, history) {
		t.Errorf("modifying a statement changed l.History: %v", l.History)
	}

	if got, err := Statement(l, "nope"); !errors.Is(err, ErrNoAccount) || got != nil {
		t.Errorf("Statement(nope) = %v, %v, want nil, ErrNoAccount", got, err)
	}
	fresh := setup(t)
	if got, err := Statement(fresh, "b2"); err != nil || len(got) != 0 {
		t.Errorf("Statement(new b2) = %v, %v, want empty, nil", got, err)
	}
}

func TestTotal(t *testing.T) {
	l := newLedger(t)
	if got := Total(l); got != 0 {
		t.Errorf("Total(empty) = %d, want 0", got)
	}
	l = setup(t)
	_ = Open(l, "c3", "Cy", 1_234)
	_ = Transfer(l, "a1", "c3", 5_000)
	_ = Transfer(l, "c3", "b2", 6_000)
	_ = Deposit(l, "b2", 66)
	_ = Withdraw(l, "a1", 300)
	if got, want := Total(l), int64(10_000+1_234+66-300); got != want {
		t.Errorf("Total = %d, want %d", got, want)
	}
}

func TestFormatCents(t *testing.T) {
	tests := []struct {
		c    int64
		want string
	}{
		{0, "0.00"},
		{5, "0.05"},
		{10, "0.10"},
		{99, "0.99"},
		{100, "1.00"},
		{1234, "12.34"},
		{7450, "74.50"},
		{1_000_000_01, "1000000.01"},
		{-5, "-0.05"},
		{-150, "-1.50"},
		{-100, "-1.00"},
		{math.MaxInt64, "92233720368547758.07"},
		{math.MinInt64, "-92233720368547758.08"},
	}
	for _, tt := range tests {
		if got := FormatCents(tt.c); got != tt.want {
			t.Errorf("FormatCents(%d) = %q, want %q", tt.c, got, tt.want)
		}
	}
}
