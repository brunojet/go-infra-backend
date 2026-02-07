package providers

import (
	"context"
	"testing"

	"github.com/brunojet/go-infra-backend/internal/observability/exporters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOTLPTracerProvider(t *testing.T) {
	exp := exporters.NewOTLPNoopTracerExporter()

	tp, shutdown, err := NewOTLPTracerProvider(context.Background(), exp)
	require.NoError(t, err)
	assert.NotNil(t, tp)
	require.NotNil(t, shutdown)
	assert.NoError(t, shutdown(context.Background()))
}
