package middlewares

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	otellog "go.opentelemetry.io/otel/log"
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

	require.Equal(t, http.StatusBadRequest, w.Code)

	recs := logExp.snapshot()
	require.NotEmpty(t, recs, "expected at least one log record")

	// The SDK can link logs to trace/span via ctx. We assert TraceID is set.
	require.True(t, recs[0].TraceID().IsValid(), "expected log record to have valid TraceID")
}

func TestSeverityForHTTPStatus_AllCases(t *testing.T) {
	require.Equal(t, otellog.SeverityInfo, severityForHTTPStatus(200))
	require.Equal(t, otellog.SeverityWarn, severityForHTTPStatus(400))
	require.Equal(t, otellog.SeverityError, severityForHTTPStatus(500))
}

func TestOTLPErrorLogMiddleware_EmitsWhenGinErrorsEvenIfStatusBelowMin(t *testing.T) {
	prevProvider := otellogglobal.GetLoggerProvider()
	defer otellogglobal.SetLoggerProvider(prevProvider)

	logExp := &memoryLogExporter{}
	lp := sdklog.NewLoggerProvider(sdklog.WithProcessor(sdklog.NewSimpleProcessor(logExp)))
	otellogglobal.SetLoggerProvider(lp)
	defer lp.Shutdown(context.Background())

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(OtelGinMiddleware())
	r.Use(OTLPErrorLogMiddleware()) // default minStatus = 500

	r.GET("/ok-with-error", func(c *gin.Context) {
		c.Error(errors.New("boom-error"))
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest("GET", "/ok-with-error", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	recs := logExp.snapshot()
	require.NotEmpty(t, recs, "expected at least one log record")

	// body should include the gin error text
	require.Contains(t, recs[0].Body().AsString(), "boom-error")

	// severity for status 200 should be Info
	require.Equal(t, otellog.SeverityInfo, recs[0].Severity())
}

func TestOTLPErrorLogMiddleware_DefaultMinStatus_EmitsOn500(t *testing.T) {
	prevProvider := otellogglobal.GetLoggerProvider()
	defer otellogglobal.SetLoggerProvider(prevProvider)

	logExp := &memoryLogExporter{}
	lp := sdklog.NewLoggerProvider(sdklog.WithProcessor(sdklog.NewSimpleProcessor(logExp)))
	otellogglobal.SetLoggerProvider(lp)
	defer lp.Shutdown(context.Background())

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(OtelGinMiddleware())
	r.Use(OTLPErrorLogMiddleware()) // default minStatus = 500

	r.GET("/internal", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusInternalServerError)
	})

	req := httptest.NewRequest("GET", "/internal", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	recs := logExp.snapshot()
	require.NotEmpty(t, recs, "expected at least one log record")

	require.Equal(t, otellog.SeverityError, recs[0].Severity())
}

func TestOTLPErrorLogMiddleware_ReturnsWhenBelowMinStatusAndNoErrors(t *testing.T) {
	prevProvider := otellogglobal.GetLoggerProvider()
	defer otellogglobal.SetLoggerProvider(prevProvider)

	logExp := &memoryLogExporter{}
	lp := sdklog.NewLoggerProvider(sdklog.WithProcessor(sdklog.NewSimpleProcessor(logExp)))
	otellogglobal.SetLoggerProvider(lp)
	defer lp.Shutdown(context.Background())

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(OtelGinMiddleware())
	r.Use(OTLPErrorLogMiddleware()) // default minStatus = 500

	r.GET("/ok", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest("GET", "/ok", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	recs := logExp.snapshot()
	require.Empty(t, recs, "expected 0 log records")
}
