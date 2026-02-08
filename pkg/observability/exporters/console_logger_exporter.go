package exporters

import (
	"io"

	internalexporters "github.com/brunojet/go-infra-backend/internal/observability/exporters"
	otellog "go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

// NewStdLogWriter delegates to the internal implementation.
func NewStdLogWriter(loggerName string, severity otellog.Severity) io.Writer {
	return internalexporters.NewStdLogWriter(loggerName, severity)
}

// NewConsoleLoggerExporter delegates to the internal implementation.
func NewConsoleLoggerExporter(redirectLog bool) (sdklog.Exporter, error) {
	return internalexporters.NewConsoleLoggerExporter(redirectLog)
}
