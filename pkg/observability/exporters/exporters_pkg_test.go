package exporters

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExporters_ConstructorsHappyPath(t *testing.T) {
	// Ensure global OTLP endpoint is not set so we use console/noop exporters.
	t.Setenv(OTLPEndpointEnv, "")

	ctx := context.Background()

	logs, err := NewOTLPLoggerExporters(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, logs)

	metricExp, err := NewOTLPMetricExporter(ctx)
	require.NoError(t, err)
	require.NotNil(t, metricExp)

	traceExp, err := NewOTLPTracerExporter(ctx)
	require.NoError(t, err)
	require.NotNil(t, traceExp)
}
