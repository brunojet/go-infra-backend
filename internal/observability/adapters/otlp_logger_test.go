package adapters

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	otlploggrpc "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
)

func Test_NewOTLPLogger_ExporterError(t *testing.T) {
	origGrpc := otlploggrpcNew
	origStd := stdoutlogNew
	otlploggrpcNew = func(ctx context.Context, opts ...otlploggrpc.Option) (*otlploggrpc.Exporter, error) {
		return nil, errors.New("grpc exporter failed")
	}
	stdoutlogNew = func(opts ...stdoutlog.Option) (*stdoutlog.Exporter, error) {
		return nil, errors.New("stdout exporter failed")
	}
	defer func() { otlploggrpcNew = origGrpc; stdoutlogNew = origStd }()

	cfg := OTLPLoggerConfig{OTLPConfig: OTLPConfig{Endpoint: "localhost:4317", IsInsecure: true}}
	s, err := NewOTLPLogger(context.Background(), cfg)
	require.Error(t, err)
	require.Nil(t, s)
}
