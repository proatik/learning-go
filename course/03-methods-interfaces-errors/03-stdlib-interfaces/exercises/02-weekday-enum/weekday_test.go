package days

import (
	"fmt"
	"testing"
)

var _ fmt.Stringer = Monday

func TestString(t *testing.T) {
	tests := []struct {
		d    Weekday
		want string
	}{
		{Sunday, "Sunday"},
		{Monday, "Monday"},
		{Tuesday, "Tuesday"},
		{Wednesday, "Wednesday"},
		{Thursday, "Thursday"},
		{Friday, "Friday"},
		{Saturday, "Saturday"},
		{Weekday(7), "Weekday(7)"},
		{Weekday(9), "Weekday(9)"},
		{Weekday(-1), "Weekday(-1)"},
		{Weekday(1000), "Weekday(1000)"},
	}
	for _, tt := range tests {
		if got := fmt.Sprint(tt.d); got != tt.want {
			t.Errorf("fmt.Sprint(Weekday(%d)) = %q, want %q", int(tt.d), got, tt.want)
		}
	}
}

func TestStringInFormats(t *testing.T) {
	tests := []struct {
		expr, got, want string
	}{
		{`fmt.Sprintf("%v and %s", Saturday, Sunday)`, fmt.Sprintf("%v and %s", Saturday, Sunday), "Saturday and Sunday"},
		{`fmt.Sprint([]Weekday{Monday, Friday})`, fmt.Sprint([]Weekday{Monday, Friday}), "[Monday Friday]"},
		{`fmt.Sprintf("%d", Friday)`, fmt.Sprintf("%d", Friday), "5"},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %q, want %q", tt.expr, tt.got, tt.want)
		}
	}
}

func TestShort(t *testing.T) {
	tests := []struct {
		d    Weekday
		want string
	}{
		{Sunday, "Sun"},
		{Monday, "Mon"},
		{Tuesday, "Tue"},
		{Wednesday, "Wed"},
		{Thursday, "Thu"},
		{Friday, "Fri"},
		{Saturday, "Sat"},
		{Weekday(7), "???"},
		{Weekday(-3), "???"},
	}
	for _, tt := range tests {
		if got := tt.d.Short(); got != tt.want {
			t.Errorf("Weekday(%d).Short() = %q, want %q", int(tt.d), got, tt.want)
		}
	}
}

func TestNext(t *testing.T) {
	tests := []struct {
		d, want Weekday
	}{
		{Sunday, Monday},
		{Monday, Tuesday},
		{Friday, Saturday},
		{Saturday, Sunday},
		{Weekday(9), Weekday(9)},
		{Weekday(-1), Weekday(-1)},
	}
	for _, tt := range tests {
		if got := tt.d.Next(); got != tt.want {
			t.Errorf("Weekday(%d).Next() = %d, want %d", int(tt.d), int(got), int(tt.want))
		}
	}
}
