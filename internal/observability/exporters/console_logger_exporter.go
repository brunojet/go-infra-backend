package exporters

import (
	"os"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

// NewConsoleLoggerExporter creates an OpenTelemetry log exporter that writes
// records to stdout. Intended for development/debugging.
func NewConsoleLoggerExporter() (sdklog.Exporter, error) {
	return stdoutlog.New(
		stdoutlog.WithWriter(os.Stdout),
		// Compact output (no pretty print) to keep console logs readable.
		// Keep timestamps enabled by default.
	)
}
