package exporters

import (
	"context"
	"testing"
)

type fakeTracing struct {
	started      bool
	shutdown     bool
	lastSpanName string
}

func (f *fakeTracing) Start(ctx context.Context) error    { f.started = true; return nil }
func (f *fakeTracing) Shutdown(ctx context.Context) error { f.shutdown = true; return nil }
func (f *fakeTracing) StartSpan(ctx context.Context, name string) (context.Context, func()) {
	f.lastSpanName = name
	return ctx, func() {}
}

func TestNewTracingExporter_and_Forwarding(t *testing.T) {
	f := &fakeTracing{}
	tr := NewTracingExporter(f)
	if tr == nil {
		t.Fatal("expected non-nil exporter")
	}
	if err := tr.Start(context.Background()); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	if !f.started {
		t.Fatal("Start not forwarded")
	}
	_, end := tr.StartSpan(context.Background(), "s")
	end()
	if f.lastSpanName != "s" {
		t.Fatal("StartSpan not forwarded")
	}
	if err := tr.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}
	if !f.shutdown {
		t.Fatal("Shutdown not forwarded")
	}
}
