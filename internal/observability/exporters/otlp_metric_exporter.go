package exporters

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// NewOTLPMetricExporter cria um otlpmetric.Exporter a partir de um OTLPClient.
// Retorna como sdkmetric.Exporter para desacoplar callers do tipo concreto.
func NewOTLPMetricExporter(ctx context.Context) (sdkmetric.Exporter, error) {
	exporterConfig := GetExporterConfigFromEnv()

	if exporterConfig.Endpoint == "" {
		return NewOTLPNoopMetricExporter(), nil
	}

	return otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(exporterConfig.Endpoint),
	)
}
