package providers

import (
	"context"
	"log"

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
			if err := ff.ForceFlush(ctx); err != nil {
				log.Printf("otel metric force flush error: %v", err)
			}
		}

		if err := mp.Shutdown(ctx); err != nil {
			log.Printf("otel metric shutdown error: %v", err)
			return err
		}
		return nil
	}

	return mp, shutdown, nil
}
