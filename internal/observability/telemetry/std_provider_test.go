package telemetry

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brunojet/go-infra-backend/internal/observability/requestid"
	"github.com/brunojet/go-infra-backend/internal/observability/types"
)

func TestStdProvider_TracerMetricsLogger(t *testing.T) {
	p := NewStdProvider()
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

func TestStdTracerAndSpan(t *testing.T) {
	tracer := NewStdProvider().Tracer()
	ctx := context.Background()
	ctx2, span := tracer.Start(ctx, "test-span")
	if ctx2 == nil || span == nil {
		t.Error("Start should return context and span")
	}
	// End with no error
	span.End(nil)
	// End with error
	err := errors.New("fail")
	span.End(err)
}

func TestStdMetrics(t *testing.T) {
	metrics := NewStdProvider().Metrics()
	metrics.Inc("metric.name", 1)
	metrics.ObserveDuration("metric.name", 10*time.Millisecond)
}

func TestStdLogger_InfoError(t *testing.T) {
	logger := NewStdProvider().Logger()
	logger.Info(context.Background(), "info")
	logger.Error(context.Background(), "err", errors.New("fail"))
	logger.Error(context.Background(), "err-nil", nil)
}

func TestFormatFields(t *testing.T) {
	ctx := context.Background()
	// No request id, no fields
	if s := formatFields(ctx); s != "" {
		t.Errorf("expected empty string, got %q", s)
	}
	// With request id
	ctx = requestid.With(ctx, "abc123")
	if s := formatFields(ctx); s != " rid=abc123" {
		t.Errorf("expected rid, got %q", s)
	}
	// With fields
	fields := []types.Field{{Key: "foo", Value: 42}, {Key: "bar", Value: "baz"}}
	if s := formatFields(ctx, fields...); s != " rid=abc123 fields=foo:42,bar:baz" {
		t.Errorf("unexpected fields: %q", s)
	}
}
