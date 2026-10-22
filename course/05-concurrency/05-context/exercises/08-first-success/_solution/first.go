package first

import (
	"context"
	"errors"
)

// ErrNoFetchers is returned when First is called without fetchers.
var ErrNoFetchers = errors.New("first: no fetchers")

type result struct {
	i   int
	val string
	err error
}

// First runs every fetcher concurrently and returns the first successful
// result, canceling the others. If all fail, it returns their errors joined in
// argument order.
func First(ctx context.Context, fetchers ...func(context.Context) (string, error)) (string, error) {
	if len(fetchers) == 0 {
		return "", ErrNoFetchers
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make(chan result, len(fetchers)) // room for all: late losers never block
	for i, f := range fetchers {
		go func() {
			v, err := f(ctx)
			results <- result{i, v, err}
		}()
	}

	errs := make([]error, len(fetchers))
	for range fetchers {
		r := <-results
		if r.err == nil {
			return r.val, nil
		}
		errs[r.i] = r.err
	}
	return "", errors.Join(errs...)
}
