package exporters

import (
	"context"

	internalexporters "github.com/brunojet/go-infra-backend/internal/observability/exporters"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// NewOTLPTracerExporter creates an OTLP trace exporter configured from env.
// If OTLP is not configured, it returns a noop exporter.
func NewOTLPTracerExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	return internalexporters.NewOTLPTracerExporter(ctx)
}
