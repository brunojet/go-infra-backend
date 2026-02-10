package exporters

import (
	"context"
	"errors"
	"log"

	"github.com/brunojet/go-infra-backend/pkg/bootstrap/contracts"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

var otlploggrpcNew = otlploggrpc.New
var consoleLoggerExporterNew = NewConsoleLoggerExporter
var logPrintf = log.Printf

func logNewOTLPLoggerResult(cfg OTLPLoggerConfig, exportersCount int, exporterErr error, providerErr error) {
	logPrintf(
		"NewOTLPLogger: exporters=%d endpoint_set=%t console_enabled=%t exporter_err=%v provider_err=%v",
		exportersCount,
		cfg.Endpoint != "",
		cfg.ConsoleLogEnable,
		exporterErr,
		providerErr,
	)
}

func newEndpointExporter(ctx context.Context, cfg OTLPConfig) (sdklog.Exporter, error) {
	options := []otlploggrpc.Option{otlploggrpc.WithEndpoint(cfg.Endpoint)}
	if cfg.IsInsecure {
		options = append(options, otlploggrpc.WithInsecure())
	}
	return otlploggrpcNew(ctx, options...)
}

func buildLoggerExporters(ctx context.Context, cfg OTLPLoggerConfig) ([]sdklog.Exporter, error) {
	exporters := make([]sdklog.Exporter, 0, 2)
	var errs error

	if cfg.Endpoint != "" {
		exp, err := newEndpointExporter(ctx, cfg.OTLPConfig)
		errs = errors.Join(errs, err)
		if err == nil {
			exporters = append(exporters, exp)
		}
	}

	if cfg.ConsoleLogEnable {
		exp, err := consoleLoggerExporterNew(cfg.ConsoleRedirectEnable)
		errs = errors.Join(errs, err)
		if err == nil {
			exporters = append(exporters, exp)
		}
	}

	return exporters, errs
}

func installLoggerProvider(exporters []sdklog.Exporter) (*otlpLogger, error) {
	if len(exporters) == 0 {
		return nil, errors.New("no logger exporters configured")
	}

	options := make([]sdklog.LoggerProviderOption, 0, len(exporters))
	for _, exp := range exporters {
		options = append(options, sdklog.WithProcessor(sdklog.NewBatchProcessor(exp)))
	}

	loggerProvider := sdklog.NewLoggerProvider(options...)
	global.SetLoggerProvider(loggerProvider)
	return &otlpLogger{provider: loggerProvider}, nil
}

type otlpLogger struct {
	provider *sdklog.LoggerProvider
}

func NewOTLPLogger(ctx context.Context, cfg OTLPLoggerConfig) (contracts.Shutdown, error) {
	exporters, exporterErr := buildLoggerExporters(ctx, cfg)
	provider, err := installLoggerProvider(exporters)
	logNewOTLPLoggerResult(cfg, len(exporters), exporterErr, err)
	return provider, err
}

func (o *otlpLogger) Shutdown(ctx context.Context) error {
	return o.provider.Shutdown(ctx)
}

func NewOTLPLoggerFromEnv(ctx context.Context) (contracts.Shutdown, error) {
	cfg := getLoggerExporterConfigFromEnv()
	return NewOTLPLogger(ctx, cfg)
}
