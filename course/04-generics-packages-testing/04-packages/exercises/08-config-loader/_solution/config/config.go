// Package config loads and validates the daemon's settings.
package config

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strconv"
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
	return Config{Host: "localhost", Port: 8080, Debug: false, Timeout: 5 * time.Second}
}

// Validate reports every invalid field, joined with errors.Join, or nil.
func (c Config) Validate() error {
	var problems []error
	if c.Host == "" {
		problems = append(problems, fmt.Errorf("%w: host is empty", ErrInvalid))
	}
	if c.Port < 1 || c.Port > 65535 {
		problems = append(problems, fmt.Errorf("%w: port %d out of range", ErrInvalid, c.Port))
	}
	if c.Timeout <= 0 {
		problems = append(problems, fmt.Errorf("%w: timeout must be positive", ErrInvalid))
	}
	return errors.Join(problems...)
}

// Load applies settings on top of Default in sorted key order, then
// validates the result. All problems are joined into one error.
func Load(settings map[string]string) (Config, error) {
	c := Default()
	var errs []error
	for _, key := range slices.Sorted(maps.Keys(settings)) {
		value := settings[key]
		var err error
		switch key {
		case "host":
			c.Host = value
		case "port":
			var n int
			if n, err = strconv.Atoi(value); err == nil {
				c.Port = n
			}
		case "debug":
			var b bool
			if b, err = strconv.ParseBool(value); err == nil {
				c.Debug = b
			}
		case "timeout":
			var d time.Duration
			if d, err = time.ParseDuration(value); err == nil {
				c.Timeout = d
			}
		default:
			errs = append(errs, fmt.Errorf("%w: %q", ErrUnknownKey, key))
			continue
		}
		if err != nil {
			errs = append(errs, fmt.Errorf("%w: %s=%q", ErrBadValue, key, value))
		}
	}
	errs = append(errs, c.Validate())
	if err := errors.Join(errs...); err != nil {
		return Config{}, err
	}
	return c, nil
}
