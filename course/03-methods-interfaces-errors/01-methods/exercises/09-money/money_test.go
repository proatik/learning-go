package money

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		dollars, cents int64
		want           Money
	}{
		{12, 34, 1234},
		{0, 0, 0},
		{0, 7, 7},
		{5, 0, 500},
		{-1, -50, -150},
	}
	for _, tt := range tests {
		if got := New(tt.dollars, tt.cents); got != tt.want {
			t.Errorf("New(%d, %d) = %d, want %d", tt.dollars, tt.cents, int64(got), int64(tt.want))
		}
	}
}

func TestArithmetic(t *testing.T) {
	tests := []struct {
		desc string
		got  Money
		want Money
	}{
		{"Money(1234).Add(66)", Money(1234).Add(66), 1300},
		{"Money(0).Add(0)", Money(0).Add(0), 0},
		{"Money(500).Sub(750)", Money(500).Sub(750), -250},
		{"Money(-5).Sub(-5)", Money(-5).Sub(-5), 0},
		{"Money(125).Times(3)", Money(125).Times(3), 375},
		{"Money(125).Times(0)", Money(125).Times(0), 0},
		{"Money(125).Times(-2)", Money(125).Times(-2), -250},
		{"Money(-50).Abs()", Money(-50).Abs(), 50},
		{"Money(50).Abs()", Money(50).Abs(), 50},
		{"Money(0).Abs()", Money(0).Abs(), 0},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %d, want %d", tt.desc, int64(tt.got), int64(tt.want))
		}
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		m    Money
		want string
	}{
		{0, "$0.00"},
		{7, "$0.07"},
		{99, "$0.99"},
		{100, "$1.00"},
		{1234, "$12.34"},
		{99999, "$999.99"},
		{100000, "$1,000.00"},
		{123456, "$1,234.56"},
		{12345678, "$123,456.78"},
		{123456789, "$1,234,567.89"},
		{100000000000, "$1,000,000,000.00"},
		{-1, "-$0.01"},
		{-50, "-$0.50"},
		{-100000, "-$1,000.00"},
		{-123456789, "-$1,234,567.89"},
	}
	for _, tt := range tests {
		if got := tt.m.String(); got != tt.want {
			t.Errorf("Money(%d).String() = %q, want %q", int64(tt.m), got, tt.want)
		}
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		s    string
		want Money
		ok   bool
	}{
		{"$12.34", 1234, true},
		{"$0.07", 7, true},
		{"$0.00", 0, true},
		{"$5", 500, true},
		{"$1,234.56", 123456, true},
		{"$1234.56", 123456, true},
		{"$1,234,567.89", 123456789, true},
		{"$999,999", 99999900, true},
		{"-$0.50", -50, true},
		{"-$1,000", -100000, true},
		{"-$0", 0, true},
		{"$007.10", 710, true},
		{"", 0, false},
		{"$", 0, false},
		{"-", 0, false},
		{"12.34", 0, false},
		{"$1.2", 0, false},
		{"$1.234", 0, false},
		{"$1.", 0, false},
		{"$.50", 0, false},
		{"$1,23.00", 0, false},
		{"$12,34", 0, false},
		{"$1,,000", 0, false},
		{",$100", 0, false},
		{"$,100", 0, false},
		{"$1000,000", 0, false},
		{"$100,", 0, false},
		{"$-5", 0, false},
		{"--$1", 0, false},
		{"+$1", 0, false},
		{"$ 5", 0, false},
		{"$5 ", 0, false},
		{"$1.2a", 0, false},
		{"$١٢", 0, false},
		{"$1.00.00", 0, false},
	}
	for _, tt := range tests {
		got, ok := Parse(tt.s)
		if got != tt.want || ok != tt.ok {
			t.Errorf("Parse(%q) = (%d, %v), want (%d, %v)", tt.s, int64(got), ok, int64(tt.want), tt.ok)
		}
	}
}

func TestRoundTrip(t *testing.T) {
	for _, m := range []Money{0, 1, -1, 99, 100, 123456, -98765432, 100000000000} {
		s := m.String()
		got, ok := Parse(s)
		if !ok || got != m {
			t.Errorf("Parse(%q) = (%d, %v), want (%d, true)", s, int64(got), ok, int64(m))
		}
	}
}

func ExampleMoney_String() {
	price := New(12, 34)
	total := price.Times(3).Sub(New(0, 2))
	fmt.Println(price)
	fmt.Println(total)
	fmt.Printf("refund: %v\n", Money(-150000))
	// Output:
	// $12.34
	// $37.00
	// refund: -$1,500.00
}
