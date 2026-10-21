package convert

// Double parses s as a base-10 integer and returns twice its value.
// If s is not a valid integer it returns 0 and the parse error.
func Double(s string) (int, error) {
	return 0, nil // TODO: strconv.Atoi
}

// AddStrings parses a and b as float64 numbers and returns their sum.
// If either string is invalid it returns 0 and the parse error.
func AddStrings(a, b string) (float64, error) {
	return 0, nil // TODO: strconv.ParseFloat
}

// Label returns "item-" followed by the decimal digits of n.
func Label(n int) string {
	return "" // TODO: strconv.Itoa
}

// Binary returns n written in base 2.
func Binary(n int) string {
	return "" // TODO: strconv.FormatInt
}
