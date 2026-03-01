package errors

import internalerrors "github.com/brunojet/go-infra-backend/internal/ports/errors"

type BusinessRuleError = internalerrors.BusinessRuleError

var ErrBusinessRuleViolation = internalerrors.ErrBusinessRuleViolation

func NewBusinessRuleError(cause error) error {
	return internalerrors.NewBusinessRuleError(cause)
}

func IsBusinessRuleError(err error) bool {
	return internalerrors.IsBusinessRuleError(err)
}
