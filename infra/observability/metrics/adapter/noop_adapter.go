package adapter

import (
	"context"

	contracts "github.com/brunojet/go-infra-backend/infra/observability/metrics/contracts"
)

// Noop implementations for metrics instruments and recorder.

type noopCounter struct{}

func (noopCounter) Add(ctx context.Context, v float64) {}

type noopGauge struct{}

func (noopGauge) Set(ctx context.Context, v float64) {}

type noopHistogram struct{}

func (noopHistogram) Observe(ctx context.Context, v float64) {}

// NoopRecorder does nothing; useful for testing or when metrics are disabled.
type NoopRecorder struct{}

func NewNoopRecorder() *NoopRecorder { return &NoopRecorder{} }

func (NoopRecorder) CounterAdd(ctx context.Context, name string, value float64, labels contracts.Labels) {
}
func (NoopRecorder) GaugeSet(ctx context.Context, name string, value float64, labels contracts.Labels) {
}
func (NoopRecorder) HistogramObserve(ctx context.Context, name string, value float64, labels contracts.Labels) {
}
func (NoopRecorder) Close() error { return nil }

// NoopRegistrar provides instrument handles that are cheap to reuse.
type NoopRegistrar struct{}

func NewNoopRegistrar() *NoopRegistrar { return &NoopRegistrar{} }

func (NoopRegistrar) RegisterCounter(name string, labelKeys []string) contracts.Counter {
	return &noopCounter{}
}
func (NoopRegistrar) RegisterGauge(name string, labelKeys []string) contracts.Gauge {
	return &noopGauge{}
}
func (NoopRegistrar) RegisterHistogram(name string, labelKeys []string) contracts.Histogram {
	return &noopHistogram{}
}
