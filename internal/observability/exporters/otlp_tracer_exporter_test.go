package exporters

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOTLPTracerExporterWithEndpoint(t *testing.T) {
	os.Setenv(OTLPEndpointEnv, otlpEndpointDefault)
	defer os.Unsetenv(OTLPEndpointEnv)
	ctx := context.Background()
	exp, err := NewOTLPTracerExporter(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, exp)
}

func TestNewOTLPTracerExporterWithoutEndpoint(t *testing.T) {
	os.Unsetenv(OTLPEndpointEnv)
	ctx := context.Background()
	exp, err := NewOTLPTracerExporter(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, exp)
}
