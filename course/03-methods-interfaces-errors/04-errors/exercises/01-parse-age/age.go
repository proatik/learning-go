package ages

import "errors"

// Sentinel errors returned (possibly wrapped) by ParseAge.
var (
	ErrEmpty      = errors.New("empty age")
	ErrNotNumber  = errors.New("not a number")
	ErrOutOfRange = errors.New("out of range")
)

// ParseAge parses s (surrounding spaces allowed) as an age from 0 to 150.
func ParseAge(s string) (int, error) {
	// TODO
	return 0, nil
}
