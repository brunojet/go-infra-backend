package exporters

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestOTLPNoopMetricExporter_Basics(t *testing.T) {
	exp := NewOTLPNoopMetricExporter()
	require.NotNil(t, exp)

	// Temporality should return a valid temporality
	temp := exp.Temporality(sdkmetric.InstrumentKindCounter)
	require.Equal(t, metricdata.CumulativeTemporality, temp)

	// Aggregation should be non-nil (default selector)
	agg := exp.Aggregation(sdkmetric.InstrumentKindCounter)
	require.NotNil(t, agg)

	// Other methods should not error
	require.NoError(t, exp.Export(context.Background(), &metricdata.ResourceMetrics{}))
	require.NoError(t, exp.ForceFlush(context.Background()))
	require.NoError(t, exp.Shutdown(context.Background()))
}
