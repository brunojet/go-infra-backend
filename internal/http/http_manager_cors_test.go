package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brunojet/go-infra-backend/internal/http/contracts"
	httptypes "github.com/brunojet/go-infra-backend/internal/http/types"
)

func registerPing(r contracts.Router) error {
	r.GET("/ping", func(c contracts.Context) {
		c.String(http.StatusOK, "ok")
	})
	return nil
}

func doRequest(t *testing.T, h http.Handler, method, path, origin string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, "http://example.com"+path, nil)
	if origin != "" {
		req.Header.Set("Origin", origin)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}

func TestCORS_DefaultDisabled_NoHeaders_Gin(t *testing.T) {
	mgr, err := NewHTTPManager(HTTPParams{
		Driver:   httptypes.HTTPDriverGin,
		CORS:     httptypes.CORSConfig{Enabled: false, AllowOrigins: []string{"*"}},
		Register: true,
	})
	if err != nil {
		t.Fatalf("NewHTTPManager: %v", err)
	}
	rt, err := mgr.OpenAndRegister(registerPing)
	if err != nil {
		t.Fatalf("OpenAndRegister: %v", err)
	}

	rr := doRequest(t, rt.Handler, http.MethodGet, "/ping", "https://a.com")
	if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no Access-Control-Allow-Origin, got %q", got)
	}
}

func TestCORS_DefaultDisabled_NoHeaders_Chi(t *testing.T) {
	mgr, err := NewHTTPManager(HTTPParams{
		Driver:   httptypes.HTTPDriverChi,
		CORS:     httptypes.CORSConfig{Enabled: false, AllowOrigins: []string{"*"}},
		Register: true,
	})
	if err != nil {
		t.Fatalf("NewHTTPManager: %v", err)
	}
	rt, err := mgr.OpenAndRegister(registerPing)
	if err != nil {
		t.Fatalf("OpenAndRegister: %v", err)
	}

	rr := doRequest(t, rt.Handler, http.MethodGet, "/ping", "https://a.com")
	if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no Access-Control-Allow-Origin, got %q", got)
	}
}

func TestCORS_Enabled_ExplicitOrigins_OnlyAllowed_Gin(t *testing.T) {
	mgr, err := NewHTTPManager(HTTPParams{
		Driver: httptypes.HTTPDriverGin,
		CORS: httptypes.CORSConfig{
			Enabled:          true,
			AllowOrigins:     []string{"https://a.com"},
			AllowMethods:     []string{"GET", "OPTIONS"},
			AllowHeaders:     []string{"Content-Type"},
			AllowCredentials: true,
			MaxAge:           time.Minute,
		},
		Register: true,
	})
	if err != nil {
		t.Fatalf("NewHTTPManager: %v", err)
	}
	rt, err := mgr.OpenAndRegister(registerPing)
	if err != nil {
		t.Fatalf("OpenAndRegister: %v", err)
	}

	allowed := doRequest(t, rt.Handler, http.MethodGet, "/ping", "https://a.com")
	if got := allowed.Header().Get("Access-Control-Allow-Origin"); got != "https://a.com" {
		t.Fatalf("expected ACAO=https://a.com, got %q", got)
	}
	if got := allowed.Header().Get("Vary"); got != "Origin" {
		t.Fatalf("expected Vary=Origin, got %q", got)
	}
	if got := allowed.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Fatalf("expected Allow-Credentials=true, got %q", got)
	}

	disallowed := doRequest(t, rt.Handler, http.MethodGet, "/ping", "https://b.com")
	if got := disallowed.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no ACAO for disallowed origin, got %q", got)
	}
}

func TestCORS_Enabled_ExplicitOrigins_OnlyAllowed_Chi(t *testing.T) {
	mgr, err := NewHTTPManager(HTTPParams{
		Driver: httptypes.HTTPDriverChi,
		CORS: httptypes.CORSConfig{
			Enabled:      true,
			AllowOrigins: []string{"https://a.com"},
			AllowMethods: []string{"GET", "OPTIONS"},
			AllowHeaders: []string{"Content-Type"},
		},
		Register: true,
	})
	if err != nil {
		t.Fatalf("NewHTTPManager: %v", err)
	}
	rt, err := mgr.OpenAndRegister(registerPing)
	if err != nil {
		t.Fatalf("OpenAndRegister: %v", err)
	}

	allowed := doRequest(t, rt.Handler, http.MethodGet, "/ping", "https://a.com")
	if got := allowed.Header().Get("Access-Control-Allow-Origin"); got != "https://a.com" {
		t.Fatalf("expected ACAO=https://a.com, got %q", got)
	}
	if got := allowed.Header().Get("Vary"); got != "Origin" {
		t.Fatalf("expected Vary=Origin, got %q", got)
	}

	disallowed := doRequest(t, rt.Handler, http.MethodGet, "/ping", "https://b.com")
	if got := disallowed.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no ACAO for disallowed origin, got %q", got)
	}
}

func TestCORS_PreflightOPTIONS_204_Gin(t *testing.T) {
	mgr, err := NewHTTPManager(HTTPParams{
		Driver: httptypes.HTTPDriverGin,
		CORS: httptypes.CORSConfig{
			Enabled:      true,
			AllowOrigins: []string{"https://a.com"},
			AllowMethods: []string{"GET", "OPTIONS"},
			AllowHeaders: []string{"Content-Type"},
		},
		Register: true,
	})
	if err != nil {
		t.Fatalf("NewHTTPManager: %v", err)
	}
	rt, err := mgr.OpenAndRegister(registerPing)
	if err != nil {
		t.Fatalf("OpenAndRegister: %v", err)
	}

	rr := doRequest(t, rt.Handler, http.MethodOptions, "/ping", "https://a.com")
	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rr.Code)
	}
}

func TestCORS_PreflightOPTIONS_204_Chi(t *testing.T) {
	mgr, err := NewHTTPManager(HTTPParams{
		Driver: httptypes.HTTPDriverChi,
		CORS: httptypes.CORSConfig{
			Enabled:      true,
			AllowOrigins: []string{"https://a.com"},
			AllowMethods: []string{"GET", "OPTIONS"},
			AllowHeaders: []string{"Content-Type"},
		},
		Register: true,
	})
	if err != nil {
		t.Fatalf("NewHTTPManager: %v", err)
	}
	rt, err := mgr.OpenAndRegister(registerPing)
	if err != nil {
		t.Fatalf("OpenAndRegister: %v", err)
	}

	rr := doRequest(t, rt.Handler, http.MethodOptions, "/ping", "https://a.com")
	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rr.Code)
	}
}
