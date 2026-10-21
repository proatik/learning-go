package patterns

import "fmt"

// PrintTriangle prints a left-aligned right triangle of n rows, where row i
// contains i stars.
func PrintTriangle(n int) {
	for i := 1; i <= n; i++ { // one iteration per row
		for range i { // i stars on row i
			fmt.Print("*")
		}
		fmt.Println()
	}
}

// PrintPyramid prints a centred pyramid of n rows, where row i contains n-i
// spaces followed by 2i-1 stars.
func PrintPyramid(n int) {
	for i := 1; i <= n; i++ {
		for range n - i { // leading spaces only: nothing after the stars
			fmt.Print(" ")
		}
		for range 2*i - 1 {
			fmt.Print("*")
		}
		fmt.Println()
	}
}

// PrintTable prints an n×n multiplication table with every cell
// right-aligned in a 4-character field.
func PrintTable(n int) {
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			fmt.Printf("%4d", i*j) // the padding goes on the left, so rows never end in a space
		}
		fmt.Println()
	}
}
