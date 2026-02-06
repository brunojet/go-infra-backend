package exporters

import (
	"context"

	sdklog "go.opentelemetry.io/otel/sdk/log"
)

type noopLoggerExporter struct{}

func (n *noopLoggerExporter) Export(ctx context.Context, records []sdklog.Record) error {
	return nil
}

func (n *noopLoggerExporter) Shutdown(ctx context.Context) error {
	return nil
}

func (n *noopLoggerExporter) ForceFlush(ctx context.Context) error {
	return nil
}

// NewOTLPNoopLoggerExporter returns a no-op logger exporter for tests.
func NewOTLPNoopLoggerExporter() sdklog.Exporter {
	return &noopLoggerExporter{}
}
