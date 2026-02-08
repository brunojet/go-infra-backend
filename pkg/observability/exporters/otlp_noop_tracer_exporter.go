package exporters

import (
	internalexporters "github.com/brunojet/go-infra-backend/internal/observability/exporters"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// NewOTLPNoopTracerExporter returns a no-op SpanExporter useful for tests.
func NewOTLPNoopTracerExporter() sdktrace.SpanExporter {
	return internalexporters.NewOTLPNoopTracerExporter()
}
