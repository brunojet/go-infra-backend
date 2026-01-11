package nethttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"time"

	httptypes "github.com/brunojet/go-infra-backend/internal/http/types"
)

func TestCORS_AllHeadersAndBranches(t *testing.T) {
	cfg := httptypes.CORSConfig{
		AllowOrigins:     []string{"http://test"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"X-Test"},
		ExposeHeaders:    []string{"X-Expose"},
		AllowCredentials: true,
		MaxAge:           10 * time.Second,
	}
	mw := CORS(cfg)
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://test")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	resp := w.Result()
	headers := resp.Header
	if headers.Get("Access-Control-Allow-Origin") != "http://test" {
		t.Errorf("expected Allow-Origin header")
	}
	if headers.Get("Vary") != "Origin" {
		t.Errorf("expected Vary header")
	}
	if headers.Get("Access-Control-Allow-Methods") != "GET, POST" {
		t.Errorf("expected Allow-Methods header")
	}
	if headers.Get("Access-Control-Allow-Headers") != "X-Test" {
		t.Errorf("expected Allow-Headers header")
	}
	if headers.Get("Access-Control-Expose-Headers") != "X-Expose" {
		t.Errorf("expected Expose-Headers header")
	}
	if headers.Get("Access-Control-Allow-Credentials") != "true" {
		t.Errorf("expected Allow-Credentials header")
	}
	if headers.Get("Access-Control-Max-Age") == "" {
		t.Errorf("expected Max-Age header")
	}

	// Test OPTIONS preflight
	req = httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("Origin", "http://test")
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected 204 for OPTIONS, got %d", resp.StatusCode)
	}

	// Test with wildcard origin and credentials=false
	cfg = httptypes.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET"},
		AllowHeaders: []string{"X-Test"},
		MaxAge:       5 * time.Second,
	}
	mw = CORS(cfg)
	h = mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "any")
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)
	resp = w.Result()
	if headers.Get("Access-Control-Allow-Origin") == "*" && headers.Get("Vary") != "" {
		t.Errorf("wildcard origin should not set Vary header")
	}

	// Test with credentials=true and wildcard origin (should not set CORS headers)
	cfg = httptypes.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
	}
	mw = CORS(cfg)
	h = mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "any")
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)
	resp = w.Result()
	if resp.Header.Get("Access-Control-Allow-Origin") != "" {
		t.Errorf("should not set CORS headers when credentials=true and origin=*")
	}

	// Test with no Origin header (should not set CORS headers)
	cfg = httptypes.CORSConfig{
		AllowOrigins: []string{"http://test"},
	}
	mw = CORS(cfg)
	h = mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)
	resp = w.Result()
	if resp.Header.Get("Access-Control-Allow-Origin") != "" {
		t.Errorf("should not set CORS headers when no Origin header")
	}

	// Test with not allowed origin (should not set CORS headers)
	cfg = httptypes.CORSConfig{
		AllowOrigins: []string{"http://allowed"},
	}
	mw = CORS(cfg)
	h = mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://notallowed")
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)
	resp = w.Result()
	if resp.Header.Get("Access-Control-Allow-Origin") != "" {
		t.Errorf("should not set CORS headers for not allowed origin")
	}
}
