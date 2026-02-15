package adapters

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/brunojet/go-infra-backend/pkg/bootstrap/contracts"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

var (
	otlploggrpcNew = otlploggrpc.New
	stdoutlogNew   = stdoutlog.New
)

type otlpLogger struct {
	provider   *sdklog.LoggerProvider
	redirector *stdlogRedirector
}

func logNewOTLPLoggerResult(cfg OTLPLoggerConfig, exportersCount int, exporterErr error, providerErr error) {
	log.Printf(
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

func newConsoleExporter() (sdklog.Exporter, error) {
	return stdoutlogNew(stdoutlog.WithWriter(os.Stdout))
}

func buildLoggerProviderOptions(ctx context.Context, cfg OTLPLoggerConfig) ([]sdklog.LoggerProviderOption, error) {
	options := make([]sdklog.LoggerProviderOption, 0, 2)
	var errs error

	if cfg.Endpoint != "" {
		exp, err := newEndpointExporter(ctx, cfg.OTLPConfig)
		errs = errors.Join(errs, err)
		if err == nil {
			options = append(options, sdklog.WithProcessor(sdklog.NewBatchProcessor(exp)))
		}
	}

	if cfg.ConsoleLogEnable {
		exp, err := newConsoleExporter()
		errs = errors.Join(errs, err)
		if err == nil {
			options = append(options, sdklog.WithProcessor(sdklog.NewSimpleProcessor(exp)))
		}
	}

	return options, errs
}

func installLoggerProvider(options []sdklog.LoggerProviderOption, cfg OTLPLoggerConfig) (*otlpLogger, error) {
	if len(options) == 0 {
		return nil, errors.New("no logger exporters configured")
	}

	loggerProvider := sdklog.NewLoggerProvider(options...)
	global.SetLoggerProvider(loggerProvider)

	if cfg.ConsoleRedirectEnable {
		return &otlpLogger{
			provider:   loggerProvider,
			redirector: NewStdLogRedirector("stdlog", otellog.SeverityInfo),
		}, nil
	}

	return &otlpLogger{provider: loggerProvider}, nil
}

func NewOTLPLogger(ctx context.Context, cfg OTLPLoggerConfig) (contracts.Shutdown, error) {
	options, exporterErr := buildLoggerProviderOptions(ctx, cfg)
	provider, err := installLoggerProvider(options, cfg)
	logNewOTLPLoggerResult(cfg, len(options), exporterErr, err)
	return provider, err
}

func NewOTLPLoggerFromEnv(ctx context.Context) (contracts.Shutdown, error) {
	cfg := getLoggerExporterConfigFromEnv()
	return NewOTLPLogger(ctx, cfg)
}

func (o *otlpLogger) Shutdown(ctx context.Context) error {
	o.provider.ForceFlush(ctx)
	if o.redirector != nil {
		o.redirector.Shutdown(ctx)
	}
	return o.provider.Shutdown(ctx)
}
