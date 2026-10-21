package days

// Weekday is a day of the week, starting with Sunday = 0.
type Weekday int

// The days of the week.
const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

// String returns the English name of d, or "Weekday(N)" if d is invalid.
func (d Weekday) String() string {
	// TODO
	return ""
}

// Short returns the three-letter abbreviation of d, or "???" if d is invalid.
func (d Weekday) Short() string {
	// TODO
	return ""
}

// Next returns the day after d (Saturday wraps to Sunday). An invalid day is
// returned unchanged.
func (d Weekday) Next() Weekday {
	// TODO
	return -1
}
