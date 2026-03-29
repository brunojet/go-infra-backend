package repositories

import (
	"context"

	"github.com/brunojet/go-infra-backend/internal/ports/repositories"
	"github.com/brunojet/go-infra-backend/pkg/database"
	"github.com/brunojet/go-infra-backend/pkg/ports/errors"
	"github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
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

// ---- Errors ----
var (
	ErrDBUnavailable              = repositories.ErrDBUnavailable
	ErrInvalidTx                  = repositories.ErrInvalidTx
	ErrNotFound                   = repositories.ErrNotFound
	ErrRequiresTransaction        = repositories.ErrRequiresTransaction
	ErrBusinessRuleViolation      = repositories.ErrBusinessRuleViolation
	ErrLockValidationWhere        = repositories.ErrLockValidationWhere
	ErrConflictValidationRequired = repositories.ErrConflictValidationRequired
)

// ---- Helpers (delegating to internal) ----
func MapDbError(err error) error { return repositories.MapDbError(err) }

func MapTxError(tx *gorm.DB) error { return repositories.MapTxError(tx) }

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

func NewBusinessRuleError(cause error) error {
	return errors.NewBusinessRuleError(cause)
}

func IsBusinessRuleError(err error) bool {
	return errors.IsBusinessRuleError(err)
}

// NewGormRepository delegates to the internal implementation.
// It returns the Repository interface to avoid leaking unexported concrete types.
func NewGormRepository[E Entity](db database.DatabaseAdapter) Repository[E] {
	return repositories.NewGormRepository[E](db)
}
