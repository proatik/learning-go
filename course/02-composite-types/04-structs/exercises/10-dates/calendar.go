package calendar

// Date is a day in the Gregorian calendar.
type Date struct {
	Year, Month, Day int // Month 1–12, Day 1–31
}

// IsLeap reports whether year is a leap year.
func IsLeap(year int) bool {
	// TODO
	return false
}

// DaysInMonth returns the number of days in month of year,
// or 0 if month is not in 1–12.
func DaysInMonth(year, month int) int {
	// TODO
	return 0
}

// Valid reports whether d is a real date with a year in 1–9999.
func Valid(d Date) bool {
	// TODO
	return false
}

// Compare returns -1 if a is before b, +1 if a is after b, and 0 if equal.
func Compare(a, b Date) int {
	// TODO
	return 0
}

// IsBefore reports whether a is strictly before b.
func IsBefore(a, b Date) bool {
	// TODO
	return false
}

// AddDays returns the date n days after d (before d if n is negative).
func AddDays(d Date, n int) Date {
	// TODO
	return Date{}
}

// DaysBetween returns the number of days from a to b (negative if b is before a).
func DaysBetween(a, b Date) int {
	// TODO
	return 0
}

// Format returns d as "YYYY-MM-DD".
func Format(d Date) string {
	// TODO
	return ""
}
