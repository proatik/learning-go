package port

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// The errors ParsePort can return, always wrapped with the input.
var (
	ErrEmpty      = errors.New("empty port")
	ErrNotNumber  = errors.New("port is not a number")
	ErrOutOfRange = errors.New("port out of range")
)

// ParsePort parses a TCP port from 1 to 65535. During the bug hunt it is
// swapped for a broken version, and your tests must notice each one.
func ParsePort(s string) (int, error) { return impl(s) }

var impl = parsePort

func parsePort(s string) (int, error) {
	if s == "" {
		return 0, fmt.Errorf("port %q: %w", s, ErrEmpty)
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, fmt.Errorf("port %q: %w", s, ErrNotNumber)
		}
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 || n > 65535 {
		return 0, fmt.Errorf("port %q: %w", s, ErrOutOfRange)
	}
	return n, nil
}

var bugs = map[string]func(string) (int, error){
	"accepts port 0": func(s string) (int, error) {
		if strings.Trim(s, "0") == "" && s != "" {
			return 0, nil
		}
		return parsePort(s)
	},
	"rejects 65535": func(s string) (int, error) {
		if n, err := parsePort(s); err != nil || n < 65535 {
			return n, err
		}
		return 0, fmt.Errorf("port %q: %w", s, ErrOutOfRange)
	},
	"letters report ErrOutOfRange": func(s string) (int, error) {
		n, err := parsePort(s)
		if errors.Is(err, ErrNotNumber) {
			return 0, fmt.Errorf("port %q: %w", s, ErrOutOfRange)
		}
		return n, err
	},
	"empty input is not an error": func(s string) (int, error) {
		if s == "" {
			return 0, nil
		}
		return parsePort(s)
	},
	"accepts a leading plus": func(s string) (int, error) {
		return parsePort(strings.TrimPrefix(s, "+"))
	},
	"returns the number with the error": func(s string) (int, error) {
		n, err := parsePort(s)
		if errors.Is(err, ErrOutOfRange) {
			n, _ = strconv.Atoi(s)
		}
		return n, err
	},
}

const bugEnv = "LEARNGO_BUG"

func init() {
	if b := os.Getenv(bugEnv); b != "" {
		impl = bugs[b]
	}
}

// TestYourTestsCatchEveryBug re-runs your tests once per bug: each run must fail.
func TestYourTestsCatchEveryBug(t *testing.T) {
	if os.Getenv(bugEnv) != "" {
		t.Skip("inside the bug hunt")
	}
	run := func(bug string, args ...string) (string, error) {
		cmd := exec.Command(os.Args[0], append([]string{"-test.count=1", "-test.skip=^TestYourTestsCatchEveryBug$"}, args...)...)
		cmd.Env = append(os.Environ(), bugEnv+"="+bug)
		out, err := cmd.CombinedOutput()
		return string(out), err
	}
	list, _ := run("", "-test.list=^Test")
	if !slices.ContainsFunc(strings.Fields(list), func(n string) bool { return n != "TestYourTestsCatchEveryBug" && strings.HasPrefix(n, "Test") }) {
		t.Fatal("no Test functions found: write at least one test")
	}
	if out, err := run("", "-test.run=^Test"); err != nil {
		t.Fatalf("your tests fail on the correct ParsePort:\n%s", out)
	}
	for _, name := range slices.Sorted(maps.Keys(bugs)) {
		if _, err := run(name, "-test.run=^Test"); err == nil {
			t.Errorf("Bug %q = not caught, want caught", name)
		}
	}
}
