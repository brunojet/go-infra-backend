package contracts

import "context"

// MetricsExporter defines the port for exporting metrics.
type MetricsExporter interface {
	// Start initializes the exporter (connects to backend).
	Start(ctx context.Context) error
	// Shutdown gracefully stops the exporter.
	Shutdown(ctx context.Context) error
	// Record sends a single metric value.
	Record(ctx context.Context, name string, value float64) error
}

// TracingExporter defines the port for exporting traces.
type TracingExporter interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
	// StartSpan starts a span and returns a context containing it plus a finish function.
	StartSpan(ctx context.Context, name string) (context.Context, func())
}

// LoggingExporter defines the port for exporting logs.
type LoggingExporter interface {
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
	// Emit emits a log record with a severity and body.
	Emit(ctx context.Context, severity string, body string) error
}
