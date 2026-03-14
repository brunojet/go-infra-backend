package services

import (
	"errors"

	porterrors "github.com/brunojet/go-infra-backend/pkg/ports/errors"
)

// Centralized error variables for service-level business rule errors.
var (
	errProfileScopeIDRequired = porterrors.NewBusinessRuleError(errors.New(errTextProfileScopeIDRequired))
	errVersionScopeIDRequired = porterrors.NewBusinessRuleError(errors.New(errTextVersionScopeIDRequired))

	errNestedProfileApplicationIDRequired = porterrors.NewBusinessRuleError(errors.New(errTextNestedProfileApplicationIDScopeRequired))

	errNestedVersionApplicationIDRequired = porterrors.NewBusinessRuleError(errors.New(errTextNestedVersionApplicationIDScopeRequired))
	errNestedVersionTerminalIDRequired    = porterrors.NewBusinessRuleError(errors.New(errTextNestedVersionTerminalModelIDScopeRequired))
	errNestedVersionIDRequired            = porterrors.NewBusinessRuleError(errors.New(errTextNestedVersionIDScopeRequired))
)
