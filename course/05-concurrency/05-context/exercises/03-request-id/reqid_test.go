package reqid

import (
	"context"
	"testing"
	"time"
)

func TestRoundTrip(t *testing.T) {
	for _, id := range []string{"req-42", "", "héllo 世界", "a b c"} {
		ctx := WithRequestID(context.Background(), id)
		if ctx == nil {
			t.Fatalf("WithRequestID(ctx, %q) returned a nil context", id)
		}
		if got, ok := RequestID(ctx); got != id || !ok {
			t.Errorf("RequestID(WithRequestID(ctx, %q)) = (%q, %v), want (%q, true)", id, got, ok, id)
		}
	}
}

func TestMissing(t *testing.T) {
	if got, ok := RequestID(context.Background()); got != "" || ok {
		t.Errorf("RequestID(Background) = (%q, %v), want (\"\", false)", got, ok)
	}
	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), "other", 1))
	defer cancel()
	if got, ok := RequestID(ctx); got != "" || ok {
		t.Errorf("RequestID(ctx without an ID) = (%q, %v), want (\"\", false)", got, ok)
	}
}

func TestUsesUnexportedKey(t *testing.T) {
	ctx := WithRequestID(context.Background(), "req-1")
	if v, _ := ctx.Value(ctxKey{}).(string); v != "req-1" {
		t.Fatalf("ctx.Value(ctxKey{}) = %v, want \"req-1\" (store the ID under ctxKey{})", ctx.Value(ctxKey{}))
	}
	for _, k := range []any{"requestID", "request-id", "requestId", "id", "reqid", "ctxKey", 0, 1} {
		if v := ctx.Value(k); v != nil {
			t.Errorf("ctx.Value(%#v) = %v, want <nil> (a built-in key type can collide with other packages)", k, v)
		}
	}
}

func TestNoCollision(t *testing.T) {
	for _, k := range []string{"requestID", "request-id", "ctxKey", ""} {
		ctx := context.WithValue(context.Background(), k, "evil")
		if got, ok := RequestID(ctx); got != "" || ok {
			t.Errorf("RequestID(WithValue(ctx, %q, \"evil\")) = (%q, %v), want (\"\", false)", k, got, ok)
		}
		ctx = WithRequestID(ctx, "good")
		ctx = context.WithValue(ctx, k, "evil")
		if got, ok := RequestID(ctx); got != "good" || !ok {
			t.Errorf("RequestID after a string key %q was added = (%q, %v), want (\"good\", true)", k, got, ok)
		}
	}
}

func TestShadowAndInherit(t *testing.T) {
	ctx := WithRequestID(WithRequestID(context.Background(), "outer"), "inner")
	if got, ok := RequestID(ctx); got != "inner" || !ok {
		t.Errorf("RequestID(inner over outer) = (%q, %v), want (\"inner\", true)", got, ok)
	}
	parent := WithRequestID(context.Background(), "a")
	child, cancel := context.WithTimeout(parent, time.Hour)
	defer cancel()
	child = context.WithValue(child, "x", "y")
	if got, ok := RequestID(child); got != "a" || !ok {
		t.Errorf("RequestID(child of ctx with \"a\") = (%q, %v), want (\"a\", true)", got, ok)
	}
	if got, ok := RequestID(parent); got != "a" || !ok {
		t.Errorf("RequestID(parent) = (%q, %v), want (\"a\", true)", got, ok)
	}
}

func TestDetach(t *testing.T) {
	req, cancel := context.WithTimeout(WithRequestID(context.Background(), "req-7"), time.Hour)
	d := Detach(req)
	if d == nil {
		t.Fatalf("Detach(req) returned a nil context")
	}
	cancel()
	if req.Err() == nil {
		t.Fatalf("req.Err() after cancel = <nil>, want %v", context.Canceled)
	}
	if err := d.Err(); err != nil {
		t.Errorf("Detach(req).Err() after req canceled = %v, want <nil>", err)
	}
	if d.Done() != nil {
		t.Errorf("Detach(req).Done() = non-nil channel, want <nil> (never canceled)")
	}
	if dl, ok := d.Deadline(); ok {
		t.Errorf("Detach(req).Deadline() = (%v, true), want (zero, false)", dl)
	}
	if got, ok := RequestID(d); got != "req-7" || !ok {
		t.Errorf("RequestID(Detach(req)) = (%q, %v), want (\"req-7\", true)", got, ok)
	}
}
