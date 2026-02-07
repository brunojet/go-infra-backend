package exporters

import (
	"context"
	"testing"

	"os"

	"github.com/stretchr/testify/assert"
)

func TestNewOTLPLoggerExporterWithEndpoint(t *testing.T) {
	os.Setenv(OTLPEndpointEnv, otlpEndpointDefault)
	defer os.Unsetenv(OTLPEndpointEnv)
	ctx := context.Background()
	exp, err := NewOTLPLoggerExporter(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, exp)
}

func TestNewOTLPLoggerExporterWithoutEndpoint(t *testing.T) {
	os.Unsetenv(OTLPEndpointEnv)
	ctx := context.Background()
	exp, err := NewOTLPLoggerExporter(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, exp)
}

func TestNewOTLPLoggerExporterWithoutEndpointAndConsoleDisabled_ReturnsNoop(t *testing.T) {
	t.Setenv(OTLPEndpointEnv, "")
	t.Setenv(otlpLoggerEndpointKey, "")
	t.Setenv(otlpConsoleLogEnv, "false")
	t.Setenv(consoleLoggerRedirectEnv, "false")

	ctx := context.Background()
	exp, err := NewOTLPLoggerExporter(ctx)
	assert.NoError(t, err)
	assert.Len(t, exp, 1)
	assert.IsType(t, &noopLoggerExporter{}, exp[0])
}
