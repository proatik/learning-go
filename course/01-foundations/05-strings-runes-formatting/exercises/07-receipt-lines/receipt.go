package receipt

// Header returns the column titles of the receipt.
func Header() string {
	// TODO
	return ""
}

// Separator returns a line of 37 dashes.
func Separator() string {
	// TODO
	return ""
}

// Line returns one receipt line: item (cut to 12 runes), quantity,
// unit price and line total.
func Line(item string, qty int, unit float64) string {
	// TODO
	return ""
}

// Total returns the final TOTAL line.
func Total(total float64) string {
	// TODO
	return ""
}
