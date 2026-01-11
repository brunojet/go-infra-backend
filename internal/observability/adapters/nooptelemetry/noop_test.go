package noopadapter

import (
	"context"
	"testing"
	"time"

	"github.com/brunojet/go-infra-backend/internal/observability/types"
)

func TestNoopProvider_TracerMetricsLogger(t *testing.T) {
	p := NoopProvider{}
	if p.Tracer() == nil {
		t.Error("Tracer should not be nil")
	}
	if p.Metrics() == nil {
		t.Error("Metrics should not be nil")
	}
	if p.Logger() == nil {
		t.Error("Logger should not be nil")
	}
}

func TestNoopTracerAndSpan(t *testing.T) {
	tracer := NoopProvider{}.Tracer()
	ctx := context.Background()
	ctx2, span := tracer.Start(ctx, "test-span")
	if ctx2 != ctx {
		t.Error("NoopTracer should return same context")
	}
	if span == nil {
		t.Error("NoopTracer should return non-nil span")
	}
	span.End(nil)
	span.End(nil, types.Field{Key: "foo", Value: 1})
}

func TestNoopMetrics(t *testing.T) {
	metrics := NoopProvider{}.Metrics()
	metrics.Inc("metric.name", 1)
	metrics.ObserveDuration("metric.name", 10*time.Millisecond)
}

func TestNoopLogger(t *testing.T) {
	logger := NoopProvider{}.Logger()
	logger.Info(context.Background(), "info")
	logger.Error(context.Background(), "err", nil)
}
