package repositories

import (
	"errors"

	infradb "github.com/brunojet/go-infra-backend/internal/infra/database"
	infraerrors "github.com/brunojet/go-infra-backend/internal/infra/errors"
	"gorm.io/gorm"
)

// ErrConstraintViolation is the sentinel for any database constraint violation
// (UNIQUE, CHECK, FK). Use errors.Is to detect it; the original database error
// is also preserved in the chain via fmt.Errorf with multiple %%w.
var ErrConstraintViolation = infradb.ErrConstraintViolation

var (
	ErrDBUnavailable = infradb.ErrDBUnavailable
	ErrInvalidTx     = infradb.ErrInvalidTx
	ErrNotFound      = infradb.ErrNotFound
	// Validation / user-level repository errors
	ErrInvalidConflictColumns     = errors.New("invalid conflict columns")
	ErrInvalidConflictColumnName  = errors.New("invalid conflict column: field name cannot be empty")
	ErrInvalidScope               = errors.New("invalid scope: field has invalid value")
	ErrEmptyScopes                = errors.New("invalid scope: scopes must be non-empty and contain valid field names and values")
	ErrOrderByMissing             = errors.New("orderBy must be provided")
	ErrInvalidPage                = errors.New("page must be greater than zero")
	ErrInvalidPageSize            = errors.New("pageSize must be greater than zero")
	ErrRequiresTransaction        = errors.New("operation must run inside a transaction")
	ErrLockValidationWhere        = errors.New("where clause must be provided for lock validation")
	ErrConflictValidationRequired = errors.New("conflict validation is required for this operation")
	ErrConflictValidationFailed   = infraerrors.NewDatabaseError(infraerrors.DBErrConflictValidation)
)

// MapDbError delegates to the infra layer which owns the GORM/driver mapping logic.
func MapDbError(err error) error { return infradb.MapDbError(err) }

// MapTxError inspects the GORM transaction result and maps it to a port sentinel.
// Lives here rather than in infra/database because it is specific to the GORM
// repository pattern: it interprets *gorm.DB fields (Error, RowsAffected) that
// are only meaningful in the context of repository operations.
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
