package exporters

import (
	"context"
	"testing"
)

type fakeLogging struct {
	started      bool
	shutdown     bool
	lastSeverity string
	lastBody     string
}

func (f *fakeLogging) Start(ctx context.Context) error    { f.started = true; return nil }
func (f *fakeLogging) Shutdown(ctx context.Context) error { f.shutdown = true; return nil }
func (f *fakeLogging) Emit(ctx context.Context, severity string, body string) error {
	f.lastSeverity = severity
	f.lastBody = body
	return nil
}

func TestNewLoggingExporter_and_Forwarding(t *testing.T) {
	f := &fakeLogging{}
	l := NewLoggingExporter(f)
	if l == nil {
		t.Fatal("expected non-nil exporter")
	}
	if err := l.Start(context.Background()); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	if !f.started {
		t.Fatal("Start not forwarded")
	}
	if err := l.Emit(context.Background(), "info", "b"); err != nil {
		t.Fatalf("Emit failed: %v", err)
	}
	if f.lastSeverity != "info" || f.lastBody != "b" {
		t.Fatal("Emit not forwarded")
	}
	if err := l.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}
	if !f.shutdown {
		t.Fatal("Shutdown not forwarded")
	}
}
