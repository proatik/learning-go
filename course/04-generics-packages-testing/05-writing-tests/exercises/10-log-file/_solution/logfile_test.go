package logfile

import (
	"errors"
	"io/fs"
	"path/filepath"
	"slices"
	"testing"
)

// openLog opens a Log at path and closes it when the test ends.
func openLog(t *testing.T, path string) *Log {
	t.Helper()
	l, err := Open(path)
	if err != nil {
		t.Fatalf("Open(%q) error = %v, want nil", path, err)
	}
	t.Cleanup(func() { l.Close() })
	return l
}

// write writes each message to l and stops the test on the first error.
func write(t *testing.T, l *Log, msgs ...string) {
	t.Helper()
	for _, m := range msgs {
		if err := l.Write(m); err != nil {
			t.Fatalf("Write(%q) error = %v, want nil", m, err)
		}
	}
}

// readLines reads the log at path and compares it with want.
func readLines(t *testing.T, path string, want []string) {
	t.Helper()
	got, err := ReadLines(path)
	if err != nil {
		t.Fatalf("ReadLines() error = %v, want nil", err)
	}
	if !slices.Equal(got, want) {
		t.Errorf("ReadLines() = %q, want %q", got, want)
	}
}

// TestLog writes messages to a log in a temporary directory and reads them back.
func TestLog(t *testing.T) {
	t.Run("write and read", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "app.log")
		l := openLog(t, path)
		readLines(t, path, nil)
		write(t, l, "started", "listening on :8080")
		readLines(t, path, []string{"started", "listening on :8080"})
	})

	t.Run("reopen appends", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "app.log")
		l := openLog(t, path)
		write(t, l, "one", "two")
		l.Close()
		write(t, openLog(t, path), "three")
		readLines(t, path, []string{"one", "two", "three"})
	})

	t.Run("write after close", func(t *testing.T) {
		l := openLog(t, filepath.Join(t.TempDir(), "app.log"))
		l.Close()
		if err := l.Write("late"); !errors.Is(err, ErrClosed) {
			t.Errorf("Write after Close error = %v, want %v", err, ErrClosed)
		}
	})

	t.Run("missing file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "nope.log")
		if _, err := ReadLines(path); !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("ReadLines(missing) error = %v, want %v", err, fs.ErrNotExist)
		}
	})
}
