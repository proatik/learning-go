package store

// Store is a destination for lines that must be closed after use.
type Store interface {
	Write(line string) error
	Close() error
}

// Save opens a store, writes lines to it in order and always closes it.
// The first error wins: open, then write, then close.
func Save(open func() (Store, error), lines []string) (err error) {
	// TODO
	return nil
}

// SaveBatches saves each batch to its own store, closing each store before
// opening the next. It stops at the first failure ("batch <i>: ...").
func SaveBatches(open func(i int) (Store, error), batches [][]string) error {
	// TODO
	return nil
}
