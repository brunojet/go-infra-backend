package exporters

import (
	"context"

	internalexporters "github.com/brunojet/go-infra-backend/internal/observability/exporters"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Re-export the env var key so external apps can configure OTLP without touching internal.
const OTLPEndpointEnv = internalexporters.OTLPEndpointEnv

func NewOTLPLoggerExporters(ctx context.Context) ([]sdklog.Exporter, error) {
	return internalexporters.NewOTLPLoggerExporters(ctx)
}

func NewOTLPMetricExporter(ctx context.Context) (sdkmetric.Exporter, error) {
	return internalexporters.NewOTLPMetricExporter(ctx)
}

func NewOTLPTracerExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	return internalexporters.NewOTLPTracerExporter(ctx)
}
