package exporters

import (
	internalexporters "github.com/brunojet/go-infra-backend/internal/observability/exporters"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// NewOTLPNoopMetricExporter returns a no-op metric exporter useful for tests.
func NewOTLPNoopMetricExporter() sdkmetric.Exporter {
	return internalexporters.NewOTLPNoopMetricExporter()
}
