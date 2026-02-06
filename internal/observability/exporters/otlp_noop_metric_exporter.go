package exporters

import (
	"context"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type noopMetricExporter struct{}

func (n *noopMetricExporter) Temporality(ik sdkmetric.InstrumentKind) metricdata.Temporality {
	return metricdata.CumulativeTemporality
}

func (n *noopMetricExporter) Aggregation(ik sdkmetric.InstrumentKind) sdkmetric.Aggregation {
	return sdkmetric.DefaultAggregationSelector(ik)
}

func (n *noopMetricExporter) Export(ctx context.Context, rm *metricdata.ResourceMetrics) error {
	return nil
}

func (n *noopMetricExporter) ForceFlush(ctx context.Context) error {
	return nil
}

func (n *noopMetricExporter) Shutdown(ctx context.Context) error {
	return nil
}

// NewOTLPNoopMetricExporter returns a no-op metric exporter for tests.
func NewOTLPNoopMetricExporter() sdkmetric.Exporter {
	return &noopMetricExporter{}
}
