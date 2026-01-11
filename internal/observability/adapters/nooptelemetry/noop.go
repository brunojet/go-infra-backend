package noopadapter

import (
	"context"
	"time"

	"github.com/brunojet/go-infra-backend/internal/observability/contracts"
	"github.com/brunojet/go-infra-backend/internal/observability/types"
)

// NoopProvider is a safe default Provider implementation.
//
// It does nothing and can be used when no real telemetry is configured.
type NoopProvider struct{}

func (NoopProvider) Tracer() contracts.Tracer   { return noopTracer{} }
func (NoopProvider) Metrics() contracts.Metrics { return noopMetrics{} }
func (NoopProvider) Logger() contracts.Logger   { return noopLogger{} }

type noopTracer struct{}

type noopSpan struct{}

type noopMetrics struct{}

type noopLogger struct{}

func (noopTracer) Start(ctx context.Context, _ string, _ ...types.Field) (context.Context, contracts.Span) {
	return ctx, noopSpan{}
}

func (noopSpan) End(_ error, _ ...types.Field)              {}
func (noopMetrics) Inc(_ string, _ int64, _ ...types.Field) {}

func (noopMetrics) ObserveDuration(_ string, _ time.Duration, _ ...types.Field) {}

func (noopLogger) Info(_ context.Context, _ string, _ ...types.Field)           {}
func (noopLogger) Error(_ context.Context, _ string, _ error, _ ...types.Field) {}
