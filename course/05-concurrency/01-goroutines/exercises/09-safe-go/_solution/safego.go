package safego

import (
	"fmt"
	"sync"
)

// SafeRun runs every task in its own goroutine, waits for all of them, and
// returns one entry per task: nil, or an error describing its panic.
func SafeRun(tasks []func()) []error {
	errs := make([]error, len(tasks))
	var wg sync.WaitGroup
	for i, task := range tasks {
		wg.Go(func() { errs[i] = protect(i, task) })
	}
	wg.Wait()
	return errs
}

// protect calls task and turns a panic into an error. It must run in the
// same goroutine as task: recover only sees panics of its own goroutine.
func protect(i int, task func()) (err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		if e, ok := r.(error); ok {
			err = fmt.Errorf("task %d panicked: %w", i, e)
		} else {
			err = fmt.Errorf("task %d panicked: %v", i, r)
		}
	}()
	task()
	return nil
}
