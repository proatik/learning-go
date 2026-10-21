package temperature

import "fmt"

// Describe formats a Celsius temperature together with its Fahrenheit value.
func Describe(c float64) string {
	return fmt.Sprintf("%.1f°C = %.1f°F", c, CToF(c))
}
