package must

import (
	"errors"
	"fmt"
	"strconv"
)

// ErrPortRange is wrapped by MustPort's panic value for ports outside 1..65535.
var ErrPortRange = errors.New("port out of range")

// Must returns n, or panics with err if err is not nil.
func Must(n int, err error) int {
	if err != nil {
		panic(err)
	}
	return n
}

// MustAtoi is strconv.Atoi that panics with an error wrapping the strconv
// error, e.g. `MustAtoi("12a"): strconv.Atoi: parsing "12a": invalid syntax`.
func MustAtoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Errorf("MustAtoi(%q): %w", s, err))
	}
	return n
}

// MustPort parses s with MustAtoi and panics with an error wrapping
// ErrPortRange unless the result is in 1..65535.
func MustPort(s string) int {
	n := MustAtoi(s)
	if n < 1 || n > 65535 {
		panic(fmt.Errorf("MustPort(%q): %w", s, ErrPortRange))
	}
	return n
}
