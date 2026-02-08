package providers

import (
	"context"
	"testing"
	"time"

	pubexp "github.com/brunojet/go-infra-backend/pkg/observability/exporters"
	"github.com/stretchr/testify/require"
)

func TestProviders_ConstructorsAndShutdownHappyPath(t *testing.T) {
	// Ensure global OTLP endpoint is not set so we use console/noop exporters.
	t.Setenv(pubexp.OTLPEndpointEnv, "")

	ctx := context.Background()

	logExporters, err := pubexp.NewOTLPLoggerExporters(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, logExporters)

	lp, logShutdown, err := NewOTLPLoggerProvider(ctx, logExporters...)
	require.NoError(t, err)
	require.NotNil(t, lp)
	require.NotNil(t, logShutdown)

	metricExporter, err := pubexp.NewOTLPMetricExporter(ctx)
	require.NoError(t, err)
	require.NotNil(t, metricExporter)

	mp, metricShutdown, err := NewOTLPMetricProvider(ctx, metricExporter)
	require.NoError(t, err)
	require.NotNil(t, mp)
	require.NotNil(t, metricShutdown)

	traceExporter, err := pubexp.NewOTLPTracerExporter(ctx)
	require.NoError(t, err)
	require.NotNil(t, traceExporter)

	tp, traceShutdown, err := NewOTLPTracerProvider(ctx, traceExporter)
	require.NoError(t, err)
	require.NotNil(t, tp)
	require.NotNil(t, traceShutdown)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Shutdown order doesn't matter for this unit test; we just verify they succeed.
	require.NoError(t, traceShutdown(shutdownCtx))
	require.NoError(t, metricShutdown(shutdownCtx))
	require.NoError(t, logShutdown(shutdownCtx))
}
