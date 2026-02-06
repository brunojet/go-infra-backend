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

	repo := NewGormRepository[TestEntity](db.GormDB())
	ctx := context.Background()

	err := repo.WithTx(ctx, func(ctx context.Context) error {
		if TxFromContext(ctx) == nil {
			return errors.New("tx not found in context")
		}
		return nil
	})
	assert.NoError(t, err)
}
