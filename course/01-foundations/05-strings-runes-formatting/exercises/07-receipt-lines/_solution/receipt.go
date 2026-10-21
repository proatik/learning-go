package receipt

import (
	"fmt"
	"strings"
)

// Column layout (37 runes): item 12 | space | qty 3 | " x " | unit 7 | " = " | total 8

// Header returns the column titles of the receipt.
func Header() string {
	return fmt.Sprintf("%-12s %3s   %7s   %8s", "Item", "Qty", "Price", "Total")
}

// Separator returns a line of 37 dashes.
func Separator() string {
	return strings.Repeat("-", 37)
}

// Line returns one receipt line: item (cut to 12 runes), quantity,
// unit price and line total.
func Line(item string, qty int, unit float64) string {
	// %-12.12s: left-aligned, padded to 12 runes, truncated to 12 runes.
	return fmt.Sprintf("%-12.12s %3d x %7.2f = %8.2f", item, qty, unit, float64(qty)*unit)
}

// Total returns the final TOTAL line.
func Total(total float64) string {
	return fmt.Sprintf("%-29s%8.2f", "TOTAL", total)
}
