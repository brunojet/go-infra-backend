package errors

import (
	stdErrors "errors"
	"fmt"
)

// httpError carries an HTTP status code so BFF adapters can forward upstream
// status codes without importing any transport or handler package directly.
// Construct via NewHTTPError; inspect via IsHTTPError and HTTPStatusCode.
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

// IsHTTPError reports whether err (or any error in its chain) is an httpError.
func IsHTTPError(err error) bool {
	var he *httpError
	return stdErrors.As(err, &he)
}

// HTTPStatusCode returns the status code carried by the first httpError in the
// chain, and true. Returns 0, false if err is not an httpError.
func HTTPStatusCode(err error) (int, bool) {
	var he *httpError
	if stdErrors.As(err, &he) {
		return he.statusCode, true
	}
	return 0, false
}
