package clock

import "fmt"

const (
	secondsPerMinute = 60
	secondsPerHour   = 60 * secondsPerMinute
	secondsPerDay    = 24 * secondsPerHour
)

// Split breaks a number of seconds into hours, minutes and seconds.
// For negative totals every part is negative (or zero).
func Split(total int) (h, m, s int) {
	h = total / secondsPerHour
	m = total % secondsPerHour / secondsPerMinute
	s = total % secondsPerMinute
	return h, m, s
}

// Clock formats a duration in seconds as "HH:MM:SS", with a leading "-" for
// negative durations. Hours may use more than two digits.
func Clock(total int) string {
	sign := ""
	if total < 0 {
		sign = "-"
		total = -total
	}
	h, m, s := Split(total)
	return fmt.Sprintf("%s%02d:%02d:%02d", sign, h, m, s)
}

// TimeOfDay formats total seconds after midnight as a wall-clock time
// "HH:MM:SS", wrapping around so the result is always within one day.
func TimeOfDay(total int) string {
	// % keeps the sign of total, so the first remainder is in
	// (-secondsPerDay, secondsPerDay). Adding a day and taking % again maps
	// it into [0, secondsPerDay).
	t := (total%secondsPerDay + secondsPerDay) % secondsPerDay
	h, m, s := Split(t)
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
