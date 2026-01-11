package ginmw

import (
	"net/http/httptest"
	"testing"

	noopadapter "github.com/brunojet/go-infra-backend/internal/observability/adapters/nooptelemetry"
	"github.com/gin-gonic/gin"
)

func TestTelemetryMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Telemetry(noopadapter.NoopProvider{}))
	r.GET("/baz", func(c *gin.Context) { c.String(200, "baz") })
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/baz", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestTelemetryMiddleware_InternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Telemetry(noopadapter.NoopProvider{}))
	r.GET("/fail", func(c *gin.Context) { c.Status(500) })
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/fail", nil)
	r.ServeHTTP(w, req)
	if w.Code != 500 {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestTelemetryMiddleware_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Telemetry(noopadapter.NoopProvider{}))
	// Não registra rota, para forçar 404
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/notfound", nil)
	r.ServeHTTP(w, req)
	if w.Code != 404 {
		t.Errorf("expected 404, got %d", w.Code)
	}
}
