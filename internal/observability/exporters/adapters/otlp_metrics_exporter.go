package adapters

import (
	"context"

	"go.opentelemetry.io/otel"
	otlpmetricgrpc "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"google.golang.org/grpc"

	"github.com/brunojet/go-infra-backend/internal/observability/exporters/contracts"
)

// OTLPMetricsExporter implements contracts.MetricsExporter using OpenTelemetry OTLP.
type OTLPMetricsExporter struct {
	mp   *metric.MeterProvider
	conn *grpc.ClientConn
}

// NewOTLPMetricsExporter returns a new metrics exporter (endpoint currently fixed to localhost:4317).
func NewOTLPMetricsExporter() contracts.MetricsExporter {
	return &OTLPMetricsExporter{mp: nil}
}

// Start initializes the metrics exporter.
func (o *OTLPMetricsExporter) Start(ctx context.Context) error {
	endpoint := EndpointFromEnv()
	res, _ := resource.Merge(resource.Default(), resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String("go-infra-backend")))

	// Try to dial OTLP collector and create exporter using existing gRPC conn.
	cc, err := DialOTLP(ctx, endpoint)
	if err == nil && cc != nil {
		exporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(cc))
		if err == nil {
			reader := metric.NewPeriodicReader(exporter)
			mp := metric.NewMeterProvider(metric.WithReader(reader), metric.WithResource(res))
			o.mp = mp
			o.conn = cc
			otel.SetMeterProvider(mp)
			return nil
		}
		_ = cc.Close()
	}

	mp := metric.NewMeterProvider(metric.WithResource(res))
	o.mp = mp
	otel.SetMeterProvider(mp)
	return nil
}

// Shutdown shuts down the meter provider.
func (o *OTLPMetricsExporter) Shutdown(ctx context.Context) error {
	if o.mp == nil {
		return nil
	}
	if err := o.mp.Shutdown(ctx); err != nil {
		return err
	}
	if o.conn != nil {
		return o.conn.Close()
	}
	return nil
}

// Record records a measurement using a simple counter instrument.
func (o *OTLPMetricsExporter) Record(ctx context.Context, name string, value float64) error {
	meter := otel.Meter("otlp-exporter")
	// Use float64 counter
	counter, err := meter.Float64Counter(name)
	if err != nil {
		return err
	}
	// record value without additional options to keep cross-version compatibility
	counter.Add(ctx, value)
	return nil
}

var _ contracts.MetricsExporter = (*OTLPMetricsExporter)(nil)
