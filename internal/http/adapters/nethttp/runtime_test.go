package nethttp

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRuntime_Default(t *testing.T) {
	router, handler := NewRuntime()
	if router == nil || handler == nil {
		t.Error("expected non-nil router and handler")
	}
}

func TestNewRuntime_WithMiddleware(t *testing.T) {
	mwCalled := false
	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mwCalled = true
			next.ServeHTTP(w, r)
		})
	}
	router, handler := NewRuntime(WithMiddlewares(mw))
	if router == nil || handler == nil {
		t.Error("expected non-nil router and handler")
	}
	// Registra uma rota para garantir que o handler e o middleware sejam chamados
	if chiRouter, ok := router.(*ChiRouter); ok {
		chiRouter.r.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
		}))
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(w, r)
	if !mwCalled {
		t.Error("expected middleware to be called")
	}
}
