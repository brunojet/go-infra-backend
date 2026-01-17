package adapters_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"

	ada "github.com/brunojet/go-infra-backend/internal/observability/middlewares/http/adapters"
)

func TestOtelGinMiddleware_CreatesSpan(t *testing.T) {
	// arrange: in-memory tracer
	sr := tracetest.NewSpanRecorder()
	tp := trace.NewTracerProvider(trace.WithSpanProcessor(sr))
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	r := gin.New()
	r.Use(ada.OtelGinMiddleware("test-service"))
	r.GET("/ok", func(c *gin.Context) { c.Status(204) })

	// act
	req := httptest.NewRequest("GET", "/ok", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	// assert
	if rec.Code != 204 {
		t.Fatalf("expected 204 got %d", rec.Code)
	}
	if len(sr.Ended()) == 0 {
		t.Fatalf("expected otel span recorded")
	}
}
