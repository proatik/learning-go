package collect

import (
	"errors"
	"sync"
)

// RunAll runs every task concurrently, waits for all of them and returns
// their errors joined in input order, or nil if none failed.
func RunAll(tasks []func() error) error {
	errs := make([]error, len(tasks))
	var wg sync.WaitGroup
	for i, task := range tasks {
		wg.Go(func() { errs[i] = task() })
	}
	wg.Wait()
	return errors.Join(errs...)
}
