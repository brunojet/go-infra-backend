package errors

import (
	"errors"
	stdErrors "errors"
)

var ErrBusinessRuleViolation = stdErrors.New("business rule violation")

type BusinessRuleError struct {
	cause error
}

func (e *BusinessRuleError) Error() string {
	if e == nil || e.cause == nil {
		return ErrBusinessRuleViolation.Error()
	}

	return e.cause.Error()
}

func (e *BusinessRuleError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.cause
}

func NewBusinessRuleError(message string) error {
	cause := errors.New(message)
	return &BusinessRuleError{cause: cause}
}

func IsBusinessRuleError(err error) bool {
	var businessErr *BusinessRuleError
	return stdErrors.As(err, &businessErr)
}
