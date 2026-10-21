package main

import (
	"math"
	"testing"
)

func TestCToF(t *testing.T) {
	tests := []struct {
		c    Celsius
		want Fahrenheit
	}{
		{-40, -40},
		{0, 32},
		{37, 98.6},
		{100, 212},
		{-273.15, -459.67},
		{21.5, 70.7},
	}
	for _, tt := range tests {
		if got := CToF(tt.c); math.Abs(float64(got-tt.want)) > 1e-9 {
			t.Errorf("CToF(%v) = %v, want %v", tt.c, got, tt.want)
		}
	}
}

func TestCToK(t *testing.T) {
	tests := []struct {
		c    Celsius
		want Kelvin
	}{
		{-273.15, 0},
		{0, 273.15},
		{100, 373.15},
		{-40, 233.15},
		{26.85, 300},
	}
	for _, tt := range tests {
		if got := CToK(tt.c); math.Abs(float64(got-tt.want)) > 1e-9 {
			t.Errorf("CToK(%v) = %v, want %v", tt.c, got, tt.want)
		}
	}
}

// Example runs main and compares its output with the comment below.
func Example() {
	main()
	// Output:
	// Celsius  Fahrenheit   Kelvin
	// -------  ----------  -------
	//   -40.0       -40.0   233.15
	//   -20.0        -4.0   253.15
	//     0.0        32.0   273.15
	//    20.0        68.0   293.15
	//    40.0       104.0   313.15
	//    60.0       140.0   333.15
	//    80.0       176.0   353.15
	//   100.0       212.0   373.15
	// 8 rows
}
