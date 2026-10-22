package fetch

import (
	"context"
	"fmt"
	"time"
)

// Fetch calls f with a child of ctx that times out after d (no extra timeout
// if d <= 0), cancels that child before returning, and wraps f's error as
// "fetch: <err>".
func Fetch(ctx context.Context, d time.Duration, f func(context.Context) (string, error)) (string, error) {
	var cancel context.CancelFunc
	if d > 0 {
		ctx, cancel = context.WithTimeout(ctx, d)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	s, err := f(ctx)
	if err != nil {
		return "", fmt.Errorf("fetch: %w", err)
	}
	return s, nil
}
