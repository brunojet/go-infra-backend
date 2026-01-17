package contracts

import "context"

// Fields represents structured logging fields.
type Fields map[string]interface{}

// Logger is an implementation-agnostic logging interface suitable for core code
// and adapters (zerolog, logrus, etc.). Methods now accept a `context.Context`
// so adapters can extract trace/span identifiers and attach them to log entries
// (necessary to correlate logs with traces and metrics).
type Logger interface {
	Debug(ctx context.Context, msg string, fields Fields)
	Info(ctx context.Context, msg string, fields Fields)
	Warn(ctx context.Context, msg string, fields Fields)
	Error(ctx context.Context, msg string, err error, fields Fields)

	// WithFields returns a derived logger that includes the provided fields on every entry.
	// The derived logger retains methods that accept a context.
	WithFields(fields Fields) Logger
	// WithContext returns a derived logger enriched with fields extracted from the
	// provided context (for example trace_id/span_id). This is a convenience for
	// call sites that prefer not to pass `ctx` on every logging call.
	WithContext(ctx context.Context) Logger
}
