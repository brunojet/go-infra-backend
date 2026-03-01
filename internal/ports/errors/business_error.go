package errors

import stdErrors "errors"

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

func NewBusinessRuleError(cause error) error {
	if cause == nil {
		return nil
	}

	var businessErr *BusinessRuleError
	if stdErrors.As(cause, &businessErr) {
		return cause
	}

	return &BusinessRuleError{cause: cause}
}

func IsBusinessRuleError(err error) bool {
	var businessErr *BusinessRuleError
	return stdErrors.As(err, &businessErr)
}
