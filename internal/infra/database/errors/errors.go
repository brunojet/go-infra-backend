package errors

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// DBErrorKind categorizes database errors for port-level handling without
// exposing specific driver or ORM implementations.
type DBErrorKind uint8

const (
	DBErrUnavailable       DBErrorKind = iota + 1 // Results an internal server error (HTTP 500) and should be exposed to clients.
	DBErrInvalidParameters                        // Results an internal server error (HTTP 500) and should be exposed to clients.
	DbErrRuntime                                  // Results an internal server error (HTTP 500) and should be exposed to clients.
	DBErrConstraint                               // Results a conflict error (HTTP 409) and should be exposed to clients with the message from the wrapped error.
	DBErrNotFound                                 // Results a not found error (HTTP 404) and should be exposed to clients.
)

// DB-layer sentinels — stable port identifiers, independent of any specific driver.
// Port packages alias these vars so application code never imports GORM directly.
var (
	ErrDBUnavailable              = NewDatabaseError(DBErrUnavailable)
	ErrConflictColumnsMissing     = NewDatabaseError(DBErrInvalidParameters, errors.New("conflict columns missing"))
	ErrInvalidConflictColumnName  = NewDatabaseError(DBErrInvalidParameters, errors.New("invalid conflict column name"))
	ErrInvalidScope               = NewDatabaseError(DBErrInvalidParameters, errors.New("scope field has invalid value"))
	ErrScopesMissing              = NewDatabaseError(DBErrInvalidParameters, errors.New("scopes is missing"))
	ErrOrderByMissing             = NewDatabaseError(DBErrInvalidParameters, errors.New("orderBy is missing"))
	ErrInvalidPage                = NewDatabaseError(DBErrInvalidParameters, errors.New("page must be greater than zero"))
	ErrInvalidPageSize            = NewDatabaseError(DBErrInvalidParameters, errors.New("pageSize must be greater than zero"))
	ErrRequiresTransaction        = NewDatabaseError(DbErrRuntime, errors.New("operation must run inside a transaction"))
	ErrLockValidationWhere        = NewDatabaseError(DbErrRuntime, errors.New("where clause must be provided for lock validation"))
	ErrInvalidTx                  = NewDatabaseError(DbErrRuntime, errors.New("invalid transaction"))
	ErrConflictValidationRequired = NewDatabaseError(DBErrConstraint, errors.New("validation is required for this operation"))
	ErrConflictValidationFailed   = NewDatabaseError(DBErrConstraint, errors.New("another transaction has modified the same entity"))
	ErrConstraintViolation        = NewDatabaseError(DBErrConstraint)
	ErrNotFound                   = NewDatabaseError(DBErrNotFound)
)

// databaseError is an unexported error type that carries a DBErrorKind and an
// optional wrapped cause (original driver/ORM error).
// Construct via NewDatabaseError; inspect via IsDatabaseError and DatabaseErrorKind.
type databaseError struct {
	kind    DBErrorKind
	wrapped error
}

func (e *databaseError) Error() string {
	switch e.kind {
	case DBErrUnavailable:
		return "db unavailable"
	case DBErrInvalidParameters:
		if e.wrapped != nil {
			return e.wrapped.Error()
		}
		return "invalid parameters"
	case DbErrRuntime:
		if e.wrapped != nil {
			return e.wrapped.Error()
		}
		return "db runtime error"
	case DBErrConstraint:
		if e.wrapped != nil {
			return e.wrapped.Error()
		}
		return "constraint violation"
	case DBErrNotFound:
		return "record not found"
	default:
		return fmt.Sprintf("db error: kind %d", e.kind)
	}
}

func (e *databaseError) Is(target error) bool {
	var t *databaseError
	return errors.As(target, &t) && t.kind == e.kind
}

func (e *databaseError) Unwrap() error { return e.wrapped }

func NewDatabaseError(kind DBErrorKind, wrapped ...error) error {
	var w error
	if len(wrapped) > 0 {
		w = wrapped[0]
	}
	return &databaseError{kind: kind, wrapped: w}
}

// IsDatabaseError reports whether err (or any error in its chain) is a databaseError.
func IsDatabaseError(err error) bool {
	var de *databaseError
	return errors.As(err, &de)
}

func DatabaseErrorKind(err error) (DBErrorKind, bool) {
	var de *databaseError
	if errors.As(err, &de) {
		return de.kind, true
	}
	return 0, false
}

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
		return NewDatabaseError(DBErrConstraint, err)
	}
	return err
}
