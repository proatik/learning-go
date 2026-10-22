package fetch

import (
	"context"
	"errors"
	"sync"
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
	results := make([]Result, len(urls))
	sem := make(chan struct{}, max(limit, 1))
	var wg sync.WaitGroup
	for i, url := range urls {
		results[i].URL = url
		if !acquire(ctx, sem) {
			results[i].Err = ctx.Err()
			continue
		}
		wg.Go(func() {
			defer func() { <-sem }()
			results[i] = fetch(ctx, c, url, timeout)
		})
	}
	wg.Wait()
	return results
}

// acquire takes a slot in sem, or reports false if ctx is done first.
func acquire(ctx context.Context, sem chan struct{}) bool {
	if ctx.Err() != nil {
		return false
	}
	select {
	case sem <- struct{}{}:
	case <-ctx.Done():
		return false
	}
	// Both cases may have been ready; select picks randomly.
	if ctx.Err() != nil {
		<-sem
		return false
	}
	return true
}

// fetch downloads one URL, retrying once on a temporary failure.
func fetch(ctx context.Context, c Client, url string, timeout time.Duration) Result {
	r := Result{URL: url}
	for r.Attempts < 2 {
		r.Attempts++
		r.Body, r.Err = attempt(ctx, c, url, timeout)
		if !retryable(ctx, r.Err) {
			break
		}
	}
	return r
}

// attempt makes one Get call with its own timeout.
func attempt(ctx context.Context, c Client, url string, timeout time.Duration) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return c.Get(ctx, url)
}

// retryable reports whether err is worth a second attempt.
func retryable(ctx context.Context, err error) bool {
	if err == nil || ctx.Err() != nil {
		return false
	}
	return errors.Is(err, ErrTemporary) || errors.Is(err, context.DeadlineExceeded)
}
