package exporters

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// NewOTLPTracerExporter cria um otlptrace.Exporter a partir de um OTLPClient.
// Retorna como sdktrace.SpanExporter para desacoplar callers do tipo concreto.
func NewOTLPTracerExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	exporterConfig := GetExporterConfigFromEnv()

	if exporterConfig.Endpoint == "" {
		return NewOTLPNoopTracerExporter(), nil
	}

	return otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(exporterConfig.Endpoint),
	)
}
