package bridge

import (
	"context"
	"errors"
)

// ErrClosed is returned by Read after Close was called.
var ErrClosed = errors.New("reader closed")

// Reader is a blocking source that knows nothing about contexts.
type Reader interface {
	Read() (string, error) // blocks until a message arrives; returns ErrClosed once Close is called
	Close() error
}

// ReadCtx performs one r.Read, closing r if ctx is done first so the read
// returns. It returns ctx.Err() if the context won.
func ReadCtx(ctx context.Context, r Reader) (string, error) {
	// TODO
	return "", nil
}

// Merge returns a child of a that is also canceled, with b's cause, when b is
// done. Its CancelFunc cancels it with cause context.Canceled.
func Merge(a, b context.Context) (context.Context, context.CancelFunc) {
	// TODO
	return a, func() {}
}
