package ipaddr

import "fmt"

// IPAddr is an IPv4 address.
type IPAddr [4]byte

// String returns the address in dotted decimal form, e.g. "127.0.0.1".
func (ip IPAddr) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

// Color is an RGB color.
type Color struct {
	R, G, B uint8
}

// String returns the color as "#rrggbb" with lowercase hex digits.
func (c Color) String() string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}
