package leap

// IsLeap reports whether year is a leap year in the Gregorian calendar.
func IsLeap(year int) bool {
	// The most specific rule must be tested first.
	if year%400 == 0 {
		return true
	} else if year%100 == 0 {
		return false
	} else if year%4 == 0 {
		return true
	}
	return false
}

// An equivalent single expression:
//
//	return year%4 == 0 && (year%100 != 0 || year%400 == 0)

// DaysInYear returns the number of days in year: 366 for leap years, otherwise 365.
func DaysInYear(year int) int {
	if IsLeap(year) {
		return 366
	}
	return 365
}
