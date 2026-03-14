package services

import (
	"errors"

	porterrors "github.com/brunojet/go-infra-backend/pkg/ports/errors"
)

var (
	errProfileScopeIDRequired             = porterrors.NewBusinessRuleError(errors.New(errTextProfileScopeIDRequired))
	errNestedProfileApplicationIDRequired = porterrors.NewBusinessRuleError(errors.New(errTextNestedProfileApplicationIDScopeRequired))
	errNestedVersionApplicationIDRequired = porterrors.NewBusinessRuleError(errors.New(errTextNestedVersionApplicationIDScopeRequired))
	errNestedVersionIDRequired            = porterrors.NewBusinessRuleError(errors.New(errTextNestedVersionIDScopeRequired))
)
