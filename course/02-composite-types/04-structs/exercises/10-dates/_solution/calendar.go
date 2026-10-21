package calendar

import (
	"cmp"
	"fmt"
)

// Date is a day in the Gregorian calendar.
type Date struct {
	Year, Month, Day int // Month 1–12, Day 1–31
}

// IsLeap reports whether year is a leap year.
func IsLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// DaysInMonth returns the number of days in month of year,
// or 0 if month is not in 1–12.
func DaysInMonth(year, month int) int {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		if IsLeap(year) {
			return 29
		}
		return 28
	}
	return 0
}

// Valid reports whether d is a real date with a year in 1–9999.
func Valid(d Date) bool {
	return d.Year >= 1 && d.Year <= 9999 &&
		d.Day >= 1 && d.Day <= DaysInMonth(d.Year, d.Month) // 0 for a bad month
}

// Compare returns -1 if a is before b, +1 if a is after b, and 0 if equal.
func Compare(a, b Date) int {
	return cmp.Or(
		cmp.Compare(a.Year, b.Year),
		cmp.Compare(a.Month, b.Month),
		cmp.Compare(a.Day, b.Day),
	)
}

// IsBefore reports whether a is strictly before b.
func IsBefore(a, b Date) bool {
	return Compare(a, b) < 0
}

const daysPer400Years = 146097 // 400*365 + 97 leap days

// dayNumber returns the number of days from 0001-01-01 to d (0 for 0001-01-01).
func dayNumber(d Date) int {
	y := d.Year - 1 // whole years before d.Year
	n := y*365 + y/4 - y/100 + y/400
	for m := 1; m < d.Month; m++ {
		n += DaysInMonth(d.Year, m)
	}
	return n + d.Day - 1
}

// fromDayNumber is the inverse of dayNumber.
func fromDayNumber(n int) Date {
	d := Date{Year: 1, Month: 1, Day: 1}
	d.Year += 400 * (n / daysPer400Years) // skip whole 400-year cycles
	n %= daysPer400Years
	for {
		yearLen := 365
		if IsLeap(d.Year) {
			yearLen = 366
		}
		if n < yearLen {
			break
		}
		n -= yearLen
		d.Year++
	}
	for n >= DaysInMonth(d.Year, d.Month) {
		n -= DaysInMonth(d.Year, d.Month)
		d.Month++
	}
	d.Day += n
	return d
}

// AddDays returns the date n days after d (before d if n is negative).
func AddDays(d Date, n int) Date {
	return fromDayNumber(dayNumber(d) + n)
}

// DaysBetween returns the number of days from a to b (negative if b is before a).
func DaysBetween(a, b Date) int {
	return dayNumber(b) - dayNumber(a)
}

// Format returns d as "YYYY-MM-DD".
func Format(d Date) string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}
