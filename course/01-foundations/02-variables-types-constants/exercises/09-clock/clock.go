package clock

// Split breaks a number of seconds into hours, minutes and seconds.
// For negative totals every part is negative (or zero).
func Split(total int) (h, m, s int) {
	return 0, 0, 0 // TODO
}

// Clock formats a duration in seconds as "HH:MM:SS", with a leading "-" for
// negative durations. Hours may use more than two digits.
func Clock(total int) string {
	return "" // TODO
}

// TimeOfDay formats total seconds after midnight as a wall-clock time
// "HH:MM:SS", wrapping around so the result is always within one day.
func TimeOfDay(total int) string {
	return "" // TODO
}
