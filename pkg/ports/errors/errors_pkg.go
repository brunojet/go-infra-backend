package errors

import internalerrors "github.com/brunojet/go-infra-backend/internal/ports/errors"

type BusinessRuleError = internalerrors.BusinessRuleError

var ErrBusinessRuleViolation = internalerrors.ErrBusinessRuleViolation

func NewBusinessRuleError(message string) error {
	return internalerrors.NewBusinessRuleError(message)
}

func IsBusinessRuleError(err error) bool {
	return internalerrors.IsBusinessRuleError(err)
}

// NewHTTPError returns a transport-agnostic error carrying the given HTTP status
// code. The handler boundary reads it via HTTPStatusCode without importing any
// transport or infra package.
func NewHTTPError(statusCode int) error {
	return internalerrors.NewHTTPError(statusCode)
}

func IsHTTPError(err error) bool {
	return internalerrors.IsHTTPError(err)
}

func HTTPStatusCode(err error) (int, bool) {
	return internalerrors.HTTPStatusCode(err)
}
