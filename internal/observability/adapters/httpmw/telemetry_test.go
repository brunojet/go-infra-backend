package httpmw

import (
	"net/http"
	"net/http/httptest"
	"testing"

	nooptelemetry "github.com/brunojet/go-infra-backend/internal/observability/adapters/nooptelemetry"
)

func TestTelemetryMiddleware(t *testing.T) {
	h := Telemetry(nooptelemetry.NoopProvider{})
	handler := h(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	req := httptest.NewRequest("GET", "/baz", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestTelemetryMiddleware_InternalServerError(t *testing.T) {
	h := Telemetry(nooptelemetry.NoopProvider{})
	handler := h(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	req := httptest.NewRequest("GET", "/fail", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != 500 {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestTelemetryMiddleware_NotFound(t *testing.T) {
	h := Telemetry(nooptelemetry.NoopProvider{})
	handler := h(http.NotFoundHandler())
	req := httptest.NewRequest("GET", "/notfound", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != 404 {
		t.Errorf("expected 404, got %d", w.Code)
	}
}
