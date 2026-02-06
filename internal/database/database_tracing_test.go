package database

import (
	"context"
	"testing"

	"github.com/brunojet/go-infra-backend/internal/observability/exporters"
	plugins "github.com/brunojet/go-infra-backend/internal/observability/gorm_plugins"
	"github.com/brunojet/go-infra-backend/internal/observability/providers"
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

	loggerExporter, err := exporters.NewOTLPLoggerExporter(ctx)
	assert.NoError(t, err)
	_, shutdown, err := providers.NewOTLPLoggerProvider(ctx, loggerExporter)
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

	p, err := plugins.NewOtelGormPlugin()
	assert.NoError(t, err)

	db, err := NewSQLiteDatabase("memory", p)
	assert.NoError(t, err)
	defer func() {
		_ = db.Close()
	}()

	err = db.GormDB().AutoMigrate(&testEntity{})
	assert.NoError(t, err)

	e := testEntity{Name: "trace-test"}
	err = db.GormDB().Create(&e).Error
	assert.NoError(t, err)

	// Read
	var got testEntity
	err = db.GormDB().First(&got, "name = ?", "trace-test").Error
	assert.NoError(t, err)
}
