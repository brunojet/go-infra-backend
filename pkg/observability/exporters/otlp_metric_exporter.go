package exporters

import (
	"context"

	internalexporters "github.com/brunojet/go-infra-backend/internal/observability/exporters"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// NewOTLPMetricExporter creates an OTLP metric exporter configured from env.
// If OTLP is not configured, it returns a noop exporter.
func NewOTLPMetricExporter(ctx context.Context) (sdkmetric.Exporter, error) {
	return internalexporters.NewOTLPMetricExporter(ctx)
}
