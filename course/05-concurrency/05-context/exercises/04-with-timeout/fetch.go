package fetch

import (
	"context"
	"time"
)

// Fetch calls f with a child of ctx that times out after d (no extra timeout
// if d <= 0), cancels that child before returning, and wraps f's error as
// "fetch: <err>".
func Fetch(ctx context.Context, d time.Duration, f func(context.Context) (string, error)) (string, error) {
	// TODO
	return "", nil
}
