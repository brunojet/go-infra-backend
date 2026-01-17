package httpmiddleware_test

import (
	"net/http/httptest"
	"testing"

	httpmw "github.com/brunojet/go-infra-backend/infra/observability/middlewares/http"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestEnrichAppSpan_SetsAttributes(t *testing.T) {
	// in-memory tracer
	sr := tracetest.NewSpanRecorder()
	tp := trace.NewTracerProvider(trace.WithSpanProcessor(sr))
	otel.SetTracerProvider(tp)

	r := gin.New()
	// pre-middleware: create span and attach to request context (simulates otelgin)
	r.Use(func(c *gin.Context) {
		ctx, sp := otel.Tracer("test").Start(c.Request.Context(), "req")
		c.Request = c.Request.WithContext(ctx)
		c.Next()
		sp.End()
	})
	r.Use(httpmw.EnrichAppSpan())

	r.GET("/test", func(c *gin.Context) {
		c.Status(204)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-User-ID", "u123")
	req.Header.Set("X-Tenant-ID", "t456")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	ended := sr.Ended()
	if len(ended) == 0 {
		t.Fatalf("expected ended spans, got 0")
	}
	attrs := ended[0].Attributes()
	foundUser := false
	foundTenant := false
	for _, a := range attrs {
		if string(a.Key) == "app.user_id" && a.Value.AsString() == "u123" {
			foundUser = true
		}
		if string(a.Key) == "app.tenant_id" && a.Value.AsString() == "t456" {
			foundTenant = true
		}
	}
	if !foundUser || !foundTenant {
		t.Fatalf("expected attributes app.user_id and app.tenant_id set, got %v", attrs)
	}
}
