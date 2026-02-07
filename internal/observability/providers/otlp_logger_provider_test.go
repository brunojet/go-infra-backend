package providers

import (
	"context"
	"testing"

	"github.com/brunojet/go-infra-backend/internal/observability/exporters"
	"github.com/stretchr/testify/assert"
)

func TestNewOTLPLoggerProvider(t *testing.T) {
	exporter := exporters.NewOTLPNoopLoggerExporter()

	lp, shutdown, err := NewOTLPLoggerProvider(context.Background(), exporter)
	assert.NoError(t, err)
	assert.NotNil(t, lp)
	assert.NotNil(t, shutdown)
	assert.NoError(t, shutdown(context.Background()))
}

func TestNewOTLPLoggerProvider_NoExporters(t *testing.T) {
	lp, shutdown, err := NewOTLPLoggerProvider(context.Background())
	assert.Error(t, err)
	assert.Nil(t, lp)
	assert.Nil(t, shutdown)
}

func TestNewOTLPLoggerProvider_MultipleExporters(t *testing.T) {
	exp1 := exporters.NewOTLPNoopLoggerExporter()
	exp2 := exporters.NewOTLPNoopLoggerExporter()

	lp, shutdown, err := NewOTLPLoggerProvider(context.Background(), exp1, exp2)
	assert.NoError(t, err)
	assert.NotNil(t, lp)
	assert.NotNil(t, shutdown)
	assert.NoError(t, shutdown(context.Background()))
}
