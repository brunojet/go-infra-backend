package adapters

import (
	"context"
	"time"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/brunojet/go-infra-backend/internal/observability/stats"
)

// OTLPTraceExporterWrapper wraps an existing sdktrace.SpanExporter and embeds
// a TracerWrapper so it can observe export attempts. The wrapper exposes a
// generic OnqueueExport API so the same pattern can be used for spans, metrics
// and logs.
type OTLPTraceExporterWrapper struct {
	underlying sdktrace.SpanExporter
	*stats.ObservabilityStats
}

// NewOTLPTraceExporterWrapper wraps the provided underlying exporter and
// creates an embedded TracerWrapper using the provided metadata.
func NewOTLPTraceExporterWrapper(underlying sdktrace.SpanExporter, id, implementation, endpoint string, isInsecure bool, configVersion string, maxErrors int) *OTLPTraceExporterWrapper {
	tw := stats.NewObservabilityStats(id, implementation, endpoint, isInsecure, configVersion, maxErrors)
	return &OTLPTraceExporterWrapper{
		underlying:         underlying,
		ObservabilityStats: tw,
	}
}

// ExportSpans implements sdktrace.SpanExporter by delegating to the
// underlying exporter and then reporting the attempt via OnqueueExport.
func (w *OTLPTraceExporterWrapper) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	// mark items as queued/onqueue and obtain completion callback
	var done func(err error)
	if w.ObservabilityStats != nil {
		done = w.ObservabilityStats.OnqueueExport(len(spans))
	}

	start := time.Now()
	err := w.underlying.ExportSpans(ctx, spans)
	_ = time.Since(start)

	if done != nil {
		done(err)
	}

	return err
}

// Shutdown forwards the call to the embedded tracer wrapper and the underlying exporter.
func (w *OTLPTraceExporterWrapper) Shutdown(ctx context.Context) error {
	if w.ObservabilityStats != nil {
		_ = w.ObservabilityStats.Shutdown(ctx)
	}
	return w.underlying.Shutdown(ctx)
}
