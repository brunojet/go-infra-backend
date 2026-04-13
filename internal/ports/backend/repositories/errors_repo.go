package repositories

// Re-exports from the errors sub-package so that callers of
// internal/ports/backend/repositories do not need to import the sub-package.

import (
	rpoerrs "github.com/brunojet/go-infra-backend/internal/ports/backend/repositories/errors"
	"gorm.io/gorm"
)

// DBErrorKind is re-exported as a type alias.
type DBErrorKind = rpoerrs.DBErrorKind

// Error kind constants re-exported for use in handler-level HTTP status mapping.
const (
	DBErrUnavailable       = rpoerrs.DBErrUnavailable
	DBErrInvalidParameters = rpoerrs.DBErrInvalidParameters
	DBErrConstraint        = rpoerrs.DBErrConstraint
	DBErrNotFound          = rpoerrs.DBErrNotFound
)

// Sentinel errors re-exported for callers that match via errors.Is.
var (
	ErrDBUnavailable              = rpoerrs.ErrDBUnavailable
	ErrInvalidTx                  = rpoerrs.ErrInvalidTx
	ErrNotFound                   = rpoerrs.ErrNotFound
	ErrRequiresTransaction        = rpoerrs.ErrRequiresTransaction
	ErrLockValidationWhere        = rpoerrs.ErrLockValidationWhere
	ErrConflictValidationRequired = rpoerrs.ErrConflictValidationRequired
	ErrConflictValidationFailed   = rpoerrs.ErrConflictValidationFailed
	ErrConstraintViolation        = rpoerrs.ErrConstraintViolation
)

// Helper functions re-exported so pkg-level facades delegate without importing errors/.
func MapDbError(err error) error                      { return rpoerrs.MapDbError(err) }
func MapTxError(tx *gorm.DB) error                    { return rpoerrs.MapTxError(tx) }
func DatabaseErrorKind(err error) (DBErrorKind, bool) { return rpoerrs.DatabaseErrorKind(err) }
func IsDatabaseErrorKind(err error, kind DBErrorKind) bool {
	return rpoerrs.IsDatabaseErrorKind(err, kind)
}
