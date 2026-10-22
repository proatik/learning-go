package logfile

import "testing"

// openLog opens a Log at path and closes it when the test ends.
func openLog(t *testing.T, path string) *Log {
	// TODO: mark this as a helper.
	// TODO: call Open; on error stop the test with t.Fatalf.
	// TODO: register l.Close with t.Cleanup, so the file is closed even if
	//       the test fails (on Windows an open file blocks TempDir removal).
	return nil
}

// TestLog writes messages to a log in a temporary directory and reads them back.
func TestLog(t *testing.T) {
	// TODO: path := filepath.Join(t.TempDir(), "app.log")
	// TODO: write two messages, then check ReadLines returns exactly those two.
	// TODO: close, reopen, write one more: all three lines must be there.
	// TODO: Write after Close must return ErrClosed (errors.Is).
	// TODO: ReadLines on a missing file must wrap fs.ErrNotExist.
	_ = openLog
	_ = ReadLines
}
