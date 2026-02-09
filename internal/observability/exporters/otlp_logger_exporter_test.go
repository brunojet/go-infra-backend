package exporters

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOTLPLoggerExporterWithEndpoint(t *testing.T) {
	t.Setenv(OTLPEndpointEnv, otlpEndpointDefault)
	ctx := context.Background()
	exp, err := NewOTLPLoggerExporters(ctx)
	require.NoError(t, err)
	assert.NotNil(t, exp)
}

func TestNewOTLPLoggerExporterWithoutEndpoint(t *testing.T) {
	t.Setenv(OTLPEndpointEnv, "")
	ctx := context.Background()
	exp, err := NewOTLPLoggerExporters(ctx)
	require.NoError(t, err)
	assert.NotNil(t, exp)
}

func TestNewOTLPLoggerExporterWithoutEndpointAndConsoleDisabled_ReturnsNoop(t *testing.T) {
	t.Setenv(OTLPEndpointEnv, "")
	t.Setenv(otlpLoggerEndpointKey, "")
	t.Setenv(otlpConsoleLogEnv, "false")
	t.Setenv(consoleLoggerRedirectEnv, "false")

	ctx := context.Background()
	exp, err := NewOTLPLoggerExporters(ctx)
	assert.NoError(t, err)
	assert.Len(t, exp, 1)
	assert.IsType(t, &noopLoggerExporter{}, exp[0])
}
