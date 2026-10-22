package invoice_test

import (
	"errors"
	"testing"

	invoice "learngo/course/04-generics-packages-testing/04-packages/exercises/04-money-invoice"
	"learngo/course/04-generics-packages-testing/04-packages/exercises/04-money-invoice/money"
)

func amt(t *testing.T, units, cents int64) money.Amount {
	t.Helper()
	a, err := money.New(units, cents)
	if err != nil {
		t.Fatalf("money.New(%d, %d) error = %v, want <nil>", units, cents, err)
	}
	return a
}

func TestSubtotal(t *testing.T) {
	tests := []struct {
		line invoice.Line
		want string
	}{
		{invoice.Line{Item: "Pen", Price: amt(t, 1, 50), Qty: 2}, "3.00"},
		{invoice.Line{Item: "Pad", Price: amt(t, 0, 99), Qty: 1}, "0.99"},
		{invoice.Line{Item: "Box", Price: amt(t, 12, 5), Qty: 10}, "120.50"},
		{invoice.Line{}, "0.00"},
	}
	for _, tt := range tests {
		if got := tt.line.Subtotal().String(); got != tt.want {
			t.Errorf("%+v.Subtotal() = %q, want %q", tt.line, got, tt.want)
		}
	}
}

func TestEmptyInvoice(t *testing.T) {
	var inv invoice.Invoice
	if got := inv.Total().String(); got != "0.00" {
		t.Errorf("empty Total() = %q, want %q", got, "0.00")
	}
	if got := inv.String(); got != "Total: 0.00" {
		t.Errorf("empty String() = %q, want %q", got, "Total: 0.00")
	}
	if got := inv.Lines(); len(got) != 0 {
		t.Errorf("empty Lines() = %v, want []", got)
	}
}

func TestAddErrors(t *testing.T) {
	var inv invoice.Invoice
	tests := []struct {
		item string
		qty  int
		want error
	}{
		{"", 1, invoice.ErrNoItem},
		{"", 0, invoice.ErrNoItem},
		{"Pen", 0, invoice.ErrQty},
		{"Pen", -3, invoice.ErrQty},
	}
	for _, tt := range tests {
		if err := inv.Add(tt.item, amt(t, 1, 0), tt.qty); !errors.Is(err, tt.want) {
			t.Errorf("Add(%q, 1.00, %d) error = %v, want %v", tt.item, tt.qty, err, tt.want)
		}
	}
	if got := len(inv.Lines()); got != 0 {
		t.Errorf("len(Lines()) after failed Adds = %d, want 0", got)
	}
}

func TestInvoice(t *testing.T) {
	var inv invoice.Invoice
	if err := inv.Add("Pen", amt(t, 1, 50), 2); err != nil {
		t.Fatalf("Add(%q, 1.50, 2) error = %v, want <nil>", "Pen", err)
	}
	if err := inv.Add("Pad", amt(t, 0, 99), 1); err != nil {
		t.Fatalf("Add(%q, 0.99, 1) error = %v, want <nil>", "Pad", err)
	}
	if got := inv.Total().String(); got != "3.99" {
		t.Errorf("Total() = %q, want %q", got, "3.99")
	}
	want := "2 x Pen @ 1.50 = 3.00\n1 x Pad @ 0.99 = 0.99\nTotal: 3.99"
	if got := inv.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}

	lines := inv.Lines()
	if len(lines) != 2 {
		t.Fatalf("len(Lines()) = %d, want 2", len(lines))
	}
	if lines[0].Item != "Pen" || lines[1].Item != "Pad" {
		t.Errorf("Lines() items = %q, %q, want %q, %q", lines[0].Item, lines[1].Item, "Pen", "Pad")
	}
	lines[0].Qty = 1000 // must not change the invoice
	if got := inv.Total().String(); got != "3.99" {
		t.Errorf("Total() after editing Lines() copy = %q, want %q", got, "3.99")
	}
}
