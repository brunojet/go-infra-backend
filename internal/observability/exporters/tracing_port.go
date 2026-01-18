package exporters

import (
	"context"

	"github.com/brunojet/go-infra-backend/internal/observability/exporters/contracts"
)

// Tracing is a thin wrapper around a contracts.TracingExporter implementation.
type Tracing struct{ impl contracts.TracingExporter }

// NewTracingExporter constructs a Tracing exporter from a given adapter implementation.
func NewTracingExporter(a contracts.TracingExporter) *Tracing {
	if a == nil {
		return nil
	}
	return &Tracing{impl: a}
}

func (t *Tracing) Start(ctx context.Context) error    { return t.impl.Start(ctx) }
func (t *Tracing) Shutdown(ctx context.Context) error { return t.impl.Shutdown(ctx) }
func (t *Tracing) StartSpan(ctx context.Context, name string) (context.Context, func()) {
	return t.impl.StartSpan(ctx, name)
}
