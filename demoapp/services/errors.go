package services

import (
	"errors"

	porterrors "github.com/brunojet/go-infra-backend/pkg/ports/errors"
)

var (
	errProfileScopeIDRequired             = porterrors.NewBusinessRuleError(errTextProfileScopeIDRequired)
	errNestedProfileApplicationIDRequired = porterrors.NewBusinessRuleError(errTextNestedProfileApplicationIDScopeRequired)
	errNestedVersionApplicationIDRequired = porterrors.NewBusinessRuleError(errTextNestedVersionApplicationIDScopeRequired)
	errNestedVersionIDRequired            = porterrors.NewBusinessRuleError(errTextNestedVersionIDScopeRequired)

	errMapperInvalidFileHash = errors.New("invalid file hash format")
	errMapperNilModel        = errors.New("model pointer is nil in mapper")
)
