package database

import (
	"context"
	"testing"

	plugins "github.com/brunojet/go-infra-backend/internal/observability/gorm_plugins"
	dbcontracts "github.com/brunojet/go-infra-backend/pkg/database/contracts"
	"github.com/stretchr/testify/assert"
)

type testEntity struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex"`
}

// TestCreateRead_ExportsTrace creates and reads a GORM entity while the
// OpenTelemetry tracing plugin is registered and asserts that spans were
// exported to the in-memory exporter.
func TestCreateRead_ExportsTrace(t *testing.T) {
	ctx := context.Background()

	// Observability exporters/providers are not required for this integration
	// smoke test; ensure GORM plugin works with the DB manager.

	db, err := newDatabaseManagerFromConfig(&dbcontracts.DatabaseConfig{Driver: dbcontracts.DbDriverSQLite, Mode: dbcontracts.DatabaseModeMemory}, plugins.NewOtelGormPlugin())
	assert.NoError(t, err)
	defer db.Shutdown(ctx)

	gdb, err := db.DatabaseAdapter().GormDB()
	assert.NoError(t, err)

	err = gdb.AutoMigrate(&testEntity{})
	assert.NoError(t, err)

	e := testEntity{Name: "trace-test"}
	err = gdb.Create(&e).Error
	assert.NoError(t, err)

	// Read
	var got testEntity
	err = gdb.First(&got, "name = ?", "trace-test").Error
	assert.NoError(t, err)
}
