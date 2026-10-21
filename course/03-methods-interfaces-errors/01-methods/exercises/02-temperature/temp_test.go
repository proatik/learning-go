package temp

import (
	"fmt"
	"math"
	"testing"
)

func near(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestToFahrenheit(t *testing.T) {
	tests := []struct {
		c    Celsius
		want Fahrenheit
	}{
		{100, 212},
		{0, 32},
		{-40, -40},
		{37, 98.6},
		{-273.15, -459.67},
	}
	for _, tt := range tests {
		if got := tt.c.ToFahrenheit(); !near(float64(got), float64(tt.want)) {
			t.Errorf("Celsius(%v).ToFahrenheit() = %v, want %v", float64(tt.c), float64(got), float64(tt.want))
		}
	}
}

func TestToCelsius(t *testing.T) {
	tests := []struct {
		f    Fahrenheit
		want Celsius
	}{
		{32, 0},
		{212, 100},
		{-40, -40},
		{98.6, 37},
		{0, -160.0 / 9},
	}
	for _, tt := range tests {
		if got := tt.f.ToCelsius(); !near(float64(got), float64(tt.want)) {
			t.Errorf("Fahrenheit(%v).ToCelsius() = %v, want %v", float64(tt.f), float64(got), float64(tt.want))
		}
	}
}

func TestRoundTrip(t *testing.T) {
	for _, c := range []Celsius{-12.5, 0, 21.5, 1000} {
		if got := c.ToFahrenheit().ToCelsius(); !near(float64(got), float64(c)) {
			t.Errorf("Celsius(%v).ToFahrenheit().ToCelsius() = %v, want %v", float64(c), float64(got), float64(c))
		}
	}
}

func TestIsFreezing(t *testing.T) {
	tests := []struct {
		c    Celsius
		want bool
	}{
		{0, true},
		{-5, true},
		{0.1, false},
		{25, false},
	}
	for _, tt := range tests {
		if got := tt.c.IsFreezing(); got != tt.want {
			t.Errorf("Celsius(%v).IsFreezing() = %v, want %v", float64(tt.c), got, tt.want)
		}
	}
}

func TestCelsiusString(t *testing.T) {
	tests := []struct {
		c    Celsius
		want string
	}{
		{21.5, "21.5°C"},
		{0, "0.0°C"},
		{-3.25, "-3.2°C"},
		{37, "37.0°C"},
		{100.06, "100.1°C"},
	}
	for _, tt := range tests {
		if got := tt.c.String(); got != tt.want {
			t.Errorf("Celsius(%v).String() = %q, want %q", float64(tt.c), got, tt.want)
		}
	}
}

func TestFahrenheitString(t *testing.T) {
	tests := []struct {
		f    Fahrenheit
		want string
	}{
		{70.7, "70.7°F"},
		{0, "0.0°F"},
		{-40, "-40.0°F"},
		{212, "212.0°F"},
	}
	for _, tt := range tests {
		if got := tt.f.String(); got != tt.want {
			t.Errorf("Fahrenheit(%v).String() = %q, want %q", float64(tt.f), got, tt.want)
		}
	}
}

func ExampleCelsius_String() {
	body := Celsius(37)
	fmt.Println(body)
	fmt.Println(body.ToFahrenheit())
	fmt.Printf("boiling: %v\n", Celsius(100).ToFahrenheit())
	// Output:
	// 37.0°C
	// 98.6°F
	// boiling: 212.0°F
}
