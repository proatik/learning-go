package bytesize

import "fmt"

// ExampleFormatBytes shows FormatBytes for sizes in every unit.
func ExampleFormatBytes() {
	fmt.Println(FormatBytes(512))
	// TODO: print more sizes, one per line. Pick sizes near each unit edge
	//       and some that need the decimal (1536 is 1.5 KiB).
	// TODO: end the function with an output comment: a line "// Output:"
	//       followed by one "// line" per printed line, exactly as printed.
}
