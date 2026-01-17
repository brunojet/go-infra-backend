package contracts_test

import (
	"context"
	"testing"

	mm "github.com/brunojet/go-infra-backend/internal/observability/metrics/contracts"
)

type mockCounter struct{}

func (m *mockCounter) Add(ctx context.Context, v float64) {}

type mockGauge struct{}

func (m *mockGauge) Set(ctx context.Context, v float64) {}

type mockHistogram struct{}

func (m *mockHistogram) Observe(ctx context.Context, v float64) {}

type mockRecorder struct{}

func (m *mockRecorder) CounterAdd(ctx context.Context, name string, value float64, labels mm.Labels) {
}
func (m *mockRecorder) GaugeSet(ctx context.Context, name string, value float64, labels mm.Labels) {}
func (m *mockRecorder) HistogramObserve(ctx context.Context, name string, value float64, labels mm.Labels) {
}
func (m *mockRecorder) Close() error { return nil }
func (m *mockRecorder) RegisterCounter(name string, labelKeys []string) mm.Counter {
	return &mockCounter{}
}
func (m *mockRecorder) RegisterGauge(name string, labelKeys []string) mm.Gauge { return &mockGauge{} }
func (m *mockRecorder) RegisterHistogram(name string, labelKeys []string) mm.Histogram {
	return &mockHistogram{}
}

func TestMetricsRecorderInterface(t *testing.T) {
	var _ mm.MetricsRecorder = (*mockRecorder)(nil)
	var _ mm.MetricsRegister = (*mockRecorder)(nil)
}
