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
		// BUG: this line has a problem. Find and fix it.
		port, err := strconv.Atoi(s)
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
		f, err := strconv.Atoi(factor)
		if err != nil {
			return 0, err
		}
		if f == 0 {
			return 0, fmt.Errorf("scale factor must not be zero")
		}
	}
	return value * f, nil
}
