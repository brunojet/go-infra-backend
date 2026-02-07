package providers

import (
	"context"
	"fmt"

	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

func activateExporterIfNeeded(ctx context.Context, exporters []sdklog.Exporter, loggerName string, severity otellog.Severity) {
	for _, exporter := range exporters {
		if a, ok := exporter.(interface {
			Activate(context.Context, string, otellog.Severity)
		}); ok {
			a.Activate(ctx, loggerName, severity)
		}
	}
}

func deactivateExporterIfNeeded(ctx context.Context, exporters []sdklog.Exporter) {
	for _, exporter := range exporters {
		if d, ok := exporter.(interface{ Deactivate(context.Context) }); ok {
			d.Deactivate(ctx)
		}
	}
}

// NewOTLPLoggerProvider builds and registers a LoggerProvider using the
// provided exporter and the resource from NewOTLPResource. The exporter must
// implement sdklog.Exporter (allows passing noop or real exporters in tests).
func NewOTLPLoggerProvider(ctx context.Context, exporters ...sdklog.Exporter) (*sdklog.LoggerProvider, func(context.Context) error, error) {
	// --- Logger Provider ---
	options := make([]sdklog.LoggerProviderOption, 0, len(exporters))
	for _, exp := range exporters {
		options = append(options, sdklog.WithProcessor(sdklog.NewBatchProcessor(exp)))
	}
	if len(options) == 0 {
		return nil, nil, fmt.Errorf("at least one logger exporter is required")
	}
	loggerProvider := sdklog.NewLoggerProvider(options...)

	// Register as global provider
	global.SetLoggerProvider(loggerProvider)

	activateExporterIfNeeded(ctx, exporters, "otel-logger", otellog.SeverityInfo)

	shutdown := func(ctx context.Context) error {
		if ff, ok := interface{}(loggerProvider).(interface{ ForceFlush(context.Context) error }); ok {
			_ = ff.ForceFlush(ctx)
		}

		deactivateExporterIfNeeded(ctx, exporters)

		_ = loggerProvider.Shutdown(ctx)

		return nil
	}

	return loggerProvider, shutdown, nil
}
