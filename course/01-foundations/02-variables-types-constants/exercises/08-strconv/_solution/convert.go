package convert

import "strconv"

// Double parses s as a base-10 integer and returns twice its value.
// If s is not a valid integer it returns 0 and the parse error.
func Double(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return n * 2, nil
}

// AddStrings parses a and b as float64 numbers and returns their sum.
// If either string is invalid it returns 0 and the parse error.
func AddStrings(a, b string) (float64, error) {
	x, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return 0, err
	}
	// err already exists; := is legal because y is new (err is reused).
	y, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return 0, err
	}
	return x + y, nil
}

// Label returns "item-" followed by the decimal digits of n.
func Label(n int) string {
	return "item-" + strconv.Itoa(n)
}

// Binary returns n written in base 2.
func Binary(n int) string {
	// FormatInt takes an int64, so convert explicitly.
	return strconv.FormatInt(int64(n), 2)
}
