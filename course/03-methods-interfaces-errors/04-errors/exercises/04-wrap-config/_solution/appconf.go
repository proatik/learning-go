package appconf

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Config is a server address.
type Config struct {
	Host string
	Port int
}

// Sentinel errors returned (wrapped) by LoadConfig.
var (
	ErrNotFound    = errors.New("file not found")
	ErrMissingHost = errors.New("missing host")
	ErrMissingPort = errors.New("missing port")
	ErrPortRange   = errors.New("port out of range")
)

// parse parses "host:port". Its errors don't mention any file name.
func parse(content string) (Config, error) {
	host, p, ok := strings.Cut(strings.TrimSpace(content), ":")
	if !ok {
		return Config{}, ErrMissingPort
	}
	if host == "" {
		return Config{}, ErrMissingHost
	}
	n, err := strconv.Atoi(p)
	if err != nil {
		return Config{}, fmt.Errorf("port %q: %w", p, err)
	}
	if n < 1 || n > 65535 {
		return Config{}, fmt.Errorf("port %d: %w", n, ErrPortRange)
	}
	return Config{Host: host, Port: n}, nil
}

// LoadConfig parses the file called name from files, wrapping every error
// with `config "<name>": `.
func LoadConfig(files map[string]string, name string) (Config, error) {
	content, ok := files[name]
	if !ok {
		return Config{}, fmt.Errorf("config %q: %w", name, ErrNotFound)
	}
	cfg, err := parse(content)
	if err != nil {
		return Config{}, fmt.Errorf("config %q: %w", name, err)
	}
	return cfg, nil
}
