package exporters

import (
	"context"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type noopSpanExporter struct{}

func (n *noopSpanExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	return nil
}

func (n *noopSpanExporter) Shutdown(ctx context.Context) error {
	return nil
}

// NewOTLPNoopTracerExporter returns a no-op SpanExporter useful for tests.
func NewOTLPNoopTracerExporter() sdktrace.SpanExporter {
	return &noopSpanExporter{}
}
