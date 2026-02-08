package repositories

import (
	"context"

	internalrepos "github.com/brunojet/go-infra-backend/internal/ports/repositories"
	db "github.com/brunojet/go-infra-backend/pkg/database"
	"github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"gorm.io/gorm"
)

// ---- Contracts ----

type ListParams = contracts.ListParams

type Entity = contracts.Entity

type Repository[E Entity] = contracts.Repository[E]

// ---- Errors ----

var (
	ErrDBUnavailable = internalrepos.ErrDBUnavailable
	ErrInvalidTx     = internalrepos.ErrInvalidTx
	ErrNotFound      = internalrepos.ErrNotFound
)

// ---- Helpers (delegating to internal) ----

func MapDbError(err error) error { return internalrepos.MapDbError(err) }

func MapTxError(tx *gorm.DB) error { return internalrepos.MapTxError(tx) }

func ContextWithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return internalrepos.ContextWithTx(ctx, tx)
}

func TxFromContext(ctx context.Context) *gorm.DB { return internalrepos.TxFromContext(ctx) }

// NewGormRepository delegates to the internal implementation.
// It returns the Repository interface to avoid leaking unexported concrete types.
func NewGormRepository[E Entity](db db.DatabaseAdapter) Repository[E] {
	return internalrepos.NewGormRepository[E](db)
}
