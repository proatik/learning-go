package weekday

// Weekday is a day of the week, starting with Sunday = 0.
type Weekday int

// The days of the week. Lines without an expression repeat "Weekday = iota".
const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

// daysPerWeek is an untyped constant, so it works with Weekday arithmetic.
const daysPerWeek = 7

// IsWeekend reports whether d is Saturday or Sunday.
func IsWeekend(d Weekday) bool {
	return d == Saturday || d == Sunday
}

// Next returns the day after d. The day after Saturday is Sunday.
func Next(d Weekday) Weekday {
	return (d + 1) % daysPerWeek
}

// Prev returns the day before d. The day before Sunday is Saturday.
func Prev(d Weekday) Weekday {
	// (d - 1) % 7 would give -1 for Sunday, because % keeps the sign of the
	// left operand. Adding 6 (≡ -1 mod 7) keeps everything non-negative.
	return (d + daysPerWeek - 1) % daysPerWeek
}
