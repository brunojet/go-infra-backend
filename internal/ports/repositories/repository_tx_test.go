package repositories

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGormRepository_WithTx(t *testing.T) {
	db, cleanup := openMemoryDB(t)
	defer cleanup()

	repo := NewGormRepository[TestEntity](db)
	ctx := context.Background()

	err := repo.WithTx(ctx, func(ctx context.Context) error {
		if TxFromContext(ctx) == nil {
			return errors.New("tx not found in context")
		}
		return nil
	})
	assert.NoError(t, err)
}

func TestTxFromContext_AbsentReturnsNil(t *testing.T) {
	ctx := context.Background()
	tx := TxFromContext(ctx)
	assert.Nil(t, tx)
}

func TestTxFromContext_WithGormDB(t *testing.T) {
	db, cleanup := openMemoryDB(t)
	defer cleanup()
	gdb := mustGormDB(t, db)
	ctx := ContextWithTx(context.Background(), gdb)
	got := TxFromContext(ctx)
	assert.Same(t, gdb, got)
}

func TestTxFromContext_WrongTypeReturnsNil(t *testing.T) {
	// store a value under the same key but with wrong type
	ctx := context.WithValue(context.Background(), ctxKeyTx{}, "not-a-tx")
	got := TxFromContext(ctx)
	assert.Nil(t, got)
}
