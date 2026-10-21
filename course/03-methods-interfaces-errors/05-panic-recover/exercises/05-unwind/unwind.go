package unwind

import "io"

// Unwind prints "enter i" for nested levels 1..depth, panics at the deepest
// level, prints "exit i" from deferred calls while unwinding, and finally
// "recovered: bottom at <depth>". For depth <= 0 it prints "nothing to unwind".
func Unwind(w io.Writer, depth int) {
	// TODO
}

// Acquire prints "acquire <name>" for each name and defers "release <name>".
// It panics with "failed at <name>" right after acquiring names[failAt],
// recovers, and prints "error: <value>" last. Without a failure it prints
// "all done" before the releases.
func Acquire(w io.Writer, names []string, failAt int) {
	// TODO
}
