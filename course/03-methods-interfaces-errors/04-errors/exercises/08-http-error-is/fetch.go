package fetch

import "errors"

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
	// TODO
	return ""
}

// Is reports whether e belongs to the category target, or is an *HTTPError
// with the same code.
func (e *HTTPError) Is(target error) bool {
	// TODO
	return false
}

// Check turns an HTTP status code into an error (nil for 100..399).
func Check(url string, code int) error {
	// TODO
	return nil
}

// Advice tells the user what to do about err.
func Advice(err error) string {
	// TODO
	return ""
}
