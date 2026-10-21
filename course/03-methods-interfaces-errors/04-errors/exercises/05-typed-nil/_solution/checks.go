package checks

import (
	"errors"
	"strings"
)

// CheckError describes one failed check.
type CheckError struct {
	Name    string
	Problem string
}

// Error returns "<name>: <problem>".
func (e *CheckError) Error() string {
	return e.Name + ": " + e.Problem
}

// CheckHost reports an empty host.
func CheckHost(host string) error {
	if strings.TrimSpace(host) == "" {
		return &CheckError{Name: "host", Problem: "must not be empty"}
	}
	return nil
}

// checkPort reports a port outside 1..65535. It returns the error
// interface, so a success is a true nil.
func checkPort(port int) error {
	if port < 1 || port > 65535 {
		return &CheckError{Name: "port", Problem: "must be 1-65535"}
	}
	return nil
}

// CheckPort reports a port outside 1..65535.
func CheckPort(port int) error {
	return checkPort(port)
}

// CheckAll runs every check and joins the failures (nil if all pass).
func CheckAll(host string, port int) error {
	var errs []error
	if err := CheckHost(host); err != nil {
		errs = append(errs, err)
	}
	if err := CheckPort(port); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

// FirstFailure returns the first non-nil result, or nil if there is none.
func FirstFailure(results []*CheckError) error {
	for _, r := range results {
		if r != nil {
			return r
		}
	}
	return nil
}
