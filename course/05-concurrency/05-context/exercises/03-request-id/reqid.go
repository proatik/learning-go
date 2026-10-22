package reqid

import "context"

// ctxKey is the context key for request IDs. It is unexported, so no other
// package can build the same key.
type ctxKey struct{}

// WithRequestID returns a copy of ctx that carries id.
func WithRequestID(ctx context.Context, id string) context.Context {
	// TODO
	return ctx
}

// RequestID returns the request ID stored in ctx, and whether there was one.
func RequestID(ctx context.Context) (string, bool) {
	// TODO
	return "", false
}

// Detach returns a context with ctx's values that is never canceled and has
// no deadline, for work that must outlive the request.
func Detach(ctx context.Context) context.Context {
	// TODO
	return ctx
}
