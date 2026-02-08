package exporters

import (
	internalexporters "github.com/brunojet/go-infra-backend/internal/observability/exporters"

	sdklog "go.opentelemetry.io/otel/sdk/log"
)

// NewOTLPNoopLoggerExporter returns a no-op logger exporter useful for tests.
func NewOTLPNoopLoggerExporter() sdklog.Exporter {
	return internalexporters.NewOTLPNoopLoggerExporter()
}
