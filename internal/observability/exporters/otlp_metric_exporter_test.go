package exporters

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOTLPMetricExporterWithEndpoint(t *testing.T) {
	os.Setenv(OTLPEndpointEnv, otlpEndpointDefault)
	defer os.Unsetenv(OTLPEndpointEnv)
	ctx := context.Background()
	exp, err := NewOTLPMetricExporter(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, exp)
}

func TestNewOTLPMetricExporterWithoutEndpoint(t *testing.T) {
	os.Unsetenv(OTLPEndpointEnv)
	ctx := context.Background()
	exp, err := NewOTLPMetricExporter(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, exp)
}
