package kinds_test

import (
	"fmt"
	"testing"

	kinds "learngo/course/03-methods-interfaces-errors/02-interfaces/exercises/04-describe"
)

type celsius float64

func (c celsius) String() string { return fmt.Sprintf("%g°C", float64(c)) }

type point struct{ X, Y int }

type ticket struct{ ID int }

func (t *ticket) String() string { return fmt.Sprintf("ticket #%d", t.ID) }

func TestDescribe(t *testing.T) {
	var nilInts []int
	var nilPtr *int
	tests := []struct {
		in   any
		want string
	}{
		{nil, "nil"},
		{42, "int 42"},
		{0, "int 0"},
		{-3, "int -3"},
		{int8(-128), "sized int -128"},
		{int16(300), "sized int 300"},
		{'A', "sized int 65"}, // a rune is an int32
		{int64(-7), "sized int -7"},
		{uint(3), "unsigned 3"},
		{uint8(255), "unsigned 255"}, // a byte is a uint8
		{uint64(1 << 40), "unsigned 1099511627776"},
		{3.5, "float64 3.5"},
		{0.0, "float64 0"},
		{float32(1.5), "unknown float32"},
		{"héllo", `string "héllo" with 5 runes`},
		{"", `string "" with 0 runes`},
		{"a\tb", `string "a\tb" with 3 runes`},
		{true, "bool true"},
		{false, "bool false"},
		{[]int{1, 2, 3}, "[]int with 3 elements"},
		{nilInts, "[]int with 0 elements"},
		{[]string{"x"}, "unknown []string"},
		{map[string]int{"a": 1}, "map[string]int with 1 keys"},
		{map[string]int{}, "map[string]int with 0 keys"},
		{map[string]bool{}, "unknown map[string]bool"},
		{celsius(21), "stringer 21°C"},
		{&ticket{ID: 7}, "stringer ticket #7"},
		{ticket{ID: 7}, "unknown kinds_test.ticket"},
		{point{1, 2}, "unknown kinds_test.point"},
		{struct{}{}, "unknown struct {}"},
		{nilPtr, "unknown *int"},
	}
	for _, tt := range tests {
		if got := kinds.Describe(tt.in); got != tt.want {
			t.Errorf("Describe(%#v) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
