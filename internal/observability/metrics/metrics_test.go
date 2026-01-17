package metrics_test

import (
	"context"
	"sync"
	"testing"

	"github.com/brunojet/go-infra-backend/internal/observability/metrics"
	"github.com/brunojet/go-infra-backend/internal/observability/metrics/contracts"
)

type mockCounter struct{}

func (m *mockCounter) Add(ctx context.Context, v float64) {}

type mockGauge struct{}

func (m *mockGauge) Set(ctx context.Context, v float64) {}

type mockHistogram struct{}

func (m *mockHistogram) Observe(ctx context.Context, v float64) {}

type mockRegistrar struct {
	mu       sync.Mutex
	counters map[string]int
}

func (r *mockRegistrar) RegisterCounter(name string, labelKeys []string) contracts.Counter {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.counters == nil {
		r.counters = make(map[string]int)
	}
	r.counters[name]++
	return &mockCounter{}
}
func (r *mockRegistrar) RegisterGauge(name string, labelKeys []string) contracts.Gauge {
	return &mockGauge{}
}
func (r *mockRegistrar) RegisterHistogram(name string, labelKeys []string) contracts.Histogram {
	return &mockHistogram{}
}

type mockRecorder struct{}

func (m *mockRecorder) CounterAdd(ctx context.Context, name string, value float64, labels contracts.Labels) {
}
func (m *mockRecorder) GaugeSet(ctx context.Context, name string, value float64, labels contracts.Labels) {
}
func (m *mockRecorder) HistogramObserve(ctx context.Context, name string, value float64, labels contracts.Labels) {
}
func (m *mockRecorder) Close() error { return nil }

func TestCore_RegisterAndUseInstruments(t *testing.T) {
	reg := &mockRegistrar{}
	rec := &mockRecorder{}
	c := metrics.NewMetrics(reg, rec)

	// using recorder path should not panic
	c.CounterAdd(context.Background(), "requests_total", 1, contracts.Labels{"path": "/"})
	c.GaugeSet(context.Background(), "inflight", 3, nil)
	c.HistogramObserve(context.Background(), "latency_ms", 123.4, nil)

	// ensure registration via registrar (via RegisterCounter)
	_ = c.RegisterCounter("my_counter", []string{"k"})
	if _, ok := c.GetCounter("my_counter"); !ok {
		t.Fatal("expected my_counter to be registered")
	}
}
