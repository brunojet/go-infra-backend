package repositories

import (
	"context"

	internalrepos "github.com/brunojet/go-infra-backend/internal/ports/repositories"
	db "github.com/brunojet/go-infra-backend/pkg/database"
	errs "github.com/brunojet/go-infra-backend/pkg/ports/errors"
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
	BusinessRuleError            = errs.BusinessRuleError
)

// ---- Errors ----
var (
	ErrDBUnavailable              = internalrepos.ErrDBUnavailable
	ErrInvalidTx                  = internalrepos.ErrInvalidTx
	ErrNotFound                   = internalrepos.ErrNotFound
	ErrRequiresTransaction        = internalrepos.ErrRequiresTransaction
	ErrBusinessRuleViolation      = internalrepos.ErrBusinessRuleViolation
	ErrLockValidationWhere        = internalrepos.ErrLockValidationWhere
	ErrConflictValidationRequired = internalrepos.ErrConflictValidationRequired
)

// ---- Helpers (delegating to internal) ----
func MapDbError(err error) error { return internalrepos.MapDbError(err) }

func MapTxError(tx *gorm.DB) error { return internalrepos.MapTxError(tx) }

// ContextWithTx is a convenience wrapper that annotates a context with a
// *gorm.DB transaction so downstream code can retrieve it via TxFromContext.
func ContextWithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return internalrepos.ContextWithTx(ctx, tx)
}

func TxFromContext(ctx context.Context) (*gorm.DB, error) { return internalrepos.TxFromContext(ctx) }

func ValidateTxWithUpdateLock[E Entity](tx *gorm.DB, spec contracts.LockValidationSpec[E]) error {
	return internalrepos.ValidateTxWithUpdateLock(tx, spec)
}

func AddOnConflictDoNothing(tx *gorm.DB, columnNames ...string) error {
	return internalrepos.AddOnConflictDoNothing(tx, columnNames...)
}

func AddOnConflictUpdateAll(tx *gorm.DB, columnNames ...string) error {
	return internalrepos.AddOnConflictUpdateAll(tx, columnNames...)
}

func WhereOnConflict(tx *gorm.DB, scopes ...contracts.ConflictScope) *gorm.DB {
	return internalrepos.WhereOnConflict(tx, scopes...)
}

func NewBusinessRuleError(cause error) error {
	return errs.NewBusinessRuleError(cause)
}

func IsBusinessRuleError(err error) bool {
	return errs.IsBusinessRuleError(err)
}

// NewGormRepository delegates to the internal implementation.
// It returns the Repository interface to avoid leaking unexported concrete types.
func NewGormRepository[E Entity](db db.DatabaseAdapter) Repository[E] {
	return internalrepos.NewGormRepository[E](db)
}
