package exporters

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOTLPMetricExporterWithEndpoint(t *testing.T) {
	t.Setenv(OTLPEndpointEnv, otlpEndpointDefault)
	ctx := context.Background()
	exp, err := NewOTLPMetricExporter(ctx)
	require.NoError(t, err)
	assert.NotNil(t, exp)
}

func TestNewOTLPMetricExporterWithoutEndpoint(t *testing.T) {
	t.Setenv(OTLPEndpointEnv, "")
	ctx := context.Background()
	exp, err := NewOTLPMetricExporter(ctx)
	require.NoError(t, err)
	assert.NotNil(t, exp)
}
