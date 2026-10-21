package clock

// AddMinutes adds delta minutes to the time h:m and returns the result as
// "HH:MM" on a 24-hour clock, wrapping around midnight in both directions.
func AddMinutes(h, m, delta int) string {
	// TODO
	return ""
}

// ParseClock parses a strict "HH:MM" string (00:00 to 23:59).
func ParseClock(s string) (h, m int, err error) {
	// TODO
	return 0, 0, nil
}
