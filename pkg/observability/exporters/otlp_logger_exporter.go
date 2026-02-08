package exporters

import (
	"context"

	internalexporters "github.com/brunojet/go-infra-backend/internal/observability/exporters"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

// NewOTLPLoggerExporter delegates to the internal implementation.
func NewOTLPLoggerExporter(ctx context.Context) ([]sdklog.Exporter, error) {
	return internalexporters.NewOTLPLoggerExporter(ctx)
}

// NewOTLPLoggerExporters is kept for compatibility with early pkg experiments.
// Prefer NewOTLPLoggerExporter.
func NewOTLPLoggerExporters(ctx context.Context) ([]sdklog.Exporter, error) {
	return NewOTLPLoggerExporter(ctx)
}
