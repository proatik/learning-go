package first

import (
	"context"
	"errors"
)

// ErrNoFetchers is returned when First is called without fetchers.
var ErrNoFetchers = errors.New("first: no fetchers")

// First runs every fetcher concurrently and returns the first successful
// result, canceling the others. If all fail, it returns their errors joined in
// argument order.
func First(ctx context.Context, fetchers ...func(context.Context) (string, error)) (string, error) {
	// TODO
	return "", nil
}
