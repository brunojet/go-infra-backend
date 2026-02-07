package bootstrap

import (
	"context"

	bootcontracts "github.com/brunojet/go-infra-backend/internal/bootstrap/contracts"
	"github.com/brunojet/go-infra-backend/internal/observability/exporters"
	"github.com/brunojet/go-infra-backend/internal/observability/providers"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type loggerExporterFunc func(ctx context.Context) ([]sdklog.Exporter, error)
type metricExporterFunc func(ctx context.Context) (sdkmetric.Exporter, error)
type tracerExporterFunc func(ctx context.Context) (sdktrace.SpanExporter, error)

type loggerProviderFunc func(ctx context.Context, exporters ...sdklog.Exporter) (*sdklog.LoggerProvider, func(context.Context) error, error)
type metricProviderFunc func(ctx context.Context, exporter sdkmetric.Exporter) (*sdkmetric.MeterProvider, func(context.Context) error, error)
type tracerProviderFunc func(ctx context.Context, exporter sdktrace.SpanExporter) (*sdktrace.TracerProvider, func(context.Context) error, error)

type InitLoggerFuncs struct {
	ExporterFunc loggerExporterFunc
	ProviderFunc loggerProviderFunc
}

type InitMetricFuncs struct {
	ExporterFunc metricExporterFunc
	ProviderFunc metricProviderFunc
}

type InitTracerFuncs struct {
	ExporterFunc tracerExporterFunc
	ProviderFunc tracerProviderFunc
}

type InitObservabilityFuncs struct {
	InitLoggerFuncs *InitLoggerFuncs
	InitMetricFuncs *InitMetricFuncs
	InitTracerFuncs *InitTracerFuncs
}

func initLogger(sm bootcontracts.ShutdownManager, loggerFuncs *InitLoggerFuncs) error {
	if loggerFuncs == nil {
		return nil
	}
	ctx := sm.GetContext()
	exportersSlice, err := loggerFuncs.ExporterFunc(ctx)
	if err != nil {
		return err
	}
	_, shutdown, err := loggerFuncs.ProviderFunc(ctx, exportersSlice...)
	if err != nil {
		return err
	}
	sm.RegisterFunc("otel-logger", shutdown)
	return nil
}

func initMetrics(sm bootcontracts.ShutdownManager, metricFuncs *InitMetricFuncs) error {
	if metricFuncs == nil {
		return nil
	}
	ctx := sm.GetContext()

	metricExporter, err := metricFuncs.ExporterFunc(ctx)
	if err != nil {
		return err
	}
	_, metricShutdown, err := metricFuncs.ProviderFunc(ctx, metricExporter)
	if err != nil {
		return err
	}
	sm.RegisterFunc("otel-metric", metricShutdown)
	return nil
}

func initTracing(sm bootcontracts.ShutdownManager, tracerFuncs *InitTracerFuncs) error {
	if tracerFuncs == nil {
		return nil
	}
	ctx := sm.GetContext()

	spanExporter, err := tracerFuncs.ExporterFunc(ctx)
	if err != nil {
		return err
	}
	_, tracerShutdown, err := tracerFuncs.ProviderFunc(ctx, spanExporter)
	if err != nil {
		return err
	}
	sm.RegisterFunc("otel-tracer", tracerShutdown)
	return nil
}

func initObservability(sm bootcontracts.ShutdownManager, initFuncs InitObservabilityFuncs) error {
	if err := initLogger(sm, initFuncs.InitLoggerFuncs); err != nil {
		return err
	}
	if err := initMetrics(sm, initFuncs.InitMetricFuncs); err != nil {
		return err
	}
	if err := initTracing(sm, initFuncs.InitTracerFuncs); err != nil {
		return err
	}

	return nil
}

func InitLogger(sm bootcontracts.ShutdownManager) error {
	return initLogger(sm, &InitLoggerFuncs{
		ExporterFunc: exporters.NewOTLPLoggerExporter,
		ProviderFunc: providers.NewOTLPLoggerProvider,
	})
}

func InitMetrics(sm bootcontracts.ShutdownManager) error {
	return initMetrics(sm, &InitMetricFuncs{
		ExporterFunc: exporters.NewOTLPMetricExporter,
		ProviderFunc: providers.NewOTLPMetricProvider,
	})
}

func InitTracing(sm bootcontracts.ShutdownManager) error {
	return initTracing(sm, &InitTracerFuncs{
		ExporterFunc: exporters.NewOTLPTracerExporter,
		ProviderFunc: providers.NewOTLPTracerProvider,
	})
}

func InitObservability(sm bootcontracts.ShutdownManager) error {
	return initObservability(sm, InitObservabilityFuncs{
		InitLoggerFuncs: &InitLoggerFuncs{
			ExporterFunc: exporters.NewOTLPLoggerExporter,
			ProviderFunc: providers.NewOTLPLoggerProvider,
		},
		InitMetricFuncs: &InitMetricFuncs{
			ExporterFunc: exporters.NewOTLPMetricExporter,
			ProviderFunc: providers.NewOTLPMetricProvider,
		},
		InitTracerFuncs: &InitTracerFuncs{
			ExporterFunc: exporters.NewOTLPTracerExporter,
			ProviderFunc: providers.NewOTLPTracerProvider,
		},
	})
}
