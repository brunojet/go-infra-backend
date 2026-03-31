package repositories

import (
	"context"
	"testing"

	"github.com/brunojet/go-infra-backend/pkg/database"
	porterrors "github.com/brunojet/go-infra-backend/pkg/ports/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type pkgTestEntity struct {
	ID   string `gorm:"primaryKey"`
	Name string
}

func (pkgTestEntity) TableName() string { return "pkg_test_entities" }
func (e pkgTestEntity) WhereOnConflict(tx *gorm.DB) *gorm.DB {
	return tx
}

func TestValidateTxWithUpdateLock_ExposedForPkgConsumers(t *testing.T) {
	adapter, err := database.NewSQLite("memory")
	require.NoError(t, err)
	defer func() { _ = adapter.Close() }()

	db, err := adapter.GormDB()
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&pkgTestEntity{}))
	require.NoError(t, db.Create(&pkgTestEntity{ID: "1", Name: "alpha"}).Error)
	repo := NewGormRepository[pkgTestEntity](adapter)

	err = repo.WithTx(context.Background(), func(txCtx context.Context) error {
		tx, txErr := TxFromContext(txCtx)
		require.NoError(t, txErr)
		require.NotNil(t, tx)
		tx = tx.WithContext(txCtx)
		return ValidateTxWithUpdateLock(tx, LockValidationSpec[pkgTestEntity]{
			SelectColumns: []string{"id", "name"},
			WhereSQL:      "id = ?",
			WhereArgs:     []any{"1"},
			BlockIfFound: func(found *pkgTestEntity) error {
				if found.Name == "alpha" {
					return nil
				}
				return porterrors.ErrBusinessRuleViolation
			},
		})
	})
	require.NoError(t, err)

	err = repo.WithTx(context.Background(), func(txCtx context.Context) error {
		tx, txErr := TxFromContext(txCtx)
		require.NoError(t, txErr)
		require.NotNil(t, tx)
		tx = tx.WithContext(txCtx)
		return ValidateTxWithUpdateLock(tx, LockValidationSpec[pkgTestEntity]{
			WhereSQL:  "id = ?",
			WhereArgs: []any{"1"},
		})
	})
	assert.ErrorIs(t, err, porterrors.ErrBusinessRuleViolation)
}

func TestPkgFacade_AddOnConflictWrappers(t *testing.T) {
	adapter, err := database.NewSQLite("memory")
	require.NoError(t, err)
	defer func() { _ = adapter.Close() }()

	db, err := adapter.GormDB()
	require.NoError(t, err)

	err = AddOnConflictDoNothing(nil, "id")
	assert.ErrorIs(t, err, ErrInvalidTx)

	err = AddOnConflictUpdateAll(nil, "id")
	assert.ErrorIs(t, err, ErrInvalidTx)

	err = db.Transaction(func(tx *gorm.DB) error {
		if txErr := AddOnConflictDoNothing(tx, "id"); txErr != nil {
			return txErr
		}
		return AddOnConflictUpdateAll(tx, "id")
	})
	assert.NoError(t, err)
}
