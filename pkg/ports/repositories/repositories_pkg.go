package repositories

import (
	"context"

	"github.com/brunojet/go-infra-backend/internal/ports/backend/repositories"
	rpoerrs "github.com/brunojet/go-infra-backend/internal/ports/backend/repositories/errors"
	"github.com/brunojet/go-infra-backend/pkg/infra/database"
	"github.com/brunojet/go-infra-backend/pkg/ports/backend/repositories/contracts"
	"github.com/brunojet/go-infra-backend/pkg/ports/errors"
	"gorm.io/gorm"
)

// ---- Contracts ----
type (
	ListParams                   = contracts.ListParams
	Entity                       = contracts.Entity
	LockValidationSpec[E Entity] = contracts.LockValidationSpec[E]
	Repository[E Entity]         = contracts.Repository[E]
	ConflictScope                = contracts.ConflictScope
	BusinessRuleError            = errors.BusinessRuleError
)

// ---- Error kinds ----
type DBErrorKind = rpoerrs.DBErrorKind

const (
	DBErrUnavailable       = rpoerrs.DBErrUnavailable
	DBErrInvalidParameters = rpoerrs.DBErrInvalidParameters
	DBErrConstraint        = rpoerrs.DBErrConstraint
	DBErrNotFound          = rpoerrs.DBErrNotFound
)

// ---- Helpers (delegating to internal) ----
func MapDbError(err error) error { return rpoerrs.MapDbError(err) }

func MapTxError(tx *gorm.DB) error { return rpoerrs.MapTxError(tx) }

func IsDatabaseErrorKind(err error, kind DBErrorKind) bool {
	return rpoerrs.IsDatabaseErrorKind(err, kind)
}

func DatabaseErrorKind(err error) (DBErrorKind, bool) { return rpoerrs.DatabaseErrorKind(err) }

// ContextWithTx is a convenience wrapper that annotates a context with a
// *gorm.DB transaction so downstream code can retrieve it via TxFromContext.
func ContextWithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return repositories.ContextWithTx(ctx, tx)
}

func TxFromContext(ctx context.Context) (*gorm.DB, error) { return repositories.TxFromContext(ctx) }

func ValidateTxWithUpdateLock[E Entity](tx *gorm.DB, spec LockValidationSpec[E]) error {
	return repositories.ValidateTxWithUpdateLock(tx, spec)
}

func AddOnConflictDoNothing(tx *gorm.DB, columnNames ...string) error {
	return repositories.AddOnConflictDoNothing(tx, columnNames...)
}

func AddOnConflictUpdateAll(tx *gorm.DB, columnNames ...string) error {
	return repositories.AddOnConflictUpdateAll(tx, columnNames...)
}

func WhereOnConflict(tx *gorm.DB, scopes ...ConflictScope) *gorm.DB {
	return repositories.WhereOnConflict(tx, scopes...)
}

// NewGormRepository delegates to the internal implementation.
// It returns the Repository interface to avoid leaking unexported concrete types.
func NewGormRepository[E Entity](db database.DatabaseAdapter) Repository[E] {
	return repositories.NewGormRepository[E](db)
}
