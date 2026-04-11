package adapters

import (
	"context"
	"errors"
	"log"

	"github.com/brunojet/go-infra-backend/pkg/bootstrap/contracts"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var otlptracegrpcNew = otlptracegrpc.New

type otlpTracer struct {
	provider *sdktrace.TracerProvider
}

func logNewOTLPTracerResult(cfg OTLPConfig, tracerErr error) {
	log.Printf(
		"NewOTLPTracer: endpoint_set=%t insecure=%t tracerErr=%v",
		cfg.Endpoint != "",
		cfg.IsInsecure,
		tracerErr,
	)
}

func newTraceExporter(ctx context.Context, cfg OTLPConfig) (sdktrace.SpanExporter, error) {
	options := []otlptracegrpc.Option{otlptracegrpc.WithEndpoint(cfg.Endpoint)}
	if cfg.IsInsecure {
		options = append(options, otlptracegrpc.WithInsecure())
	}
	return otlptracegrpcNew(ctx, options...)
}

func installTraceProvider(exporter sdktrace.SpanExporter) (*otlpTracer, error) {
	if exporter == nil {
		return nil, errors.New("nil exporter provided")
	}
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
	)
	otel.SetTracerProvider(provider)
	return &otlpTracer{provider: provider}, nil
}

func NewOTLPTracer(ctx context.Context, cfg OTLPConfig) (contracts.Shutdown, error) {
	var tracerErr error
	var tracer *otlpTracer
	exporter, tracerErr := newTraceExporter(ctx, cfg)
	if tracerErr == nil {
		tracer, tracerErr = installTraceProvider(exporter)
	}
	logNewOTLPTracerResult(cfg, tracerErr)
	return tracer, tracerErr
}

func NewOTLPTracerFromEnv(ctx context.Context) (contracts.Shutdown, error) {
	cfg := getTraceExporterConfigFromEnv()
	return NewOTLPTracer(ctx, cfg)
}

func (o *otlpTracer) Shutdown(ctx context.Context) error {
	o.provider.ForceFlush(ctx)
	return o.provider.Shutdown(ctx)
}
