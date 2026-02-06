package providers

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

// NewOTLPLoggerProvider builds and registers a LoggerProvider using the
// provided exporter and the resource from NewOTLPResource. The exporter must
// implement sdklog.Exporter (allows passing noop or real exporters in tests).
func NewOTLPLoggerProvider(ctx context.Context, exporters ...sdklog.Exporter) (*sdklog.LoggerProvider, func(context.Context) error, error) {
	_ = ctx

	processors := make([]sdklog.Processor, 0, len(exporters))
	for _, exp := range exporters {
		if exp == nil {
			continue
		}
		processors = append(processors, sdklog.NewBatchProcessor(exp))
	}
	if len(processors) == 0 {
		return nil, nil, fmt.Errorf("at least one logger exporter is required")
	}

	// --- Logger Provider ---
	options := make([]sdklog.LoggerProviderOption, 0, len(processors))
	for _, p := range processors {
		options = append(options, sdklog.WithProcessor(p))
	}
	loggerProvider := sdklog.NewLoggerProvider(options...)

	// Register as global provider
	global.SetLoggerProvider(loggerProvider)

	shutdown := func(ctx context.Context) error {
		if ff, ok := interface{}(loggerProvider).(interface{ ForceFlush(context.Context) error }); ok {
			if err := ff.ForceFlush(ctx); err != nil {
				log.Printf("otel logger force flush error: %v", err)
			}
		}

		if err := loggerProvider.Shutdown(ctx); err != nil {
			log.Printf("otel logger shutdown error: %v", err)
			return err
		}
		return nil
	}

	return loggerProvider, shutdown, nil
}
