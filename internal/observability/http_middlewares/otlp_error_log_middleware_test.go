package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	otellogglobal "go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

type memoryLogExporter struct {
	mu      sync.Mutex
	records []sdklog.Record
}

func (m *memoryLogExporter) Export(ctx context.Context, records []sdklog.Record) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, r := range records {
		m.records = append(m.records, r.Clone())
	}
	return nil
}

func (m *memoryLogExporter) Shutdown(ctx context.Context) error { return nil }

func (m *memoryLogExporter) ForceFlush(ctx context.Context) error { return nil }

func (m *memoryLogExporter) snapshot() []sdklog.Record {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]sdklog.Record, len(m.records))
	copy(out, m.records)
	return out
}

func TestOTLPErrorLogMiddleware_EmitsLogWithTraceContext(t *testing.T) {
	// tracer provider: in-memory exporter to ensure spans exist
	exp := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sdktrace.NewSimpleSpanProcessor(exp)))
	otel.SetTracerProvider(tp)
	defer tp.Shutdown(context.Background())

	// logger provider: memory exporter
	logExp := &memoryLogExporter{}
	lp := sdklog.NewLoggerProvider(sdklog.WithProcessor(sdklog.NewSimpleProcessor(logExp)))
	otellogglobal.SetLoggerProvider(lp)
	defer lp.Shutdown(context.Background())

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(OtelGinMiddleware())
	r.Use(OTLPErrorLogMiddleware(WithMinStatus(400)))

	r.GET("/bad", func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bad request"})
	})

	req := httptest.NewRequest("GET", "/bad", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: %d", w.Code)
	}

	recs := logExp.snapshot()
	if len(recs) == 0 {
		t.Fatalf("expected at least one log record, got 0")
	}

	// The SDK can link logs to trace/span via ctx. We assert TraceID is set.
	if !recs[0].TraceID().IsValid() {
		t.Fatalf("expected log record to have valid TraceID")
	}
}
