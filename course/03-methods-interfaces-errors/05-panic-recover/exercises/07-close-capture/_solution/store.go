package store

import "fmt"

// Store is a destination for lines that must be closed after use.
type Store interface {
	Write(line string) error
	Close() error
}

// Save opens a store, writes lines to it in order and always closes it.
// The first error wins: open, then write, then close.
func Save(open func() (Store, error), lines []string) (err error) {
	s, err := open()
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer func() {
		// Runs on every exit path, including a panic in Write.
		if cerr := s.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("close: %w", cerr)
		}
	}()
	for _, line := range lines {
		if err := s.Write(line); err != nil {
			return fmt.Errorf("write %q: %w", line, err)
		}
	}
	return nil
}

// SaveBatches saves each batch to its own store, closing each store before
// opening the next. It stops at the first failure ("batch <i>: ...").
func SaveBatches(open func(i int) (Store, error), batches [][]string) error {
	for i, lines := range batches {
		// Save closes its store before returning: no defers pile up here.
		if err := Save(func() (Store, error) { return open(i) }, lines); err != nil {
			return fmt.Errorf("batch %d: %w", i, err)
		}
	}
	return nil
}
