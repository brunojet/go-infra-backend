package providers

import (
	"context"
	"testing"

	"github.com/brunojet/go-infra-backend/internal/observability/exporters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOTLPMetricProvider(t *testing.T) {
	t.Setenv(exporters.OTLPEndpointEnv, "")
	ctx := context.Background()
	exporter, err := exporters.NewOTLPMetricExporter(ctx)
	require.NoError(t, err)
	mp, shutdown, err := NewOTLPMetricProvider(ctx, exporter)
	require.NoError(t, err)
	assert.NotNil(t, mp)
	require.NotNil(t, shutdown)
	assert.NoError(t, shutdown(ctx))
}
