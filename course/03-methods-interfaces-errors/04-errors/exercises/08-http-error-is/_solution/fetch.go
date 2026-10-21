package fetch

import (
	"errors"
	"fmt"
)

// Error categories an *HTTPError can match with errors.Is.
var (
	ErrClient        = errors.New("client error")        // any 4xx
	ErrServer        = errors.New("server error")        // any 5xx
	ErrRetryable     = errors.New("retryable")           // 429, 502, 503, 504
	ErrInvalidStatus = errors.New("invalid status code") // outside 100..599
)

// HTTPError is a failed HTTP response.
type HTTPError struct {
	URL  string
	Code int
}

// Error returns "<url>: status <code>".
func (e *HTTPError) Error() string {
	return fmt.Sprintf("%s: status %d", e.URL, e.Code)
}

// Is reports whether e belongs to the category target, or is an *HTTPError
// with the same code.
func (e *HTTPError) Is(target error) bool {
	switch target {
	case ErrClient:
		return e.Code >= 400 && e.Code <= 499
	case ErrServer:
		return e.Code >= 500 && e.Code <= 599
	case ErrRetryable:
		switch e.Code {
		case 429, 502, 503, 504:
			return true
		}
		return false
	}
	t, ok := target.(*HTTPError)
	return ok && t != nil && t.Code == e.Code
}

// Check turns an HTTP status code into an error (nil for 100..399).
func Check(url string, code int) error {
	switch {
	case code >= 100 && code <= 399:
		return nil
	case code >= 400 && code <= 599:
		return &HTTPError{URL: url, Code: code}
	default:
		return fmt.Errorf("%s: status %d: %w", url, code, ErrInvalidStatus)
	}
}

// Advice tells the user what to do about err.
func Advice(err error) string {
	switch {
	case err == nil:
		return "ok"
	case errors.Is(err, ErrRetryable):
		return "retry later"
	case errors.Is(err, ErrClient):
		return "fix the request"
	case errors.Is(err, ErrServer):
		return "server problem"
	default:
		return "unknown error"
	}
}
