package exporters

import (
	"context"
	"errors"

	"github.com/brunojet/go-infra-backend/pkg/bootstrap/contracts"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

var otlpmetricgrpcNew = otlpmetricgrpc.New

type otlpMetric struct {
	provider *sdkmetric.MeterProvider
}

func newMetricExporter(ctx context.Context, cfg OTLPConfig) (sdkmetric.Exporter, error) {
	options := []otlpmetricgrpc.Option{otlpmetricgrpc.WithEndpoint(cfg.Endpoint)}
	if cfg.IsInsecure {
		options = append(options, otlpmetricgrpc.WithInsecure())
	}
	return otlpmetricgrpcNew(ctx, options...)
}

func logNewOTLPMetricsResult(cfg OTLPConfig, metricsErr error) {
	logPrintf(
		"NewOTLPMetrics: endpoint_set=%t insecure=%t metricsErr=%v",
		cfg.Endpoint != "",
		cfg.IsInsecure,
		metricsErr,
	)
}

func installMetricProvider(exporter sdkmetric.Exporter) (*otlpMetric, error) {
	if exporter == nil {
		return nil, errors.New("nil exporter provided")
	}
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)))
	otel.SetMeterProvider(provider)
	return &otlpMetric{provider: provider}, nil
}

func NewOTLPMetrics(ctx context.Context, cfg OTLPConfig) (contracts.Shutdown, error) {
	var metricsErr error
	var metrics *otlpMetric
	exporter, metricsErr := newMetricExporter(ctx, cfg)
	if metricsErr == nil {
		metrics, metricsErr = installMetricProvider(exporter)
	}
	logNewOTLPMetricsResult(cfg, metricsErr)
	return metrics, metricsErr
}

func (o *otlpMetric) Shutdown(ctx context.Context) error {
	return o.provider.Shutdown(ctx)
}

func NewOTLPMetricsFromEnv(ctx context.Context) (contracts.Shutdown, error) {
	cfg := getMetricExporterConfigFromEnv()
	return NewOTLPMetrics(ctx, cfg)
}
