package database

import (
	"database/sql"
	"errors"
	"strings"

	infraerrors "github.com/brunojet/go-infra-backend/internal/infra/errors"
	"gorm.io/gorm"
)

// Note: MapTxError lives in internal/ports/backend/repositories because it
// interprets *gorm.DB result fields (Error, RowsAffected) that are specific
// to the GORM repository pattern, not a general database infrastructure concern.

// DB-layer sentinels — stable port identifiers, independent of any specific driver.
// Port packages alias these vars so application code never imports GORM directly.
var (
	ErrNotFound            = infraerrors.NewDatabaseError(infraerrors.DBErrNotFound)
	ErrDBUnavailable       = infraerrors.NewDatabaseError(infraerrors.DBErrUnavailable)
	ErrConstraintViolation = infraerrors.NewDatabaseError(infraerrors.DBErrConstraintViolation)
	ErrInvalidTx           = infraerrors.NewDatabaseError(infraerrors.DBErrInvalidTx)
)

// MapDbError translates raw database / driver errors into stable port sentinels.
// Unknown errors are returned as-is so callers can still inspect them.
func MapDbError(err error) error {
	if err == nil {
		return nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	} else if errors.Is(err, sql.ErrConnDone) || strings.Contains(err.Error(), "database is closed") {
		return ErrDBUnavailable
	} else if errors.Is(err, gorm.ErrDuplicatedKey) ||
		strings.Contains(err.Error(), "UNIQUE constraint failed") ||
		strings.Contains(err.Error(), "CHECK constraint failed") ||
		strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
		return infraerrors.NewDatabaseError(infraerrors.DBErrConstraintViolation, err)
	}
	return err
}
