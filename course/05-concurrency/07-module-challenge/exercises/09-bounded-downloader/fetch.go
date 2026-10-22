package fetch

import (
	"context"
	"errors"
	"time"
)

// Client fetches one URL.
type Client interface {
	Get(ctx context.Context, url string) ([]byte, error)
}

// ErrTemporary marks failures worth one retry.
var ErrTemporary = errors.New("fetch: temporary failure")

// Result is the outcome of downloading one URL.
type Result struct {
	URL      string
	Body     []byte
	Err      error
	Attempts int // Get calls made for this URL: 0, 1 or 2
}

// Download fetches every URL with at most limit requests in flight, a
// timeout per attempt and one retry for temporary failures. Results are in
// the same order as urls.
func Download(ctx context.Context, c Client, urls []string, limit int, timeout time.Duration) []Result {
	// TODO
	return nil
}
