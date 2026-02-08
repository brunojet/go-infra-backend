package repositories

import (
	"context"

	internaldbcontracts "github.com/brunojet/go-infra-backend/internal/database/contracts"
	internalrepos "github.com/brunojet/go-infra-backend/internal/ports/repositories"
	"github.com/brunojet/go-infra-backend/pkg/ports/repositories/contracts"
	"gorm.io/gorm"
)

var (
	ErrDBUnavailable = internalrepos.ErrDBUnavailable
	ErrInvalidTx     = internalrepos.ErrInvalidTx
	ErrNotFound      = internalrepos.ErrNotFound
)

func MapDbError(err error) error { return internalrepos.MapDbError(err) }

func MapTxError(tx *gorm.DB) error { return internalrepos.MapTxError(tx) }

func ContextWithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return internalrepos.ContextWithTx(ctx, tx)
}

func TxFromContext(ctx context.Context) *gorm.DB { return internalrepos.TxFromContext(ctx) }

// NewGormRepository delegates to the internal implementation.
// It returns the Repository interface to avoid leaking unexported concrete types.
func NewGormRepository[E contracts.Entity](db internaldbcontracts.DatabaseAdapter) contracts.Repository[E] {
	return internalrepos.NewGormRepository[E](db)
}
