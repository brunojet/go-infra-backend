package providers

import (
	"context"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// NewOTLPTracerProvider builds and registers a TracerProvider using the OTLP
// tracer exporter and the resource from NewOTLPResource. Returns the provider
// and a shutdown function.
func NewOTLPTracerProvider(ctx context.Context, exporter sdktrace.SpanExporter) (*sdktrace.TracerProvider, func(context.Context) error, error) {
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
	)

	otel.SetTracerProvider(tp)

	shutdown := func(ctx context.Context) error {
		if ff, ok := interface{}(tp).(interface{ ForceFlush(context.Context) error }); ok {
			ff.ForceFlush(ctx)
		}
		tp.Shutdown(ctx)
		return nil
	}

	return tp, shutdown, nil
}
