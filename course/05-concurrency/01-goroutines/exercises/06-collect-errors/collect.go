package collect

// RunAll runs every task concurrently, waits for all of them and returns
// their errors joined in input order, or nil if none failed.
func RunAll(tasks []func() error) error {
	// TODO
	return nil
}
