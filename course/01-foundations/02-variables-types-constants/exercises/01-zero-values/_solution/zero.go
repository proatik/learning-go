package zero

import "fmt"

// PrintZeroValues declares one variable of several basic types without
// initializing them and prints each variable's type and zero value.
func PrintZeroValues() {
	var (
		i  int
		i8 int8
		u  uint
		f  float64
		s  string
		ok bool
		r  rune // alias for int32
		b  byte // alias for uint8
	)

	fmt.Printf("%T = %v\n", i, i)
	fmt.Printf("%T = %v\n", i8, i8)
	fmt.Printf("%T = %v\n", u, u)
	fmt.Printf("%T = %v\n", f, f)
	fmt.Printf("%T = %q\n", s, s) // %q shows the empty string as ""
	fmt.Printf("%T = %v\n", ok, ok)
	fmt.Printf("%T = %v\n", r, r)
	fmt.Printf("%T = %v\n", b, b)
}
