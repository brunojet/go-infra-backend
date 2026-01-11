package contracts

import (
	"context"
	"time"

	types "github.com/brunojet/go-infra-backend/internal/observability/types"
)

// Span represents an in-flight operation.
//
// End must be safe to call exactly once.
type Span interface {
	End(err error, fields ...types.Field)
}

// Tracer starts spans.
type Tracer interface {
	Start(ctx context.Context, name string, fields ...types.Field) (context.Context, Span)
}

// Metrics records counters and durations.
type Metrics interface {
	Inc(name string, value int64, fields ...types.Field)
	ObserveDuration(name string, d time.Duration, fields ...types.Field)
}

// Logger writes application logs.
type Logger interface {
	Info(ctx context.Context, msg string, fields ...types.Field)
	Error(ctx context.Context, msg string, err error, fields ...types.Field)
}

// Provider groups telemetry primitives.
type Provider interface {
	Tracer() Tracer
	Metrics() Metrics
	Logger() Logger
}
