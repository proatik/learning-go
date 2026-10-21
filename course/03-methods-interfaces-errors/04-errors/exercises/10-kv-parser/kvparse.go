package kvparse

import "errors"

// Problems a ParseError can report.
var (
	ErrMissingEquals = errors.New("missing '='")
	ErrEmptyKey      = errors.New("empty key")
	ErrInvalidKey    = errors.New("invalid key character")
	ErrDuplicateKey  = errors.New("duplicate key")
)

// ParseError reports where in the input parsing failed.
type ParseError struct {
	Line int   // 1-based
	Col  int   // 1-based, counted in runes
	Err  error // one of the sentinels above
}

// Error returns "line <Line>, col <Col>: <Err>".
func (e *ParseError) Error() string {
	// TODO
	return ""
}

// Unwrap returns the underlying sentinel.
func (e *ParseError) Unwrap() error {
	// TODO
	return nil
}

// Parse parses lines of "key=value;key2=value2" into a map.
func Parse(input string) (map[string]string, error) {
	// TODO
	return nil, nil
}
