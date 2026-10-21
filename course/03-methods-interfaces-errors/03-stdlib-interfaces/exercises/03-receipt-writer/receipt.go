package receipt

import "io"

// Item is one line of a receipt.
type Item struct {
	Name  string
	Qty   int
	Cents int // price of one unit, in cents
}

// WriteReceipt writes one line per item followed by a TOTAL line to w.
// It returns the first write error, or nil.
func WriteReceipt(w io.Writer, items []Item) error {
	// TODO
	return nil
}
