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

type ListParams = contracts.ListParams

type Entity = contracts.Entity

type LockValidationSpec[E Entity] = contracts.LockValidationSpec[E]
type ConflictAction = contracts.ConflictAction

type Repository[E Entity] = contracts.Repository[E]

type BusinessRuleError = errs.BusinessRuleError

// ---- Errors ----

var (
	ErrDBUnavailable         = internalrepos.ErrDBUnavailable
	ErrInvalidTx             = internalrepos.ErrInvalidTx
	ErrNotFound              = internalrepos.ErrNotFound
	ErrRequiresTransaction   = internalrepos.ErrRequiresTransaction
	ErrBusinessRuleViolation = internalrepos.ErrBusinessRuleViolation
	ErrLockValidationWhere   = internalrepos.ErrLockValidationWhere
)

const (
	ConflictActionError  = contracts.ConflictActionError
	ConflictActionIgnore = contracts.ConflictActionIgnore
	ConflictActionUpdate = contracts.ConflictActionUpdate
)

// ---- Helpers (delegating to internal) ----

func MapDbError(err error) error { return internalrepos.MapDbError(err) }

func MapTxError(tx *gorm.DB) error { return internalrepos.MapTxError(tx) }

func ContextWithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return internalrepos.ContextWithTx(ctx, tx)
}

func TxFromContext(ctx context.Context) (*gorm.DB, error) { return internalrepos.TxFromContext(ctx) }

func GetContextFromTx(tx *gorm.DB) context.Context { return internalrepos.GetContextFromTx(tx) }

func ValidateTxWithUpdateLock[E Entity](tx *gorm.DB, spec contracts.LockValidationSpec[E]) error {
	return internalrepos.ValidateTxWithUpdateLock(tx, spec)
}

func AddOnConflict(tx *gorm.DB, action ConflictAction, columnNames ...string) error {
	return internalrepos.AddOnConflict(tx, action, columnNames...)
}

func AddOnConflictDoNothing(tx *gorm.DB, columnNames ...string) error {
	return internalrepos.AddOnConflictDoNothing(tx, columnNames...)
}

func AddOnConflictUpdateAll(tx *gorm.DB, columnNames ...string) error {
	return internalrepos.AddOnConflictUpdateAll(tx, columnNames...)
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
