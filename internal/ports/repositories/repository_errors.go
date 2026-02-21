package repositories

import (
	"database/sql"
	"errors"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrDBUnavailable = sql.ErrConnDone
	ErrInvalidTx     = errors.New("invalid transaction")
	ErrNotFound      = gorm.ErrRecordNotFound
	// Validation / user-level repository errors
	ErrInvalidConflictColumns    = errors.New("invalid conflict columns")
	ErrInvalidConflictColumnName = errors.New("invalid conflict column: field name cannot be empty")
	ErrInvalidScope              = errors.New("invalid scope: field has invalid value")
	ErrEmptyScopes               = errors.New("invalid scope: scopes must be non-empty and contain valid field names and values")
	ErrOrderByMissing            = errors.New("orderBy must be provided")
	ErrInvalidPage               = errors.New("page must be greater than zero")
	ErrInvalidPageSize           = errors.New("pageSize must be greater than zero")
)

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
