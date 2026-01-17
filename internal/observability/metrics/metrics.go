package metrics

import (
	"context"
	"sync"

	contracts "github.com/brunojet/go-infra-backend/internal/observability/metrics/contracts"
)

// Metrics provides a small thread-safe facade for registering and using metrics
// instruments. It accepts a MetricsRegistrar (to obtain instrument handles)
// and an optional MetricsRecorder for direct name-based recording.
type Metrics struct {
	register contracts.MetricsRegister
	recorder contracts.MetricsRecorder

	mu         sync.RWMutex
	counters   map[string]contracts.Counter
	gauges     map[string]contracts.Gauge
	histograms map[string]contracts.Histogram
	closed     bool
}

// NewMetrics creates a metrics core. `recorder` may be nil; in that case Metrics
// will fall back to instrument handles returned by the registrar.
func NewMetrics(registrar contracts.MetricsRegister, recorder contracts.MetricsRecorder) *Metrics {
	return &Metrics{
		register:   registrar,
		recorder:   recorder,
		counters:   make(map[string]contracts.Counter),
		gauges:     make(map[string]contracts.Gauge),
		histograms: make(map[string]contracts.Histogram),
	}
}

func (c *Metrics) RegisterCounter(name string, labelKeys []string) contracts.Counter {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil
	}
	if inst, ok := c.counters[name]; ok {
		return inst
	}
	inst := c.register.RegisterCounter(name, labelKeys)
	c.counters[name] = inst
	return inst
}

func (c *Metrics) GetCounter(name string) (contracts.Counter, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	inst, ok := c.counters[name]
	return inst, ok
}

// CounterAdd records a counter. If a MetricsRecorder was provided it is used
// (and labels are forwarded). Otherwise a registered Counter handle is used
// and labels are ignored.
func (c *Metrics) CounterAdd(ctx context.Context, name string, value float64, labels contracts.Labels) {
	if c.recorder != nil {
		c.recorder.CounterAdd(ctx, name, value, labels)
		return
	}
	// fallback to handle
	c.mu.RLock()
	inst, ok := c.counters[name]
	c.mu.RUnlock()
	if !ok {
		inst = c.RegisterCounter(name, nil)
		if inst == nil {
			return
		}
	}
	inst.Add(ctx, value)
}

// RegisterGauge and GaugeSet
func (c *Metrics) RegisterGauge(name string, labelKeys []string) contracts.Gauge {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil
	}
	if inst, ok := c.gauges[name]; ok {
		return inst
	}
	inst := c.register.RegisterGauge(name, labelKeys)
	c.gauges[name] = inst
	return inst
}

func (c *Metrics) GaugeSet(ctx context.Context, name string, value float64, labels contracts.Labels) {
	if c.recorder != nil {
		c.recorder.GaugeSet(ctx, name, value, labels)
		return
	}
	c.mu.RLock()
	inst, ok := c.gauges[name]
	c.mu.RUnlock()
	if !ok {
		inst = c.RegisterGauge(name, nil)
		if inst == nil {
			return
		}
	}
	inst.Set(ctx, value)
}

// RegisterHistogram and HistogramObserve
func (c *Metrics) RegisterHistogram(name string, labelKeys []string) contracts.Histogram {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil
	}
	if inst, ok := c.histograms[name]; ok {
		return inst
	}
	inst := c.register.RegisterHistogram(name, labelKeys)
	c.histograms[name] = inst
	return inst
}

func (c *Metrics) HistogramObserve(ctx context.Context, name string, value float64, labels contracts.Labels) {
	if c.recorder != nil {
		c.recorder.HistogramObserve(ctx, name, value, labels)
		return
	}
	c.mu.RLock()
	inst, ok := c.histograms[name]
	c.mu.RUnlock()
	if !ok {
		inst = c.RegisterHistogram(name, nil)
		if inst == nil {
			return
		}
	}
	inst.Observe(ctx, value)
}

// Close shuts down the core; if a recorder exists its Close is called.
func (c *Metrics) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil
	}
	c.closed = true
	if c.recorder != nil {
		return c.recorder.Close()
	}
	return nil
}
