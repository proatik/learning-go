package duration

import (
	"strconv"
	"strings"
)

// Unit sizes in seconds.
const (
	Minute = 60
	Hour   = 60 * Minute
	Day    = 24 * Hour
)

// HumanDuration formats seconds as "1d 2h 3m 4s", omitting zero parts.
// Zero is "0s"; negative durations get a leading "-".
func HumanDuration(seconds int) string {
	if seconds == 0 {
		return "0s"
	}

	var sb strings.Builder
	if seconds < 0 {
		sb.WriteByte('-')
		seconds = -seconds
	}

	// add appends "<n><unit>" when n is non-zero, preceded by a space
	// if an earlier part was already written. The closure captures
	// sb and wrote.
	wrote := false
	add := func(n int, unit string) {
		if n == 0 {
			return
		}
		if wrote {
			sb.WriteByte(' ')
		}
		wrote = true
		sb.WriteString(strconv.Itoa(n))
		sb.WriteString(unit)
	}

	add(seconds/Day, "d")
	add(seconds%Day/Hour, "h")
	add(seconds%Hour/Minute, "m")
	add(seconds%Minute, "s")
	return sb.String()
}
