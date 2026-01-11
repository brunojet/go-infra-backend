package ginmw

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAccessLog(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AccessLog())
	r.GET("/foo", func(c *gin.Context) { c.String(200, "ok") })
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/foo", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequestID())
	r.GET("/bar", func(c *gin.Context) {
		id := c.Writer.Header().Get("X-Request-ID")
		if id == "" {
			t.Error("expected X-Request-ID header to be set")
		}
		c.String(200, id)
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/bar", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
	id := w.Header().Get("X-Request-ID")
	if id == "" {
		t.Error("expected X-Request-ID header in response")
	}
	// Test with header provided
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("GET", "/bar", nil)
	req2.Header.Set("X-Request-ID", "abc123")
	r.ServeHTTP(w2, req2)
	if w2.Header().Get("X-Request-ID") != "abc123" {
		t.Error("should preserve provided X-Request-ID")
	}
}
