package exporters

import (
	"context"
	"testing"
)

type fakeMetrics struct {
	started   bool
	shutdown  bool
	lastName  string
	lastValue float64
}

func (f *fakeMetrics) Start(ctx context.Context) error    { f.started = true; return nil }
func (f *fakeMetrics) Shutdown(ctx context.Context) error { f.shutdown = true; return nil }
func (f *fakeMetrics) Record(ctx context.Context, name string, value float64) error {
	f.lastName = name
	f.lastValue = value
	return nil
}

func TestNewMetricsExporter_and_Forwarding(t *testing.T) {
	f := &fakeMetrics{}
	m := NewMetricsExporter(f)
	if m == nil {
		t.Fatal("expected non-nil exporter")
	}
	if err := m.Start(context.Background()); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	if !f.started {
		t.Fatal("Start not forwarded")
	}
	if err := m.Record(context.Background(), "m", 1.5); err != nil {
		t.Fatalf("Record failed: %v", err)
	}
	if f.lastName != "m" || f.lastValue != 1.5 {
		t.Fatal("Record not forwarded")
	}
	if err := m.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}
	if !f.shutdown {
		t.Fatal("Shutdown not forwarded")
	}
}
