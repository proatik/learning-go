package main

// Celsius, Fahrenheit and Kelvin are temperatures on the three scales.
// Distinct named types stop you from mixing them up by accident.
type (
	Celsius    float64
	Fahrenheit float64
	Kelvin     float64
)

// CToF converts Celsius to Fahrenheit.
func CToF(c Celsius) Fahrenheit {
	// TODO
	return 0
}

// CToK converts Celsius to Kelvin.
func CToK(c Celsius) Kelvin {
	// TODO
	return 0
}

func main() {
	// TODO: print the header, the dashes line, one row per temperature
	// from -40 to 100 in steps of 20, and the "8 rows" footer.
}
