package adapter

import (
	"context"
	"sync"

	contracts "github.com/brunojet/go-infra-backend/internal/observability/metrics/contracts"
)

// Mem-based recorder that mimics behavior of a metrics adapter. This is a
// lightweight, test-friendly implementation used as a minimal OpenTelemetry
// adapter substitute for unit tests.

type MemRecorder struct {
	mu         sync.Mutex
	counters   map[string]float64
	gauges     map[string]float64
	histograms map[string][]float64
}

func NewMemRecorder() *MemRecorder {
	return &MemRecorder{counters: map[string]float64{}, gauges: map[string]float64{}, histograms: map[string][]float64{}}
}

func (m *MemRecorder) CounterAdd(ctx context.Context, name string, value float64, labels contracts.Labels) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters[name] += value
}
func (m *MemRecorder) GaugeSet(ctx context.Context, name string, value float64, labels contracts.Labels) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.gauges[name] = value
}
func (m *MemRecorder) HistogramObserve(ctx context.Context, name string, value float64, labels contracts.Labels) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.histograms[name] = append(m.histograms[name], value)
}
func (m *MemRecorder) Close() error { return nil }

// Accessors return copies of internal maps for inspection in tests.
func (m *MemRecorder) Counters() map[string]float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make(map[string]float64, len(m.counters))
	for k, v := range m.counters {
		out[k] = v
	}
	return out
}

func (m *MemRecorder) Gauges() map[string]float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make(map[string]float64, len(m.gauges))
	for k, v := range m.gauges {
		out[k] = v
	}
	return out
}

func (m *MemRecorder) Histograms() map[string][]float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make(map[string][]float64, len(m.histograms))
	for k, v := range m.histograms {
		out[k] = append([]float64(nil), v...)
	}
	return out
}

// Registrar returns lightweight handles that update the MemRecorder.

type MemRegister struct{ rec *MemRecorder }

func NewMemRegister(rec *MemRecorder) *MemRegister { return &MemRegister{rec: rec} }

type memCounterHandle struct {
	name string
	rec  *MemRecorder
}

func (h *memCounterHandle) Add(ctx context.Context, v float64) { h.rec.CounterAdd(ctx, h.name, v, nil) }

type memGaugeHandle struct {
	name string
	rec  *MemRecorder
}

func (h *memGaugeHandle) Set(ctx context.Context, v float64) { h.rec.GaugeSet(ctx, h.name, v, nil) }

type memHistogramHandle struct {
	name string
	rec  *MemRecorder
}

func (h *memHistogramHandle) Observe(ctx context.Context, v float64) {
	h.rec.HistogramObserve(ctx, h.name, v, nil)
}

func (m *MemRegister) RegisterCounter(name string, labelKeys []string) contracts.Counter {
	return &memCounterHandle{name: name, rec: m.rec}
}
func (m *MemRegister) RegisterGauge(name string, labelKeys []string) contracts.Gauge {
	return &memGaugeHandle{name: name, rec: m.rec}
}
func (m *MemRegister) RegisterHistogram(name string, labelKeys []string) contracts.Histogram {
	return &memHistogramHandle{name: name, rec: m.rec}
}
