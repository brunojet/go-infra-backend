package nethttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	contracts "github.com/brunojet/go-infra-backend/internal/http/contracts"
	"github.com/go-chi/chi/v5"
)

func TestNewChiRouter_BasicRoutes(t *testing.T) {
	r := chi.NewRouter()
	router := NewChiRouter(r)
	called := ""
	router.GET("/foo", contracts.HandlerFunc(func(ctx contracts.Context) { called = "get" }))
	router.POST("/bar", contracts.HandlerFunc(func(ctx contracts.Context) { called = "post" }))
	router.PATCH("/baz", contracts.HandlerFunc(func(ctx contracts.Context) { called = "patch" }))
	router.DELETE("/qux", contracts.HandlerFunc(func(ctx contracts.Context) { called = "delete" }))

	tests := []struct {
		method, path, expect string
	}{
		{"GET", "/foo", "get"},
		{"POST", "/bar", "post"},
		{"PATCH", "/baz", "patch"},
		{"DELETE", "/qux", "delete"},
	}
	for _, tt := range tests {
		called = ""
		req, _ := http.NewRequest(tt.method, tt.path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if called != tt.expect {
			t.Errorf("expected %s handler to be called, got %s", tt.expect, called)
		}
	}
}

func TestChiRouter_Group(t *testing.T) {
	r := chi.NewRouter()
	router := NewChiRouter(r)
	group := router.Group("/api")
	if group == nil {
		t.Error("expected non-nil group router")
	}
}
