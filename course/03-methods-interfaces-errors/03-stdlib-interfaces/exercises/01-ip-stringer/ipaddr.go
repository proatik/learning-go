package ipaddr

// IPAddr is an IPv4 address.
type IPAddr [4]byte

// String returns the address in dotted decimal form, e.g. "127.0.0.1".
func (ip IPAddr) String() string {
	// TODO
	return ""
}

// Color is an RGB color.
type Color struct {
	R, G, B uint8
}

// String returns the color as "#rrggbb" with lowercase hex digits.
func (c Color) String() string {
	// TODO
	return ""
}
