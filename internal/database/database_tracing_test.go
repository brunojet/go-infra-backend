package database

import (
	"context"
	"testing"

	"github.com/brunojet/go-infra-backend/internal/observability/exporters"
	plugins "github.com/brunojet/go-infra-backend/internal/observability/gorm_plugins"
	"github.com/brunojet/go-infra-backend/internal/observability/providers"
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

	loggerExporter, err := exporters.NewOTLPLoggerExporters(ctx)
	assert.NoError(t, err)
	_, shutdown, err := providers.NewOTLPLoggerProvider(ctx, loggerExporter...)
	assert.NoError(t, err)
	defer shutdown(ctx)
	metricExporter, err := exporters.NewOTLPMetricExporter(ctx)
	assert.NoError(t, err)
	_, shutdown, err = providers.NewOTLPMetricProvider(ctx, metricExporter)
	assert.NoError(t, err)
	defer shutdown(ctx)
	tracerExporter, err := exporters.NewOTLPTracerExporter(ctx)
	assert.NoError(t, err)
	_, shutdown, err = providers.NewOTLPTracerProvider(ctx, tracerExporter)
	assert.NoError(t, err)
	defer shutdown(ctx)

	db, err := newDatabaseManagerFromConfig(&dbcontracts.DatabaseConfig{Driver: dbcontracts.DbDriverSQLiteMemory}, plugins.NewOtelGormPlugin())
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
