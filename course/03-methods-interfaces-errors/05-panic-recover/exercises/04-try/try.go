package catch

// Try calls f and converts a panic into the error "recovered: <value>",
// wrapping the value with %w when it is an error. It returns nil if f
// returns normally.
func Try(f func()) (err error) {
	// TODO
	return nil
}

// TryInt is Try for a function with an int result; it returns 0 on panic.
func TryInt(f func() int) (n int, err error) {
	// TODO
	return 0, nil
}

// TryAll runs every job and returns one error (or nil) per job.
func TryAll(jobs []func()) []error {
	// TODO
	return nil
}
