// Package daemon starts a background service from its configuration.
package daemon

import (
	"fmt"

	"learngo/course/04-generics-packages-testing/04-packages/exercises/08-config-loader/config"
)

// Describe loads settings and returns "<host>:<port> debug=<bool> timeout=<duration>",
// or "" and the error from config.Load.
func Describe(settings map[string]string) (string, error) {
	c, err := config.Load(settings)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d debug=%t timeout=%v", c.Host, c.Port, c.Debug, c.Timeout), nil
}
