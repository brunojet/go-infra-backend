package contracts

import "context"

// Labels represents a set of key/value label pairs for metrics.
type Labels map[string]string

// MetricsRecorder defines a small, implementation-agnostic API for recording
// metrics from the core application. Adapters (Prometheus/OpenTelemetry/etc.)
// should implement this interface.
type MetricsRecorder interface {
	// CounterAdd increments (or decrements if negative) a named counter.
	CounterAdd(ctx context.Context, name string, value float64, labels Labels)

	// GaugeSet sets a gauge to a specific value.
	GaugeSet(ctx context.Context, name string, value float64, labels Labels)

	// HistogramObserve records an observation for a histogram (or summary).
	HistogramObserve(ctx context.Context, name string, value float64, labels Labels)

	// Close releases resources used by the recorder (optional for adapters).
	Close() error
}

// Counter is a handle to a counter instrument.
type Counter interface {
	Add(ctx context.Context, value float64)
}

// Gauge is a handle to a gauge instrument.
type Gauge interface {
	Set(ctx context.Context, value float64)
}

// Histogram is a handle to a histogram (or summary) instrument.
type Histogram interface {
	Observe(ctx context.Context, value float64)
}

// Register* methods allow adapters to return instrument handles which are
// efficient to reuse in hot paths instead of using string lookups every time.
type MetricsRegister interface {
	RegisterCounter(name string, labelKeys []string) Counter
	RegisterGauge(name string, labelKeys []string) Gauge
	RegisterHistogram(name string, labelKeys []string) Histogram
}
