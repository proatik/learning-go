package logfile

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"slices"
	"strings"
	"testing"
)

// ErrClosed is returned by Write after Close.
var ErrClosed = errors.New("log is closed")

// Log appends one line per message to a file. During the bug hunt one of the
// named bugs below is switched on, and your tests must notice each one.
type Log struct {
	f      *os.File
	closed bool
}

// Open opens the log at path for appending, creating the file if needed.
func Open(path string) (*Log, error) {
	flags := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	if bug == "Open truncates the file" {
		flags = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	}
	f, err := os.OpenFile(path, flags, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open log: %w", err)
	}
	return &Log{f: f}, nil
}

// Write appends msg and a newline. After Close it returns ErrClosed.
func (l *Log) Write(msg string) error {
	if l.closed {
		if bug == "Write after Close succeeds" {
			return nil
		}
		return ErrClosed
	}
	if bug != "Write forgets the newline" {
		msg += "\n"
	}
	_, err := l.f.WriteString(msg)
	return err
}

// Close closes the file. Closing twice is allowed.
func (l *Log) Close() error {
	if l.closed {
		return nil
	}
	l.closed = true
	return l.f.Close()
}

// ReadLines returns the lines of the log at path, without their newlines.
// An empty file has no lines. A missing file is an error that wraps
// fs.ErrNotExist.
func ReadLines(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if bug == "ReadLines hides a missing file" {
			return nil, nil
		}
		return nil, fmt.Errorf("read log: %w", err)
	}
	s := string(data)
	if bug != "ReadLines keeps an empty last line" {
		s = strings.TrimSuffix(s, "\n")
	}
	if s == "" {
		return nil, nil
	}
	lines := strings.Split(s, "\n")
	if bug == "ReadLines drops the last line" {
		lines = lines[:len(lines)-1]
	}
	return lines, nil
}

var bugs = []string{
	"Open truncates the file",
	"Write forgets the newline",
	"Write after Close succeeds",
	"ReadLines keeps an empty last line",
	"ReadLines drops the last line",
	"ReadLines hides a missing file",
}

const (
	bugEnv    = "LEARNGO_BUG"
	learnFile = "logfile_test.go"
)

var bug = os.Getenv(bugEnv)

// calls returns the names of the methods your file calls, like "TempDir".
func calls(t *testing.T) map[string]bool {
	t.Helper()
	f, err := parser.ParseFile(token.NewFileSet(), learnFile, nil, 0)
	if err != nil {
		t.Fatalf("parse %s: %v", learnFile, err)
	}
	seen := map[string]bool{}
	ast.Inspect(f, func(n ast.Node) bool {
		if c, ok := n.(*ast.CallExpr); ok {
			if s, ok := c.Fun.(*ast.SelectorExpr); ok {
				seen[s.Sel.Name] = true
			}
		}
		return true
	})
	return seen
}

// dirEntries lists the exercise directory, to spot files your tests leave behind.
func dirEntries(t *testing.T) []string {
	t.Helper()
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return names
}

// TestYourTestsCatchEveryBug re-runs your tests once per bug: each run must
// fail. Your tests must use t.TempDir and t.Cleanup, and leave no files behind.
func TestYourTestsCatchEveryBug(t *testing.T) {
	if bug != "" {
		t.Skip("inside the bug hunt")
	}
	used := calls(t)
	for _, m := range []string{"TempDir", "Cleanup"} {
		if !used[m] {
			t.Errorf("calls to t.%s = 0, want at least 1", m)
		}
	}
	before := dirEntries(t)
	t.Cleanup(func() {
		if after := dirEntries(t); !slices.Equal(before, after) {
			t.Errorf("files in the exercise folder = %v, want %v (write only inside t.TempDir())", after, before)
		}
	})
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
		t.Fatalf("your tests fail on the correct Log:\n%s", out)
	}
	for _, name := range bugs {
		if _, err := run(name, "-test.run=^Test"); err == nil {
			t.Errorf("Bug %q = not caught, want caught", name)
		}
	}
}
