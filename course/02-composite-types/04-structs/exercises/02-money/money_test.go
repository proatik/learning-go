package money

import (
	"slices"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		dollars, cents int64
		want           Money
	}{
		{12, 5, Money{Cents: 1205}},
		{0, 0, Money{}},
		{0, 99, Money{Cents: 99}},
		{3, 0, Money{Cents: 300}},
		{-1, -50, Money{Cents: -150}},
		{1, 150, Money{Cents: 250}},
	}
	for _, tt := range tests {
		if got := New(tt.dollars, tt.cents); got != tt.want {
			t.Errorf("New(%d, %d) = %+v, want %+v", tt.dollars, tt.cents, got, tt.want)
		}
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		a, b, want Money
	}{
		{Money{150}, Money{275}, Money{425}},
		{Money{}, Money{}, Money{}},
		{Money{100}, Money{-250}, Money{-150}},
		{Money{-5}, Money{5}, Money{}},
	}
	for _, tt := range tests {
		if got := Add(tt.a, tt.b); got != tt.want {
			t.Errorf("Add(%+v, %+v) = %+v, want %+v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestSum(t *testing.T) {
	tests := []struct {
		ms   []Money
		want Money
	}{
		{nil, Money{}},
		{[]Money{}, Money{}},
		{[]Money{{199}}, Money{199}},
		{[]Money{{100}, {250}, {1}}, Money{351}},
		{[]Money{{100}, {-100}, {-1}}, Money{-1}},
	}
	for _, tt := range tests {
		if got := Sum(tt.ms); got != tt.want {
			t.Errorf("Sum(%+v) = %+v, want %+v", tt.ms, got, tt.want)
		}
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		m    Money
		want string
	}{
		{Money{1205}, "$12.05"},
		{Money{7}, "$0.07"},
		{Money{}, "$0.00"},
		{Money{100}, "$1.00"},
		{Money{99}, "$0.99"},
		{Money{123456789}, "$1234567.89"},
		{Money{-50}, "-$0.50"},
		{Money{-1}, "-$0.01"},
		{Money{-12345}, "-$123.45"},
	}
	for _, tt := range tests {
		if got := Format(tt.m); got != tt.want {
			t.Errorf("Format(%+v) = %q, want %q", tt.m, got, tt.want)
		}
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		m    Money
		n    int
		want []Money
	}{
		{Money{100}, 3, []Money{{34}, {33}, {33}}},
		{Money{100}, 4, []Money{{25}, {25}, {25}, {25}}},
		{Money{5}, 1, []Money{{5}}},
		{Money{2}, 4, []Money{{1}, {1}, {0}, {0}}},
		{Money{0}, 2, []Money{{0}, {0}}},
		{Money{1001}, 3, []Money{{334}, {334}, {333}}},
		{Money{5}, 0, nil},
		{Money{5}, -2, nil},
	}
	for _, tt := range tests {
		got := Split(tt.m, tt.n)
		if !slices.Equal(got, tt.want) {
			t.Errorf("Split(%+v, %d) = %+v, want %+v", tt.m, tt.n, got, tt.want)
		}
		if tt.want == nil && got != nil {
			t.Errorf("Split(%+v, %d) = %#v, want nil", tt.m, tt.n, got)
		}
	}
}
