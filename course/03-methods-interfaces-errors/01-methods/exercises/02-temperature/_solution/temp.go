package temp

import "fmt"

// Celsius is a temperature in degrees Celsius.
type Celsius float64

// Fahrenheit is a temperature in degrees Fahrenheit.
type Fahrenheit float64

// ToFahrenheit converts c to Fahrenheit.
func (c Celsius) ToFahrenheit() Fahrenheit {
	return Fahrenheit(c*9/5 + 32)
}

// ToCelsius converts f to Celsius.
func (f Fahrenheit) ToCelsius() Celsius {
	return Celsius((f - 32) * 5 / 9)
}

// IsFreezing reports whether c is at or below 0°C.
func (c Celsius) IsFreezing() bool {
	return c <= 0
}

// String formats c like "21.5°C".
func (c Celsius) String() string {
	return fmt.Sprintf("%.1f°C", float64(c))
}

// String formats f like "70.7°F".
func (f Fahrenheit) String() string {
	return fmt.Sprintf("%.1f°F", float64(f))
}
