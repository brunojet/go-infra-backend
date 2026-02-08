package providers

import (
	"context"

	internalproviders "github.com/brunojet/go-infra-backend/internal/observability/providers"
	"github.com/brunojet/go-infra-backend/pkg/observability/contracts"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// NewOTLPMetricProvider builds and registers a MeterProvider using the provided
// OTLP metric exporter.
func NewOTLPMetricProvider(ctx context.Context, exporter sdkmetric.Exporter) (*sdkmetric.MeterProvider, contracts.ShutdownFunc, error) {
	mp, shutdown, err := internalproviders.NewOTLPMetricProvider(ctx, exporter)
	return mp, shutdown, err
}
