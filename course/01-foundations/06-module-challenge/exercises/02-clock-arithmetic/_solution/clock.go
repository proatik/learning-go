package clock

import (
	"errors"
	"fmt"
)

const minutesPerDay = 24 * 60

// AddMinutes adds delta minutes to the time h:m and returns the result as
// "HH:MM" on a 24-hour clock, wrapping around midnight in both directions.
func AddMinutes(h, m, delta int) string {
	total := h*60 + m + delta
	// In Go, % keeps the sign of the dividend, so shift negatives back
	// into [0, minutesPerDay).
	total = (total%minutesPerDay + minutesPerDay) % minutesPerDay
	return fmt.Sprintf("%02d:%02d", total/60, total%60)
}

// ParseClock parses a strict "HH:MM" string (00:00 to 23:59).
func ParseClock(s string) (h, m int, err error) {
	if len(s) != 5 || s[2] != ':' {
		return 0, 0, fmt.Errorf("parse clock %q: want HH:MM", s)
	}
	h, err = twoDigits(s[0:2])
	if err != nil {
		return 0, 0, fmt.Errorf("parse clock %q: %v", s, err)
	}
	m, err = twoDigits(s[3:5])
	if err != nil {
		return 0, 0, fmt.Errorf("parse clock %q: %v", s, err)
	}
	if h > 23 || m > 59 {
		return 0, 0, fmt.Errorf("parse clock %q: out of range", s)
	}
	return h, m, nil
}

// twoDigits converts a two-byte string of ASCII digits to a number.
func twoDigits(s string) (int, error) {
	n := 0
	for i := range len(s) {
		c := s[i]
		if c < '0' || c > '9' {
			return 0, errors.New("not a digit")
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}
