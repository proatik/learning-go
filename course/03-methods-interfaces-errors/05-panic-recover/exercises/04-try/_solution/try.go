package catch

import "fmt"

// toError formats a recovered panic value as "recovered: <value>".
func toError(r any) error {
	if e, ok := r.(error); ok {
		return fmt.Errorf("recovered: %w", e)
	}
	return fmt.Errorf("recovered: %v", r)
}

// Try calls f and converts a panic into the error "recovered: <value>",
// wrapping the value with %w when it is an error. It returns nil if f
// returns normally.
func Try(f func()) (err error) {
	defer func() {
		// Since Go 1.21 panic(nil) yields a *runtime.PanicNilError, so
		// r != nil reliably means "a panic happened".
		if r := recover(); r != nil {
			err = toError(r)
		}
	}()
	f()
	return nil
}

// TryInt is Try for a function with an int result; it returns 0 on panic.
func TryInt(f func() int) (n int, err error) {
	defer func() {
		if r := recover(); r != nil {
			n, err = 0, toError(r)
		}
	}()
	return f(), nil
}

// TryAll runs every job and returns one error (or nil) per job.
func TryAll(jobs []func()) []error {
	errs := make([]error, len(jobs))
	for i, job := range jobs {
		errs[i] = Try(job) // each Try recovers on its own, so the loop goes on
	}
	return errs
}
