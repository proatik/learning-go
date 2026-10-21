package receipt

import (
	"fmt"
	"io"
)

// Item is one line of a receipt.
type Item struct {
	Name  string
	Qty   int
	Cents int // price of one unit, in cents
}

func dollars(cents int) string {
	return fmt.Sprintf("$%d.%02d", cents/100, cents%100)
}

// WriteReceipt writes one line per item followed by a TOTAL line to w.
// It returns the first write error, or nil.
func WriteReceipt(w io.Writer, items []Item) error {
	if len(items) == 0 {
		if _, err := fmt.Fprintln(w, "(no items)"); err != nil {
			return err
		}
	}
	total := 0
	for _, it := range items {
		line := it.Qty * it.Cents
		total += line
		if _, err := fmt.Fprintf(w, "%s x%d = %s\n", it.Name, it.Qty, dollars(line)); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w, "TOTAL = %s\n", dollars(total))
	return err
}
