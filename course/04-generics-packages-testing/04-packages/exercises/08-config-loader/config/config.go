// Package config loads and validates the daemon's settings.
package config

import (
	"errors"
	"time"
)

// Config holds the daemon's settings. Every combination of field values can
// be represented, so the fields are exported; Validate reports bad ones.
type Config struct {
	Host    string
	Port    int
	Debug   bool
	Timeout time.Duration
}

// Errors wrapped by Load and Validate.
var (
	ErrUnknownKey = errors.New("config: unknown key")
	ErrBadValue   = errors.New("config: bad value")
	ErrInvalid    = errors.New("config: invalid")
)

// Default returns the default settings: localhost:8080, no debug, 5s timeout.
func Default() Config {
	// TODO
	return Config{}
}

// Validate reports every invalid field, joined with errors.Join, or nil.
func (c Config) Validate() error {
	// TODO: collect problems in a []error, then return errors.Join(problems...)
	return nil
}

// Load applies settings on top of Default in sorted key order, then
// validates the result. All problems are joined into one error.
func Load(settings map[string]string) (Config, error) {
	// TODO: switch on the key: "host", "port", "debug", "timeout", default → unknown
	return Config{}, nil
}
