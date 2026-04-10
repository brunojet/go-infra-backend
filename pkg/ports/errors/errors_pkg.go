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
