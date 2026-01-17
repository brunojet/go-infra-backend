package adapter

import (
	"context"
	"testing"
)

func TestMemMetricsRecorder_Basic(t *testing.T) {
	rec := NewMemRecorder()
	reg := NewMemRegister(rec)

	cnt := reg.RegisterCounter("requests_total", nil)
	cnt.Add(context.Background(), 1)
	cnt.Add(context.Background(), 2)

	g := reg.RegisterGauge("inflight", nil)
	g.Set(context.Background(), 3)

	h := reg.RegisterHistogram("latency_ms", nil)
	h.Observe(context.Background(), 10.5)
	h.Observe(context.Background(), 5.5)

	if rec.counters["requests_total"] != 3 {
		t.Fatalf("expected requests_total=3 got %v", rec.counters["requests_total"])
	}
	if rec.gauges["inflight"] != 3 {
		t.Fatalf("expected inflight=3 got %v", rec.gauges["inflight"])
	}
	if len(rec.histograms["latency_ms"]) != 2 {
		t.Fatalf("expected 2 histogram observations got %d", len(rec.histograms["latency_ms"]))
	}
}
