package unwind

import (
	"fmt"
	"io"
)

// Unwind prints "enter i" for nested levels 1..depth, panics at the deepest
// level, prints "exit i" from deferred calls while unwinding, and finally
// "recovered: bottom at <depth>". For depth <= 0 it prints "nothing to unwind".
func Unwind(w io.Writer, depth int) {
	if depth <= 0 {
		fmt.Fprintln(w, "nothing to unwind")
		return
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(w, "recovered:", r)
		}
	}()
	level(w, 1, depth)
}

func level(w io.Writer, i, depth int) {
	fmt.Fprintln(w, "enter", i)
	defer fmt.Fprintln(w, "exit", i) // arguments evaluated now, printed while unwinding
	if i == depth {
		panic(fmt.Sprintf("bottom at %d", depth))
	}
	level(w, i+1, depth)
}

// Acquire prints "acquire <name>" for each name and defers "release <name>".
// It panics with "failed at <name>" right after acquiring names[failAt],
// recovers, and prints "error: <value>" last. Without a failure it prints
// "all done" before the releases.
func Acquire(w io.Writer, names []string, failAt int) {
	// Deferred first, so it runs last: after every release.
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(w, "error:", r)
		}
	}()
	for i, name := range names {
		fmt.Fprintln(w, "acquire", name)
		defer fmt.Fprintln(w, "release", name)
		if i == failAt {
			panic("failed at " + name)
		}
	}
	fmt.Fprintln(w, "all done")
}
