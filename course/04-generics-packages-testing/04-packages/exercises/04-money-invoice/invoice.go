// Package invoice builds invoices priced in money.Amount.
package invoice

import (
	"errors"

	"learngo/course/04-generics-packages-testing/04-packages/exercises/04-money-invoice/money"
)

// Errors returned by Invoice.Add.
var (
	ErrNoItem = errors.New("invoice: item name is empty")
	ErrQty    = errors.New("invoice: quantity must be at least 1")
)

// Line is one invoice line: Qty units of Item at Price each.
type Line struct {
	Item  string
	Price money.Amount
	Qty   int
}

// Invoice is a list of lines. The zero value is an empty invoice.
type Invoice struct {
	lines []Line
}

// Subtotal returns Price * Qty.
func (l Line) Subtotal() money.Amount {
	// TODO
	return money.Amount{}
}

// Add appends a line. It fails with ErrNoItem for an empty item name and
// with ErrQty when qty is below 1.
func (inv *Invoice) Add(item string, price money.Amount, qty int) error {
	// TODO
	return nil
}

// Lines returns a copy of the invoice lines.
func (inv *Invoice) Lines() []Line {
	// TODO
	return nil
}

// Total returns the sum of all line subtotals.
func (inv *Invoice) Total() money.Amount {
	// TODO
	return money.Amount{}
}

// String renders one "<qty> x <item> @ <price> = <subtotal>" row per line,
// followed by "Total: <total>".
func (inv *Invoice) String() string {
	// TODO
	return ""
}
