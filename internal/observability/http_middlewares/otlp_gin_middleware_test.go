package middlewares

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestTracingMiddleware_CreatesSpan(t *testing.T) {
	// in-memory exporter to capture spans
	exp := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sdktrace.NewSimpleSpanProcessor(exp)))
	otel.SetTracerProvider(tp)
	defer tp.Shutdown(context.Background())

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(OtelGinMiddleware())
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	req := httptest.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("unexpected status: %d", w.Code)
	}

	spans := exp.GetSpans()
	if len(spans) == 0 {
		t.Fatalf("expected at least one span, got 0")
	}
}
