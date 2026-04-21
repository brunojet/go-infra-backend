package errors

import (
	"errors"
	"net/http"
)

type RestError struct {
	statusCode int
	cause      []error
}

func (e RestError) Error() string {
	if e.cause == nil {
		return ErrBusinessRuleViolation.Error()
	}
	return e.cause[0].Error()
}

func (e *RestError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause[0]
}

func NewRestError(statusCode int, cause ...error) error {
	if len(cause) == 0 {
		cause = []error{errors.New(http.StatusText(statusCode))}
	}
	return &RestError{cause: cause, statusCode: statusCode}
}

func IsRestError(err error) bool {
	var businessErr *RestError
	return errors.As(err, &businessErr)
}

func GetRestErrorStatusCode(err error) (int, bool) {
	var restErr *RestError
	if errors.As(err, &restErr) {
		return restErr.statusCode, true
	}
	return 0, false
}
