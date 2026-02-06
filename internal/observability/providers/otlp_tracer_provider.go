package providers

import (
	"context"
	"log"

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
			if err := ff.ForceFlush(ctx); err != nil {
				log.Printf("otel tracer force flush error: %v", err)
			}
		}

		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("otel tracer shutdown error: %v", err)
			return err
		}
		return nil
	}

	return tp, shutdown, nil
}
