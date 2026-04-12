package errors

import (
	stdErrors "errors"
	"fmt"
)

// httpError is an unexported error type that carries an HTTP status code to be
// forwarded to the final client. Construct via NewHTTPError; inspect via
// IsHTTPError and HTTPStatusCode — without ever accessing the struct directly.
type httpError struct {
	statusCode int
}

func (e *httpError) Error() string {
	return fmt.Sprintf("http error: status %d", e.statusCode)
}

// NewHTTPError returns an error wrapping the given HTTP status code.
func NewHTTPError(statusCode int) error {
	return &httpError{statusCode: statusCode}
}

func IsHTTPError(err error) bool {
	var he *httpError
	return stdErrors.As(err, &he)
}

func HTTPStatusCode(err error) (int, bool) {
	var he *httpError
	if stdErrors.As(err, &he) {
		return he.statusCode, true
	}
	return 0, false
}
