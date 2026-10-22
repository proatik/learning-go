// Package daemon starts a background service from its configuration.
package daemon

import "learngo/course/04-generics-packages-testing/04-packages/exercises/08-config-loader/config"

// Describe loads settings and returns "<host>:<port> debug=<bool> timeout=<duration>",
// or "" and the error from config.Load.
func Describe(settings map[string]string) (string, error) {
	_ = config.Load // TODO
	return "", nil
}
