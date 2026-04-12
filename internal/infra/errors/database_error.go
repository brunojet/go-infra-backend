package errors

import (
	stdErrors "errors"
	"fmt"
)

// DBErrorKind categorizes database errors for port-level handling without
// exposing specific driver or ORM implementations.
type DBErrorKind uint8

const (
	DBErrNotFound DBErrorKind = iota + 1
	DBErrUnavailable
	DBErrConstraintViolation
	DBErrInvalidTx
	DBErrConflictValidation
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
	case DBErrNotFound:
		return "record not found"
	case DBErrUnavailable:
		return "db unavailable"
	case DBErrConstraintViolation:
		if e.wrapped != nil {
			return fmt.Sprintf("constraint violation: %s", e.wrapped.Error())
		}
		return "constraint violation"
	case DBErrInvalidTx:
		return "invalid transaction"
	case DBErrConflictValidation:
		return "conflict validation failed: another transaction has modified the same entity"
	default:
		return fmt.Sprintf("db error: kind %d", e.kind)
	}
}

// Is enables errors.Is semantics: two databaseErrors match when they share the
// same DBErrorKind, regardless of the wrapped cause. This allows sentinel vars
// to be compared against new instances returned by MapDbError.
func (e *databaseError) Is(target error) bool {
	var t *databaseError
	return stdErrors.As(target, &t) && t.kind == e.kind
}

// Unwrap returns the original driver/ORM error, enabling errors.Is traversal to
// the underlying cause (e.g. the raw GORM constraint error).
func (e *databaseError) Unwrap() error { return e.wrapped }

// NewDatabaseError constructs a databaseError of the given kind. Optionally
// accepts the original cause to wrap — used by MapDbError for constraint violations.
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
	return stdErrors.As(err, &de)
}

// DatabaseErrorKind returns the DBErrorKind of the first databaseError in the
// chain, and true. Returns 0, false if err is not a databaseError.
func DatabaseErrorKind(err error) (DBErrorKind, bool) {
	var de *databaseError
	if stdErrors.As(err, &de) {
		return de.kind, true
	}
	return 0, false
}
