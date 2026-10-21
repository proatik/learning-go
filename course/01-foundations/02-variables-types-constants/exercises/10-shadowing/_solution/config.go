package config

import (
	"fmt"
	"strconv"
)

// DefaultPort is used when no port is configured.
const DefaultPort = 8080

// ParsePort parses a TCP port number. An empty string means DefaultPort.
// Valid ports are 1 to 65535.
func ParsePort(s string) (int, error) {
	port := DefaultPort
	if s != "" {
		// Declare err up front so plain assignment (=) can be used: the
		// result now goes into the outer port instead of a new, shadowing one.
		var err error
		port, err = strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		if port < 1 || port > 65535 {
			return 0, fmt.Errorf("port %d out of range 1-65535", port)
		}
	}
	return port, nil
}

// Scale returns value multiplied by the integer in factor. An empty factor
// means 1. A factor of zero is an error.
func Scale(value int, factor string) (int, error) {
	f := 1
	if factor != "" {
		var err error
		f, err = strconv.Atoi(factor) // assign to the outer f
		if err != nil {
			return 0, err
		}
		if f == 0 {
			return 0, fmt.Errorf("scale factor must not be zero")
		}
	}
	return value * f, nil
}
