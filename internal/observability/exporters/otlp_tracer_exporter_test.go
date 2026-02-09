package exporters

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOTLPTracerExporterWithEndpoint(t *testing.T) {
	t.Setenv(OTLPEndpointEnv, otlpEndpointDefault)
	ctx := context.Background()
	exp, err := NewOTLPTracerExporter(ctx)
	require.NoError(t, err)
	assert.NotNil(t, exp)
}

func TestNewOTLPTracerExporterWithoutEndpoint(t *testing.T) {
	t.Setenv(OTLPEndpointEnv, "")
	ctx := context.Background()
	exp, err := NewOTLPTracerExporter(ctx)
	require.NoError(t, err)
	assert.NotNil(t, exp)
}
