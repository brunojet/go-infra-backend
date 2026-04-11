package adapters

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	otlpmetricgrpc "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
)

func Test_NewOTLPMetrics_ExporterError(t *testing.T) {
	orig := otlpmetricgrpcNew
	otlpmetricgrpcNew = func(ctx context.Context, opts ...otlpmetricgrpc.Option) (*otlpmetricgrpc.Exporter, error) {
		return nil, errors.New("exporter failed")
	}
	defer func() { otlpmetricgrpcNew = orig }()

	cfg := OTLPConfig{Endpoint: "localhost:4317", IsInsecure: true}
	s, err := NewOTLPMetric(context.Background(), cfg)
	require.Error(t, err)
	require.Nil(t, s)
}
