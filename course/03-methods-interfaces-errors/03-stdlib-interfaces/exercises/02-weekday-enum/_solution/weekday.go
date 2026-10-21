package days

import "fmt"

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

var names = [...]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

func (d Weekday) valid() bool { return d >= Sunday && d <= Saturday }

// String returns the English name of d, or "Weekday(N)" if d is invalid.
func (d Weekday) String() string {
	if !d.valid() {
		return fmt.Sprintf("Weekday(%d)", int(d)) // int(d): no String method, no recursion
	}
	return names[d]
}

// Short returns the three-letter abbreviation of d, or "???" if d is invalid.
func (d Weekday) Short() string {
	if !d.valid() {
		return "???"
	}
	return names[d][:3]
}

// Next returns the day after d (Saturday wraps to Sunday). An invalid day is
// returned unchanged.
func (d Weekday) Next() Weekday {
	if !d.valid() {
		return d
	}
	return (d + 1) % 7
}
