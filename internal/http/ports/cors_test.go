package ports

import (
	"reflect"
	"testing"
	"time"

	httptypes "github.com/brunojet/go-infra-backend/internal/http/types"
)

func TestBuildCORSHeaders_EmptyOrigin(t *testing.T) {
	cfg := httptypes.CORSConfig{AllowOrigins: []string{"*"}, AllowMethods: []string{"GET"}}
	_, ok := BuildCORSHeaders(cfg, "")
	if ok {
		t.Fatalf("expected ok=false for empty origin")
	}
}

func TestBuildCORSHeaders_WildcardNoCredentials_AllowsAll(t *testing.T) {
	cfg := httptypes.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    nil,
		AllowCredentials: false,
		MaxAge:           10 * time.Minute,
	}

	h, ok := BuildCORSHeaders(cfg, "https://example.com")
	if !ok {
		t.Fatalf("expected ok=true")
	}
	if h.AllowOrigin != "*" {
		t.Fatalf("expected AllowOrigin=*, got %q", h.AllowOrigin)
	}
	if h.VaryOrigin {
		t.Fatalf("expected VaryOrigin=false")
	}
	if h.AllowMethods != "GET, POST" {
		t.Fatalf("unexpected AllowMethods: %q", h.AllowMethods)
	}
	if h.AllowHeaders != "Content-Type" {
		t.Fatalf("unexpected AllowHeaders: %q", h.AllowHeaders)
	}
	if h.AllowCredentials {
		t.Fatalf("expected AllowCredentials=false")
	}
	if h.MaxAgeSeconds != 600 {
		t.Fatalf("expected MaxAgeSeconds=600, got %d", h.MaxAgeSeconds)
	}
}

func TestBuildCORSHeaders_WildcardWithCredentials_IsRejected(t *testing.T) {
	cfg := httptypes.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           time.Minute,
	}

	_, ok := BuildCORSHeaders(cfg, "https://example.com")
	if ok {
		t.Fatalf("expected ok=false when credentials=true and origins contains '*'")
	}
}

func TestBuildCORSHeaders_ExplicitOriginWithCredentials_Varies(t *testing.T) {
	cfg := httptypes.CORSConfig{
		AllowOrigins:     []string{"https://a.com"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"x-out"},
		AllowCredentials: true,
		MaxAge:           0,
	}

	h, ok := BuildCORSHeaders(cfg, "https://a.com")
	if !ok {
		t.Fatalf("expected ok=true")
	}
	if h.AllowOrigin != "https://a.com" {
		t.Fatalf("expected AllowOrigin=https://a.com, got %q", h.AllowOrigin)
	}
	if !h.VaryOrigin {
		t.Fatalf("expected VaryOrigin=true")
	}
	if !h.AllowCredentials {
		t.Fatalf("expected AllowCredentials=true")
	}
	if h.ExposeHeaders != "x-out" {
		t.Fatalf("unexpected ExposeHeaders: %q", h.ExposeHeaders)
	}
	if h.MaxAgeSeconds != 0 {
		t.Fatalf("expected MaxAgeSeconds=0, got %d", h.MaxAgeSeconds)
	}
}

func TestBuildCORSHeaders_DisallowedOrigin(t *testing.T) {
	cfg := httptypes.CORSConfig{AllowOrigins: []string{"https://a.com"}}
	_, ok := BuildCORSHeaders(cfg, "https://b.com")
	if ok {
		t.Fatalf("expected ok=false for disallowed origin")
	}
}

func TestBuildCORSHeaders_NegativeMaxAge_IsClamped(t *testing.T) {
	cfg := httptypes.CORSConfig{AllowOrigins: []string{"*"}, MaxAge: -5 * time.Second}
	h, ok := BuildCORSHeaders(cfg, "https://example.com")
	if !ok {
		t.Fatalf("expected ok=true")
	}
	if h.MaxAgeSeconds != 0 {
		t.Fatalf("expected MaxAgeSeconds=0, got %d", h.MaxAgeSeconds)
	}
}

func TestBuildCORSHeaders_ExposeHeaders_EmptyStringWhenNil(t *testing.T) {
	cfg := httptypes.CORSConfig{AllowOrigins: []string{"*"}}
	h, ok := BuildCORSHeaders(cfg, "https://example.com")
	if !ok {
		t.Fatalf("expected ok=true")
	}
	if h.ExposeHeaders != "" {
		t.Fatalf("expected ExposeHeaders empty, got %q", h.ExposeHeaders)
	}
}

func TestCORSHeaders_MaxAgeString(t *testing.T) {
	h := CORSHeaders{MaxAgeSeconds: 123}
	if got := h.MaxAgeString(); got != "123" {
		t.Fatalf("expected 123, got %q", got)
	}
}

func TestBuildCORSHeaders_AllowHeadersAndMethods_JoinNilSlices(t *testing.T) {
	cfg := httptypes.CORSConfig{AllowOrigins: []string{"*"}, AllowMethods: nil, AllowHeaders: nil}
	h, ok := BuildCORSHeaders(cfg, "https://example.com")
	if !ok {
		t.Fatalf("expected ok=true")
	}
	if !reflect.DeepEqual(h.AllowMethods, "") {
		t.Fatalf("expected empty AllowMethods, got %q", h.AllowMethods)
	}
	if !reflect.DeepEqual(h.AllowHeaders, "") {
		t.Fatalf("expected empty AllowHeaders, got %q", h.AllowHeaders)
	}
}
