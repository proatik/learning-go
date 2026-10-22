package bytesize

import "fmt"

// ExampleFormatBytes shows FormatBytes for sizes in every unit.
func ExampleFormatBytes() {
	fmt.Println(FormatBytes(512))
	fmt.Println(FormatBytes(1000))
	fmt.Println(FormatBytes(1023))
	fmt.Println(FormatBytes(1024))
	fmt.Println(FormatBytes(1536))
	fmt.Println(FormatBytes(5 * 1024 * 1024))
	fmt.Println(FormatBytes(3 << 30))
	fmt.Println(FormatBytes(2 << 40))
	// Output:
	// 512 B
	// 1000 B
	// 1023 B
	// 1.0 KiB
	// 1.5 KiB
	// 5.0 MiB
	// 3.0 GiB
	// 2.0 TiB
}
