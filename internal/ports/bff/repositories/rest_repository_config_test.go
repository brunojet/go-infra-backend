package repositories

import (
	"errors"
	"net/http"
	"testing"
)

func assertPanics(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatalf("expected panic, got none")
		}
	}()
	fn()
}

func TestNewRestConfig_Valid(t *testing.T) {
	extractor := func(_ *http.Response) (string, error) { return "ok", nil }
	cfg := newRestConfig(
		WithBaseURL("https://example.com/api"),
		WithPathConfig("items"),
		WithEnvelopConfig("data", "meta"),
		WithMessageExtractor(extractor),
	)

	if cfg.URL.String() != "https://example.com/api" {
		t.Fatalf("unexpected URL: %s", cfg.URL.String())
	}
	if cfg.pathConfig.instancePath != "items" {
		t.Fatalf("unexpected instancePath: %s", cfg.pathConfig.instancePath)
	}
	if cfg.envelopConfig.dataField != "data" || cfg.envelopConfig.metaField != "meta" {
		t.Fatalf("unexpected envelope config: %+v", cfg.envelopConfig)
	}
	if cfg.messageExtractor == nil {
		t.Fatalf("messageExtractor should be configured")
	}
}

func TestWithBaseURL_InvalidPanics(t *testing.T) {
	assertPanics(t, func() {
		_ = newRestConfig(
			WithBaseURL("://invalid-url"),
			WithPathConfig("items"),
		)
	})
}

func TestWithPathConfig_EmptyCollectionParentPanics(t *testing.T) {
	assertPanics(t, func() {
		_ = newRestConfig(
			WithBaseURL("https://example.com"),
			WithPathConfig("items", "users", ""),
		)
	})
}

func TestWithPathConfig_BuildsCollectionParentsFmt(t *testing.T) {
	cfg := newRestConfig(
		WithBaseURL("https://example.com"),
		WithPathConfig("items", "users", "orders"),
	)

	if cfg.pathConfig.collectionParentsFmt != "users/%s/orders/%s" {
		t.Fatalf("unexpected collectionParentsFmt: %s", cfg.pathConfig.collectionParentsFmt)
	}
	if cfg.pathConfig.collectionParents != 2 {
		t.Fatalf("unexpected collectionParents: %d", cfg.pathConfig.collectionParents)
	}
}

func TestWithEnvelopConfig_MetaWithoutDataPanics(t *testing.T) {
	assertPanics(t, func() {
		_ = newRestConfig(
			WithBaseURL("https://example.com"),
			WithPathConfig("items"),
			WithEnvelopConfig("", "meta"),
		)
	})
}

func TestWithMessageExtractor_CallbackIsUsed(t *testing.T) {
	extractor := func(_ *http.Response) (string, error) {
		return "mapped", errors.New("ignored for this assertion")
	}
	cfg := newRestConfig(
		WithBaseURL("https://example.com"),
		WithPathConfig("items"),
		WithMessageExtractor(extractor),
	)
	if cfg.messageExtractor == nil {
		t.Fatalf("expected callback to be set")
	}
}
