package providers

import (
	"context"

	"go.opentelemetry.io/otel"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// NewOTLPMetricProvider builds and registers a MeterProvider using the OTLP
// metric exporter and the resource from NewOTLPResource. Returns the provider
// and a shutdown function.
func NewOTLPMetricProvider(ctx context.Context, exporter sdkmetric.Exporter) (*sdkmetric.MeterProvider, func(context.Context) error, error) {
	reader := sdkmetric.NewPeriodicReader(exporter)

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
	)

	otel.SetMeterProvider(mp)

	shutdown := func(ctx context.Context) error {
		if ff, ok := interface{}(mp).(interface{ ForceFlush(context.Context) error }); ok {
			ff.ForceFlush(ctx)
		}
		mp.Shutdown(ctx)
		return nil
	}

	return mp, shutdown, nil
}
