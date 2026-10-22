package safego

// SafeRun runs every task in its own goroutine, waits for all of them, and
// returns one entry per task: nil, or an error describing its panic.
func SafeRun(tasks []func()) []error {
	// TODO
	return nil
}
