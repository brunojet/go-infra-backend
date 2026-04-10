package adapters

import (
	"context"
	"errors"
	"log"

	"github.com/brunojet/go-infra-backend/pkg/bootstrap/contracts"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

var otlpmetricgrpcNew = otlpmetricgrpc.New

type otlpMetric struct {
	provider *sdkmetric.MeterProvider
}

func logNewOTLPMetricResult(cfg OTLPConfig, metricsErr error) {
	log.Printf(
		"NewOTLPMetric: endpoint_set=%t insecure=%t metricsErr=%v",
		cfg.Endpoint != "",
		cfg.IsInsecure,
		metricsErr,
	)
}

func newMetricExporter(ctx context.Context, cfg OTLPConfig) (sdkmetric.Exporter, error) {
	options := []otlpmetricgrpc.Option{otlpmetricgrpc.WithEndpoint(cfg.Endpoint)}
	if cfg.IsInsecure {
		options = append(options, otlpmetricgrpc.WithInsecure())
	}
	return otlpmetricgrpcNew(ctx, options...)
}

func installMetricProvider(exporter sdkmetric.Exporter) (*otlpMetric, error) {
	if exporter == nil {
		return nil, errors.New("nil exporter provided")
	}
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)))
	otel.SetMeterProvider(provider)
	return &otlpMetric{provider: provider}, nil
}

func NewOTLPMetric(ctx context.Context, cfg OTLPConfig) (contracts.Shutdown, error) {
	var metricsErr error
	var metrics *otlpMetric
	exporter, metricsErr := newMetricExporter(ctx, cfg)
	if metricsErr == nil {
		metrics, metricsErr = installMetricProvider(exporter)
	}
	logNewOTLPMetricResult(cfg, metricsErr)
	return metrics, metricsErr
}

func NewOTLPMetricFromEnv(ctx context.Context) (contracts.Shutdown, error) {
	cfg := getMetricExporterConfigFromEnv()
	return NewOTLPMetric(ctx, cfg)
}

func (o *otlpMetric) Shutdown(ctx context.Context) error {
	o.provider.ForceFlush(ctx)
	return o.provider.Shutdown(ctx)
}
