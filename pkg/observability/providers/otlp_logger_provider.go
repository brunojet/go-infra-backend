package providers

import (
	"context"

	internalproviders "github.com/brunojet/go-infra-backend/internal/observability/providers"
	"github.com/brunojet/go-infra-backend/pkg/observability/contracts"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

// NewOTLPLoggerProvider builds and registers a LoggerProvider using the provided
// exporters. The exporter must implement sdklog.Exporter (allows passing noop or
// real exporters in tests).
func NewOTLPLoggerProvider(ctx context.Context, exporters ...sdklog.Exporter) (*sdklog.LoggerProvider, contracts.ShutdownFunc, error) {
	lp, shutdown, err := internalproviders.NewOTLPLoggerProvider(ctx, exporters...)
	return lp, shutdown, err
}
