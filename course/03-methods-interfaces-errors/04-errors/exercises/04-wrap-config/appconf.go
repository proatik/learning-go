package appconf

import "errors"

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
	// TODO
	return Config{}, nil
}

// LoadConfig parses the file called name from files, wrapping every error
// with `config "<name>": `.
func LoadConfig(files map[string]string, name string) (Config, error) {
	// TODO
	return Config{}, nil
}
