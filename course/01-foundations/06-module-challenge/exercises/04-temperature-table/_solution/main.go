package main

import (
	"fmt"
	"strings"
)

// Celsius, Fahrenheit and Kelvin are temperatures on the three scales.
// Distinct named types stop you from mixing them up by accident.
type (
	Celsius    float64
	Fahrenheit float64
	Kelvin     float64
)

// Table range, in degrees Celsius.
const (
	start = -40
	end   = 100
	step  = 20
)

// Column widths.
const (
	wC = 7
	wF = 10
	wK = 7
)

// CToF converts Celsius to Fahrenheit.
func CToF(c Celsius) Fahrenheit { return Fahrenheit(c*9/5 + 32) }

// CToK converts Celsius to Kelvin.
func CToK(c Celsius) Kelvin { return Kelvin(c + 273.15) }

func main() {
	// %*s takes the width from the argument list, so the widths live in
	// one place.
	fmt.Printf("%*s  %*s  %*s\n", wC, "Celsius", wF, "Fahrenheit", wK, "Kelvin")
	fmt.Println(strings.Repeat("-", wC) + "  " + strings.Repeat("-", wF) + "  " + strings.Repeat("-", wK))

	rows := 0
	for c := Celsius(start); c <= end; c += step {
		fmt.Printf("%*.1f  %*.1f  %*.2f\n", wC, c, wF, CToF(c), wK, CToK(c))
		rows++
	}
	fmt.Printf("%d rows\n", rows)
}
