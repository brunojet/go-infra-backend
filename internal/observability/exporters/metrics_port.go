package exporters

import (
	"context"

	"github.com/brunojet/go-infra-backend/internal/observability/exporters/contracts"
)

// Metrics is a thin wrapper around a contracts.MetricsExporter implementation.
type Metrics struct {
	impl contracts.MetricsExporter
}

// NewMetricsExporter constructs a Metrics exporter from a given adapter implementation.
// It does not create or wire any concrete adapter itself, keeping the port decoupled.
func NewMetricsExporter(a contracts.MetricsExporter) *Metrics {
	if a == nil {
		return nil
	}
	return &Metrics{impl: a}
}

func (m *Metrics) Start(ctx context.Context) error    { return m.impl.Start(ctx) }
func (m *Metrics) Shutdown(ctx context.Context) error { return m.impl.Shutdown(ctx) }
func (m *Metrics) Record(ctx context.Context, name string, value float64) error {
	return m.impl.Record(ctx, name, value)
}
