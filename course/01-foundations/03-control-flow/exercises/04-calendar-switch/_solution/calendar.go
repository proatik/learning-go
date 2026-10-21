package calendar

// DayType classifies a three-letter day abbreviation such as "Mon" as
// "weekday" or "weekend". Unknown values return "invalid".
func DayType(day string) string {
	switch day {
	case "Mon", "Tue", "Wed", "Thu", "Fri":
		return "weekday"
	case "Sat", "Sun":
		return "weekend"
	default:
		return "invalid"
	}
}

// DaysInMonth returns the number of days in month (1-12) of year.
// It returns 0 for an invalid month.
func DaysInMonth(month, year int) int {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
			return 29
		}
		return 28
	default:
		return 0
	}
}
