package must

import "errors"

// ErrPortRange is wrapped by MustPort's panic value for ports outside 1..65535.
var ErrPortRange = errors.New("port out of range")

// Must returns n, or panics with err if err is not nil.
func Must(n int, err error) int {
	// TODO
	return 0
}

// MustAtoi is strconv.Atoi that panics with an error wrapping the strconv
// error, e.g. `MustAtoi("12a"): strconv.Atoi: parsing "12a": invalid syntax`.
func MustAtoi(s string) int {
	// TODO
	return 0
}

// MustPort parses s with MustAtoi and panics with an error wrapping
// ErrPortRange unless the result is in 1..65535.
func MustPort(s string) int {
	// TODO
	return 0
}
