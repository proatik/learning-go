package convert

// BMI returns the body-mass index for a weight in kilograms and a height in
// centimetres: weight divided by the square of the height in metres.
func BMI(weightKg float64, heightCm int) float64 {
	// Convert first, then divide: float64(heightCm) / 100 keeps the fraction.
	m := float64(heightCm) / 100
	return weightKg / (m * m)
}

// Truncate returns f without its fractional part (rounding toward zero).
func Truncate(f float64) int {
	return int(f) // float → int conversion truncates toward zero
}

// Pixels returns how many whole pixels span the given length in inches at
// dpi dots per inch.
func Pixels(inches float64, dpi int) int {
	return int(inches * float64(dpi))
}
