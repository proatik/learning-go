package temp

// Celsius is a temperature in degrees Celsius.
type Celsius float64

// Fahrenheit is a temperature in degrees Fahrenheit.
type Fahrenheit float64

// ToFahrenheit converts c to Fahrenheit.
func (c Celsius) ToFahrenheit() Fahrenheit {
	// TODO
	return 0
}

// ToCelsius converts f to Celsius.
func (f Fahrenheit) ToCelsius() Celsius {
	// TODO
	return 0
}

// IsFreezing reports whether c is at or below 0°C.
func (c Celsius) IsFreezing() bool {
	// TODO
	return false
}

// String formats c like "21.5°C".
func (c Celsius) String() string {
	// TODO
	return ""
}

// String formats f like "70.7°F".
func (f Fahrenheit) String() string {
	// TODO
	return ""
}
