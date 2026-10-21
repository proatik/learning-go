package weekday

// Weekday is a day of the week, starting with Sunday = 0.
type Weekday int

// TODO: use iota so the days are numbered 0 (Sunday) to 6 (Saturday).
const (
	Sunday    Weekday = 0
	Monday    Weekday = 0
	Tuesday   Weekday = 0
	Wednesday Weekday = 0
	Thursday  Weekday = 0
	Friday    Weekday = 0
	Saturday  Weekday = 0
)

// IsWeekend reports whether d is Saturday or Sunday.
func IsWeekend(d Weekday) bool {
	return false // TODO
}

// Next returns the day after d. The day after Saturday is Sunday.
func Next(d Weekday) Weekday {
	return 0 // TODO
}

// Prev returns the day before d. The day before Sunday is Saturday.
func Prev(d Weekday) Weekday {
	return 0 // TODO
}
