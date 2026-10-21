package kvparse

import (
	"errors"
	"fmt"
	"strings"
)

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
	return fmt.Sprintf("line %d, col %d: %v", e.Line, e.Col, e.Err)
}

// Unwrap returns the underlying sentinel.
func (e *ParseError) Unwrap() error {
	return e.Err
}

// Parse parses lines of "key=value;key2=value2" into a map.
func Parse(input string) (map[string]string, error) {
	out := make(map[string]string)
	for i, line := range strings.Split(input, "\n") {
		r := []rune(line) // index i in r is column i+1
		for start := 0; start <= len(r); {
			end := start
			for end < len(r) && r[end] != ';' {
				end++
			}
			if err := parsePair(r, start, end, i+1, out); err != nil {
				return nil, err
			}
			start = end + 1
		}
	}
	return out, nil
}

// parsePair parses the segment r[start:end] of line number line into out.
func parsePair(r []rune, start, end, line int, out map[string]string) error {
	first := start
	for first < end && isSpace(r[first]) {
		first++
	}
	if first == end {
		return nil // blank segment
	}
	eq := -1
	for i := first; i < end; i++ {
		if r[i] == '=' {
			eq = i
			break
		}
	}
	if eq < 0 {
		return &ParseError{Line: line, Col: first + 1, Err: ErrMissingEquals}
	}
	keyEnd := eq
	for keyEnd > first && isSpace(r[keyEnd-1]) {
		keyEnd--
	}
	if keyEnd == first {
		return &ParseError{Line: line, Col: eq + 1, Err: ErrEmptyKey}
	}
	for i := first; i < keyEnd; i++ {
		if !isKeyRune(r[i]) {
			return &ParseError{Line: line, Col: i + 1, Err: ErrInvalidKey}
		}
	}
	key := string(r[first:keyEnd])
	if _, dup := out[key]; dup {
		return &ParseError{Line: line, Col: first + 1, Err: ErrDuplicateKey}
	}
	out[key] = strings.Trim(string(r[eq+1:end]), " \t")
	return nil
}

func isSpace(c rune) bool { return c == ' ' || c == '\t' }

func isKeyRune(c rune) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' ||
		c == '_' || c == '.' || c == '-'
}
