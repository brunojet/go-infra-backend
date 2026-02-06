package exporters

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

// NewOTLPLoggerExporter cria um otlplog.Exporter a partir de um OTLPClient.
// Retorna como sdklog.Exporter para desacoplar callers do tipo concreto.
func NewOTLPLoggerExporter(ctx context.Context) (sdklog.Exporter, error) {
	exporterConfig := GetExporterConfigFromEnv()

	if exporterConfig.Endpoint == "" {
		return NewOTLPNoopLoggerExporter(), nil
	}

	return otlploggrpc.New(
		ctx,
		otlploggrpc.WithInsecure(),
		otlploggrpc.WithEndpoint(exporterConfig.Endpoint),
	)
}
