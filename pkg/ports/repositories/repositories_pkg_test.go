package repositories

import (
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

	err = db.Transaction(func(tx *gorm.DB) error {
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

	err = db.Transaction(func(tx *gorm.DB) error {
		return ValidateTxWithUpdateLock(tx, LockValidationSpec[pkgTestEntity]{
			WhereSQL:  "id = ?",
			WhereArgs: []any{"1"},
		})
	})
	assert.ErrorIs(t, err, ErrBusinessRuleViolation)
}
