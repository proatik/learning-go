package ages

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Sentinel errors returned (possibly wrapped) by ParseAge.
var (
	ErrEmpty      = errors.New("empty age")
	ErrNotNumber  = errors.New("not a number")
	ErrOutOfRange = errors.New("out of range")
)

// ParseAge parses s (surrounding spaces allowed) as an age from 0 to 150.
func ParseAge(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, ErrEmpty
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("parse age %q: %w", s, ErrNotNumber)
	}
	if n < 0 || n > 150 {
		return 0, fmt.Errorf("age %d: %w", n, ErrOutOfRange)
	}
	return n, nil
}
