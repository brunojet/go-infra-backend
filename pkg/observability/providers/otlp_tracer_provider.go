package providers

import (
	"context"

	internalproviders "github.com/brunojet/go-infra-backend/internal/observability/providers"
	"github.com/brunojet/go-infra-backend/pkg/observability/contracts"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// NewOTLPTracerProvider builds and registers a TracerProvider using the provided
// OTLP tracer exporter.
func NewOTLPTracerProvider(ctx context.Context, exporter sdktrace.SpanExporter) (*sdktrace.TracerProvider, contracts.ShutdownFunc, error) {
	tp, shutdown, err := internalproviders.NewOTLPTracerProvider(ctx, exporter)
	return tp, shutdown, err
}
