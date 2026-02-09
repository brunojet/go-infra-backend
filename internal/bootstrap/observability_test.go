package bootstrap

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestInitObservability_All_NoErrors(t *testing.T) {
	// For these unit tests, force exporters to use noop behavior.
	// internal/config.GetEnv treats empty values as unset.
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "")

	sm := NewShutdownManager(context.Background())

	assert := assert.New(t)
	assert.NoError(InitObservability(sm))
	assert.NoError(sm.ShutdownWithTimeout(2 * time.Second))
}

func TestInitLogger_NilFuncs_NoError(t *testing.T) {
	sm := NewShutdownManager(context.Background())
	require.NoError(t, initLogger(sm, nil))
}

func TestInitMetrics_NilFuncs_NoError(t *testing.T) {
	sm := NewShutdownManager(context.Background())
	require.NoError(t, initMetrics(sm, nil))
}

func TestInitTracing_NilFuncs_NoError(t *testing.T) {
	sm := NewShutdownManager(context.Background())
	require.NoError(t, initTracing(sm, nil))
}

func TestInitLogger_ExporterError(t *testing.T) {
	assert := assert.New(t)
	sm := NewShutdownManager(context.Background())

	bad := &InitLoggerFuncs{
		ExporterFunc: func(ctx context.Context) ([]sdklog.Exporter, error) {
			return nil, errors.New("exporter-fail")
		},
		ProviderFunc: func(ctx context.Context, exporters ...sdklog.Exporter) (*sdklog.LoggerProvider, func(context.Context) error, error) {
			return nil, func(ctx context.Context) error { return nil }, nil
		},
	}

	err := initLogger(sm, bad)
	assert.Error(err)
	assert.Contains(err.Error(), "exporter-fail")
}

func TestInitMetrics_ExporterError(t *testing.T) {
	assert := assert.New(t)
	sm := NewShutdownManager(context.Background())

	bad := &InitMetricFuncs{
		ExporterFunc: func(ctx context.Context) (sdkmetric.Exporter, error) {
			return nil, errors.New("metric-exporter-fail")
		},
		ProviderFunc: func(ctx context.Context, exporter sdkmetric.Exporter) (*sdkmetric.MeterProvider, func(context.Context) error, error) {
			return nil, func(ctx context.Context) error { return nil }, nil
		},
	}

	err := initMetrics(sm, bad)
	assert.Error(err)
	assert.Contains(err.Error(), "metric-exporter-fail")
}

func TestInitMetrics_ProviderError(t *testing.T) {
	assert := assert.New(t)
	sm := NewShutdownManager(context.Background())

	bad := &InitMetricFuncs{
		ExporterFunc: func(ctx context.Context) (sdkmetric.Exporter, error) {
			return &struct{ sdkmetric.Exporter }{}, nil
		},
		ProviderFunc: func(ctx context.Context, exporter sdkmetric.Exporter) (*sdkmetric.MeterProvider, func(context.Context) error, error) {
			return nil, nil, errors.New("metric-provider-fail")
		},
	}

	err := initMetrics(sm, bad)
	assert.Error(err)
	assert.Contains(err.Error(), "metric-provider-fail")
}

func TestInitTracing_ExporterError(t *testing.T) {
	assert := assert.New(t)
	sm := NewShutdownManager(context.Background())

	bad := &InitTracerFuncs{
		ExporterFunc: func(ctx context.Context) (sdktrace.SpanExporter, error) {
			return nil, errors.New("tracer-exporter-fail")
		},
		ProviderFunc: func(ctx context.Context, exporter sdktrace.SpanExporter) (*sdktrace.TracerProvider, func(context.Context) error, error) {
			return nil, nil, nil
		},
	}

	err := initTracing(sm, bad)
	assert.Error(err)
	assert.Contains(err.Error(), "tracer-exporter-fail")
}

func TestInitTracing_ProviderError(t *testing.T) {
	assert := assert.New(t)
	sm := NewShutdownManager(context.Background())

	bad := &InitTracerFuncs{
		ExporterFunc: func(ctx context.Context) (sdktrace.SpanExporter, error) {
			return &struct{ sdktrace.SpanExporter }{}, nil
		},
		ProviderFunc: func(ctx context.Context, exporter sdktrace.SpanExporter) (*sdktrace.TracerProvider, func(context.Context) error, error) {
			return nil, nil, errors.New("tracer-provider-fail")
		},
	}

	err := initTracing(sm, bad)
	assert.Error(err)
	assert.Contains(err.Error(), "tracer-provider-fail")
}

func TestInitLogger_ProviderError(t *testing.T) {
	assert := assert.New(t)
	sm := NewShutdownManager(context.Background())

	bad := &InitLoggerFuncs{
		ExporterFunc: func(ctx context.Context) ([]sdklog.Exporter, error) {
			return nil, nil
		},
		ProviderFunc: func(ctx context.Context, exporters ...sdklog.Exporter) (*sdklog.LoggerProvider, func(context.Context) error, error) {
			return nil, nil, errors.New("logger-provider-fail")
		},
	}

	err := initLogger(sm, bad)
	assert.Error(err)
	assert.Contains(err.Error(), "logger-provider-fail")
}

func TestInitObservability_LoggerErrorStopsSequence(t *testing.T) {
	assert := assert.New(t)
	sm := NewShutdownManager(context.Background())

	metricCalled := false
	tracerCalled := false

	initFuncs := InitObservabilityFuncs{
		InitLoggerFuncs: &InitLoggerFuncs{
			ExporterFunc: func(ctx context.Context) ([]sdklog.Exporter, error) {
				return nil, errors.New("logger-fail")
			},
			ProviderFunc: func(ctx context.Context, exporters ...sdklog.Exporter) (*sdklog.LoggerProvider, func(context.Context) error, error) {
				return nil, func(ctx context.Context) error { return nil }, nil
			},
		},
		InitMetricFuncs: &InitMetricFuncs{
			ExporterFunc: func(ctx context.Context) (sdkmetric.Exporter, error) {
				metricCalled = true
				return nil, nil
			},
			ProviderFunc: func(ctx context.Context, exporter sdkmetric.Exporter) (*sdkmetric.MeterProvider, func(context.Context) error, error) {
				return nil, func(ctx context.Context) error { return nil }, nil
			},
		},
		InitTracerFuncs: &InitTracerFuncs{
			ExporterFunc: func(ctx context.Context) (sdktrace.SpanExporter, error) {
				tracerCalled = true
				return nil, nil
			},
			ProviderFunc: func(ctx context.Context, exporter sdktrace.SpanExporter) (*sdktrace.TracerProvider, func(context.Context) error, error) {
				return nil, nil, nil
			},
		},
	}

	err := initObservability(sm, initFuncs)
	assert.Error(err)
	assert.Contains(err.Error(), "logger-fail")
	assert.False(metricCalled)
	assert.False(tracerCalled)
}

func TestInitObservability_MetricErrorStopsSequence(t *testing.T) {
	assert := assert.New(t)
	sm := NewShutdownManager(context.Background())

	tracerCalled := false

	initFuncs := InitObservabilityFuncs{
		InitLoggerFuncs: &InitLoggerFuncs{
			ExporterFunc: func(ctx context.Context) ([]sdklog.Exporter, error) {
				return nil, nil
			},
			ProviderFunc: func(ctx context.Context, exporters ...sdklog.Exporter) (*sdklog.LoggerProvider, func(context.Context) error, error) {
				return nil, func(ctx context.Context) error { return nil }, nil
			},
		},
		InitMetricFuncs: &InitMetricFuncs{
			ExporterFunc: func(ctx context.Context) (sdkmetric.Exporter, error) {
				return nil, errors.New("metric-fail")
			},
			ProviderFunc: func(ctx context.Context, exporter sdkmetric.Exporter) (*sdkmetric.MeterProvider, func(context.Context) error, error) {
				return nil, nil, nil
			},
		},
		InitTracerFuncs: &InitTracerFuncs{
			ExporterFunc: func(ctx context.Context) (sdktrace.SpanExporter, error) {
				tracerCalled = true
				return nil, nil
			},
			ProviderFunc: func(ctx context.Context, exporter sdktrace.SpanExporter) (*sdktrace.TracerProvider, func(context.Context) error, error) {
				return nil, nil, nil
			},
		},
	}

	err := initObservability(sm, initFuncs)
	assert.Error(err)
	assert.Contains(err.Error(), "metric-fail")
	assert.False(tracerCalled)
}

func TestInitObservability_TracerErrorReported(t *testing.T) {
	assert := assert.New(t)
	sm := NewShutdownManager(context.Background())

	initFuncs := InitObservabilityFuncs{
		InitLoggerFuncs: &InitLoggerFuncs{
			ExporterFunc: func(ctx context.Context) ([]sdklog.Exporter, error) {
				return nil, nil
			},
			ProviderFunc: func(ctx context.Context, exporters ...sdklog.Exporter) (*sdklog.LoggerProvider, func(context.Context) error, error) {
				return nil, func(ctx context.Context) error { return nil }, nil
			},
		},
		InitMetricFuncs: &InitMetricFuncs{
			ExporterFunc: func(ctx context.Context) (sdkmetric.Exporter, error) {
				return nil, nil
			},
			ProviderFunc: func(ctx context.Context, exporter sdkmetric.Exporter) (*sdkmetric.MeterProvider, func(context.Context) error, error) {
				return nil, func(ctx context.Context) error { return nil }, nil
			},
		},
		InitTracerFuncs: &InitTracerFuncs{
			ExporterFunc: func(ctx context.Context) (sdktrace.SpanExporter, error) {
				return nil, errors.New("tracer-fail")
			},
			ProviderFunc: func(ctx context.Context, exporter sdktrace.SpanExporter) (*sdktrace.TracerProvider, func(context.Context) error, error) {
				return nil, nil, nil
			},
		},
	}

	err := initObservability(sm, initFuncs)
	assert.Error(err)
	assert.Contains(err.Error(), "tracer-fail")
}
