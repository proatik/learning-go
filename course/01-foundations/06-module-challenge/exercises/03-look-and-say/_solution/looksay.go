package looksay

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Next returns the look-and-say description of the digit string s.
func Next(s string) (string, error) {
	if s == "" {
		return "", errors.New("looksay: empty input")
	}

	var sb strings.Builder
	for i := 0; i < len(s); {
		c := s[i]
		if c < '0' || c > '9' {
			return "", fmt.Errorf("looksay: invalid byte %q at index %d", c, i)
		}
		// Find the end of the run of c.
		j := i + 1
		for j < len(s) && s[j] == c {
			j++
		}
		sb.WriteString(strconv.Itoa(j - i))
		sb.WriteByte(c)
		i = j
	}
	return sb.String(), nil
}

// Term returns the n-th term (1-based) of the look-and-say sequence that
// starts with "1".
func Term(n int) (string, error) {
	if n < 1 {
		return "", fmt.Errorf("looksay: term %d: n must be at least 1", n)
	}
	s := "1"
	for range n - 1 {
		// Next cannot fail here: every term consists only of digits.
		s, _ = Next(s)
	}
	return s, nil
}
