package adapters

import (
	"context"

	"go.opentelemetry.io/otel"
	otlptracegrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"google.golang.org/grpc"

	"github.com/brunojet/go-infra-backend/internal/observability/exporters/contracts"
)

// OTLPTracingExporter implements contracts.TracingExporter using OpenTelemetry OTLP.
type OTLPTracingExporter struct {
	tp   *sdktrace.TracerProvider
	conn *grpc.ClientConn
}

// NewOTLPTracingExporter returns a new exporter configured to send to the given endpoint (host:port).
func NewOTLPTracingExporter(endpoint string) contracts.TracingExporter {
	return &OTLPTracingExporter{tp: nil}
}

// Start initializes the OTLP tracing exporter and creates a TracerProvider.
func (o *OTLPTracingExporter) Start(ctx context.Context) error {
	endpoint := EndpointFromEnv()
	// Try to dial OTLP collector; if it fails, fallback to local provider.
	cc, err := DialOTLP(ctx, endpoint)
	res, _ := sdkresource.Merge(sdkresource.Default(), sdkresource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String("go-infra-backend")))
	if err == nil && cc != nil {
		exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(cc))
		if err == nil {
			tp := sdktrace.NewTracerProvider(
				sdktrace.WithBatcher(exporter),
				sdktrace.WithResource(res),
			)
			o.tp = tp
			o.conn = cc
			otel.SetTracerProvider(tp)
			return nil
		}
		// close conn if exporter creation failed
		_ = cc.Close()
	}

	// Fallback: local provider without network exporter.
	tp := sdktrace.NewTracerProvider(sdktrace.WithResource(res))
	o.tp = tp
	otel.SetTracerProvider(tp)
	return nil
}

// Shutdown stops the tracer provider.
func (o *OTLPTracingExporter) Shutdown(ctx context.Context) error {
	if o.tp == nil {
		return nil
	}
	if err := o.tp.Shutdown(ctx); err != nil {
		return err
	}
	if o.conn != nil {
		return o.conn.Close()
	}
	return nil
}

// StartSpan starts a span and returns context and finish function.
func (o *OTLPTracingExporter) StartSpan(ctx context.Context, name string) (context.Context, func()) {
	if o.tp == nil {
		// create a noop tracer
		tracer := otel.Tracer("noop")
		ctx, span := tracer.Start(ctx, name)
		return ctx, func() { span.End() }
	}
	tracer := o.tp.Tracer("otlp-exporter")
	ctx, span := tracer.Start(ctx, name)
	return ctx, func() { span.End() }
}

var _ contracts.TracingExporter = (*OTLPTracingExporter)(nil)
