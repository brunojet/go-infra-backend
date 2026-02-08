package exporters

import (
	"context"
	"errors"
	"log"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

// NewOTLPLoggerExporters cria um otlplog.Exporter a partir de um OTLPClient.
// Retorna como sdklog.Exporter para desacoplar callers do tipo concreto.
func NewOTLPLoggerExporters(ctx context.Context) ([]sdklog.Exporter, error) {
	exporterConfig := GetExporterConfigFromEnv().loggerConfig

	exporters := make([]sdklog.Exporter, 0, 1)

	var errs error
	var err error

	if exporterConfig.EndpointLogEnable {
		var exporter sdklog.Exporter
		if exporter, err = otlploggrpc.New(ctx, otlploggrpc.WithInsecure(), otlploggrpc.WithEndpoint(exporterConfig.Endpoint)); err == nil {
			exporters = append(exporters, exporter)
		}
		errs = errors.Join(errs, err)
	}

	if exporterConfig.ConsoleLogEnable {
		var exporter sdklog.Exporter
		if exporter, err = NewConsoleLoggerExporter(exporterConfig.ConsoleRedirectEnable); err == nil {
			exporters = append(exporters, exporter)
		}
		errs = errors.Join(errs, err)
	}

	if len(exporters) == 0 && errs == nil {
		exporters = append(exporters, NewOTLPNoopLoggerExporter())
		log.Print("no logger exporter configured, using noop")
	}

	log.Printf("OTLP logger exporter configured. errorList=%v", errs)

	return exporters, errs
}
