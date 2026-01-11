package httpmw

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccessLog(t *testing.T) {
	h := AccessLog(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
	}))
	req := httptest.NewRequest("POST", "/foo", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 201 {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestRequestID(t *testing.T) {
	h := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := w.Header().Get("X-Request-ID")
		if id == "" {
			t.Error("expected X-Request-ID header to be set")
		}
	}))
	req := httptest.NewRequest("GET", "/bar", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Header().Get("X-Request-ID") == "" {
		t.Error("expected X-Request-ID header in response")
	}
	// Test with header provided
	req2 := httptest.NewRequest("GET", "/bar", nil)
	req2.Header.Set("X-Request-ID", "abc123")
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, req2)
	if w2.Header().Get("X-Request-ID") != "abc123" {
		t.Error("should preserve provided X-Request-ID")
	}
}
