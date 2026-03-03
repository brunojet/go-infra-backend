package repositories

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/brunojet/go-infra-backend/pkg/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type pkgTestEntity struct {
	ID   string `gorm:"primaryKey"`
	Name string
}

func (pkgTestEntity) TableName() string { return "pkg_test_entities" }

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
				return ErrBusinessRuleViolation
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
	assert.ErrorIs(t, err, ErrBusinessRuleViolation)
}

func TestPkgFacade_MapAndBusinessRuleHelpers(t *testing.T) {
	assert.ErrorIs(t, MapDbError(gorm.ErrRecordNotFound), ErrNotFound)
	assert.ErrorIs(t, MapDbError(sql.ErrNoRows), ErrNotFound)
	assert.Nil(t, MapDbError(nil))

	assert.ErrorIs(t, MapTxError(nil), ErrInvalidTx)
	assert.ErrorIs(t, MapTxError(&gorm.DB{RowsAffected: 0}), ErrNotFound)
	assert.NoError(t, MapTxError(&gorm.DB{RowsAffected: 1}))

	sentinel := errors.New("business-cause")
	err := NewBusinessRuleError(sentinel)
	assert.True(t, IsBusinessRuleError(err))
	assert.False(t, IsBusinessRuleError(errors.New("plain-error")))
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
