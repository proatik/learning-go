package profile

import "fmt"

// PrintProfile prints a three-line profile card.
func PrintProfile(name string, year int, active bool) {
	fmt.Printf("Name:   %s\n", name)
	fmt.Printf("Born:   %d\n", year)
	fmt.Printf("Active: %t\n", active)
}
