package exporters

import (
	"context"
	"testing"

	"os"

	"github.com/stretchr/testify/assert"
)

func TestNewOTLPLoggerExporterWithEndpoint(t *testing.T) {
	os.Setenv(OTLPEndpointEnv, OTLPEndpointDefault)
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
