package repositories

import (
	"context"

	"gorm.io/gorm"
)

type ctxKeyTx struct{}

// ContextWithTx returns a new context that carries the given *gorm.DB transaction.
func ContextWithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, ctxKeyTx{}, tx)
}

// TxFromContext extracts a *gorm.DB transaction from the context, or nil if not present.
func TxFromContext(ctx context.Context) *gorm.DB {
	if v := ctx.Value(ctxKeyTx{}); v != nil {
		if tx, ok := v.(*gorm.DB); ok {
			return tx
		}
	}
	return nil
}
