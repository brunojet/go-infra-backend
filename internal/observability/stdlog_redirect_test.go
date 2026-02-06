package observability

import (
	"context"
	"log"
	"sync"
	"testing"

	otellog "go.opentelemetry.io/otel/log"
	otellogglobal "go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

type capturingLogExporter struct {
	mu      sync.Mutex
	records []sdklog.Record
}

func (e *capturingLogExporter) Export(ctx context.Context, records []sdklog.Record) error {
	_ = ctx
	e.mu.Lock()
	defer e.mu.Unlock()
	for i := range records {
		e.records = append(e.records, records[i].Clone())
	}
	return nil
}

func (e *capturingLogExporter) Shutdown(ctx context.Context) error {
	_ = ctx
	return nil
}

func (e *capturingLogExporter) ForceFlush(ctx context.Context) error {
	_ = ctx
	return nil
}

func (e *capturingLogExporter) snapshot() []sdklog.Record {
	e.mu.Lock()
	defer e.mu.Unlock()
	out := make([]sdklog.Record, len(e.records))
	copy(out, e.records)
	return out
}

func TestStdLogWriter_EmitsOnNewlineAndBuffersChunks(t *testing.T) {
	prevProvider := otellogglobal.GetLoggerProvider()
	defer otellogglobal.SetLoggerProvider(prevProvider)

	exp := &capturingLogExporter{}
	lp := sdklog.NewLoggerProvider(sdklog.WithProcessor(sdklog.NewSimpleProcessor(exp)))
	otellogglobal.SetLoggerProvider(lp)
	defer func() { _ = lp.Shutdown(context.Background()) }()

	w := NewStdLogWriter("test", otellog.SeverityInfo)

	if _, err := w.Write([]byte("hello")); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if got := len(exp.snapshot()); got != 0 {
		t.Fatalf("expected 0 records before newline, got %d", got)
	}

	if _, err := w.Write([]byte(" world\n")); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	recs := exp.snapshot()
	if len(recs) != 1 {
		t.Fatalf("expected 1 record, got %d", len(recs))
	}
	if got := recs[0].Body().AsString(); got != "hello world" {
		t.Fatalf("unexpected body: %q", got)
	}
	if got := recs[0].Severity(); got != otellog.SeverityInfo {
		t.Fatalf("unexpected severity: %v", got)
	}
}

func TestRedirectStdLog_RoutesStdlibLogToOTel(t *testing.T) {
	prevProvider := otellogglobal.GetLoggerProvider()
	defer otellogglobal.SetLoggerProvider(prevProvider)

	exp := &capturingLogExporter{}
	lp := sdklog.NewLoggerProvider(sdklog.WithProcessor(sdklog.NewSimpleProcessor(exp)))
	otellogglobal.SetLoggerProvider(lp)
	defer func() { _ = lp.Shutdown(context.Background()) }()

	restore := RedirectStdLog("stdlib-test", otellog.SeverityWarn)
	defer restore()

	log.Print("warn")

	recs := exp.snapshot()
	if len(recs) != 1 {
		t.Fatalf("expected 1 record, got %d", len(recs))
	}
	if got := recs[0].Body().AsString(); got != "warn" {
		t.Fatalf("unexpected body: %q", got)
	}
	if got := recs[0].Severity(); got != otellog.SeverityWarn {
		t.Fatalf("unexpected severity: %v", got)
	}
}
