package repositories

import (
	"database/sql"
	"errors"
	"strings"

	porterrors "github.com/brunojet/go-infra-backend/internal/ports/errors"
	"gorm.io/gorm"
)

var (
	ErrDBUnavailable = sql.ErrConnDone
	ErrInvalidTx     = errors.New("invalid transaction")
	ErrNotFound      = gorm.ErrRecordNotFound
	// Validation / user-level repository errors
	ErrInvalidConflictColumns     = errors.New("invalid conflict columns")
	ErrInvalidConflictColumnName  = errors.New("invalid conflict column: field name cannot be empty")
	ErrInvalidScope               = errors.New("invalid scope: field has invalid value")
	ErrEmptyScopes                = errors.New("invalid scope: scopes must be non-empty and contain valid field names and values")
	ErrOrderByMissing             = errors.New("orderBy must be provided")
	ErrInvalidPage                = errors.New("page must be greater than zero")
	ErrInvalidPageSize            = errors.New("pageSize must be greater than zero")
	ErrRequiresTransaction        = errors.New("operation must run inside a transaction")
	ErrBusinessRuleViolation      = porterrors.ErrBusinessRuleViolation
	ErrLockValidationWhere        = errors.New("where clause must be provided for lock validation")
	ErrConflictValidationRequired = errors.New("conflict validation is required for this operation")
	ErrConflictValidationFailed   = errors.New("conflict validation failed: another transaction has modified the same entity")
)

type BusinessRuleError = porterrors.BusinessRuleError

func NewBusinessRuleError(cause error) error {
	return porterrors.NewBusinessRuleError(cause)
}

func IsBusinessRuleError(err error) bool {
	return porterrors.IsBusinessRuleError(err)
}

func MapDbError(err error) error {
	if err == nil {
		return nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	} else if errors.Is(err, sql.ErrConnDone) || strings.Contains(err.Error(), "database is closed") {
		return ErrDBUnavailable
	}
	return err
}

func MapTxError(tx *gorm.DB) error {
	if tx == nil {
		return ErrInvalidTx
	} else if tx.Error != nil {
		return MapDbError(tx.Error)
	} else if tx.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
