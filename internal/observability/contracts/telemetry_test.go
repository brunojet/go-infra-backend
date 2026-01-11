package contracts

import (
	"context"
	"testing"
	"time"

	obsTypes "github.com/brunojet/go-infra-backend/internal/observability/types"
)

type dummySpan struct{}

func (d *dummySpan) End(err error, fields ...obsTypes.Field) {}

type dummyTracer struct{}

func (d *dummyTracer) Start(ctx context.Context, name string, fields ...obsTypes.Field) (context.Context, Span) {
	return ctx, &dummySpan{}
}

type dummyMetrics struct{}

func (d *dummyMetrics) Inc(name string, value int64, fields ...obsTypes.Field)                   {}
func (d *dummyMetrics) ObserveDuration(name string, dur time.Duration, fields ...obsTypes.Field) {}

type dummyLogger struct{}

func (d *dummyLogger) Info(ctx context.Context, msg string, fields ...obsTypes.Field)             {}
func (d *dummyLogger) Error(ctx context.Context, msg string, err error, fields ...obsTypes.Field) {}

type dummyProvider struct{}

func (d *dummyProvider) Tracer() Tracer   { return &dummyTracer{} }
func (d *dummyProvider) Metrics() Metrics { return &dummyMetrics{} }
func (d *dummyProvider) Logger() Logger   { return &dummyLogger{} }

func TestTelemetryInterfaces(t *testing.T) {
	var _ Span = &dummySpan{}
	var _ Tracer = &dummyTracer{}
	var _ Metrics = &dummyMetrics{}
	var _ Logger = &dummyLogger{}
	var _ Provider = &dummyProvider{}
}
