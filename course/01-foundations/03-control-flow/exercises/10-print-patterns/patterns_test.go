package patterns

import "fmt"

// Each example prints "---" before and after the pattern so that leading
// spaces on the first row and trailing spaces on the last row are checked too
// (the Output comparison trims whitespace around the whole output).

func ExamplePrintTriangle() {
	fmt.Println("---")
	PrintTriangle(4)
	fmt.Println("---")
	// Output:
	// ---
	// *
	// **
	// ***
	// ****
	// ---
}

func ExamplePrintTriangle_one() {
	fmt.Println("---")
	PrintTriangle(1)
	fmt.Println("---")
	// Output:
	// ---
	// *
	// ---
}

func ExamplePrintPyramid() {
	fmt.Println("---")
	PrintPyramid(4)
	fmt.Println("---")
	// Output:
	// ---
	//    *
	//   ***
	//  *****
	// *******
	// ---
}

func ExamplePrintPyramid_two() {
	fmt.Println("---")
	PrintPyramid(2)
	PrintPyramid(1)
	fmt.Println("---")
	// Output:
	// ---
	//  *
	// ***
	// *
	// ---
}

func ExamplePrintTable() {
	fmt.Println("---")
	PrintTable(4)
	fmt.Println("---")
	// Output:
	// ---
	//    1   2   3   4
	//    2   4   6   8
	//    3   6   9  12
	//    4   8  12  16
	// ---
}

func ExamplePrintTable_nine() {
	fmt.Println("---")
	PrintTable(9)
	fmt.Println("---")
	// Output:
	// ---
	//    1   2   3   4   5   6   7   8   9
	//    2   4   6   8  10  12  14  16  18
	//    3   6   9  12  15  18  21  24  27
	//    4   8  12  16  20  24  28  32  36
	//    5  10  15  20  25  30  35  40  45
	//    6  12  18  24  30  36  42  48  54
	//    7  14  21  28  35  42  49  56  63
	//    8  16  24  32  40  48  56  64  72
	//    9  18  27  36  45  54  63  72  81
	// ---
}

func ExamplePrintTable_one() {
	fmt.Println("---")
	PrintTable(1)
	fmt.Println("---")
	// Output:
	// ---
	//    1
	// ---
}

func Example_nothingForSmallN() {
	fmt.Println("---")
	PrintTriangle(0)
	PrintPyramid(0)
	PrintTable(0)
	PrintTriangle(-3)
	PrintPyramid(-3)
	PrintTable(-3)
	fmt.Println("---")
	// Output:
	// ---
	// ---
}
