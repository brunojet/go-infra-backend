package providers

import (
	"context"

	internalproviders "github.com/brunojet/go-infra-backend/internal/observability/providers"
	pubboot "github.com/brunojet/go-infra-backend/pkg/bootstrap"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// NewOTLPLoggerProvider builds and registers a LoggerProvider using the provided
// exporters. The exporter must implement sdklog.Exporter (allows passing noop or
// real exporters in tests).
func NewOTLPLoggerProvider(ctx context.Context, exporters ...sdklog.Exporter) (*sdklog.LoggerProvider, pubboot.ShutdownFunc, error) {
	lp, shutdown, err := internalproviders.NewOTLPLoggerProvider(ctx, exporters...)
	return lp, shutdown, err
}

// NewOTLPMetricProvider builds and registers a MeterProvider using the provided
// OTLP metric exporter.
func NewOTLPMetricProvider(ctx context.Context, exporter sdkmetric.Exporter) (*sdkmetric.MeterProvider, pubboot.ShutdownFunc, error) {
	mp, shutdown, err := internalproviders.NewOTLPMetricProvider(ctx, exporter)
	return mp, shutdown, err
}

// NewOTLPTracerProvider builds and registers a TracerProvider using the provided
// OTLP tracer exporter.
func NewOTLPTracerProvider(ctx context.Context, exporter sdktrace.SpanExporter) (*sdktrace.TracerProvider, pubboot.ShutdownFunc, error) {
	tp, shutdown, err := internalproviders.NewOTLPTracerProvider(ctx, exporter)
	return tp, shutdown, err
}
