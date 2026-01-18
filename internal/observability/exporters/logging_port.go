package exporters

import (
	"context"

	"github.com/brunojet/go-infra-backend/internal/observability/exporters/contracts"
)

// Logging is a thin wrapper around a contracts.LoggingExporter implementation.
type Logging struct{ impl contracts.LoggingExporter }

// NewLoggingExporter constructs a Logging exporter from a given adapter implementation.
func NewLoggingExporter(a contracts.LoggingExporter) *Logging {
	if a == nil {
		return nil
	}
	return &Logging{impl: a}
}

func (l *Logging) Start(ctx context.Context) error    { return l.impl.Start(ctx) }
func (l *Logging) Shutdown(ctx context.Context) error { return l.impl.Shutdown(ctx) }
func (l *Logging) Emit(ctx context.Context, severity, body string) error {
	return l.impl.Emit(ctx, severity, body)
}
