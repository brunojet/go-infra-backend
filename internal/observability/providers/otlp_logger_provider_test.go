package providers

import (
	"context"
	"testing"

	"github.com/brunojet/go-infra-backend/internal/observability/exporters"
)

func TestNewOTLPLoggerProvider(t *testing.T) {
	exporter := exporters.NewOTLPNoopLoggerExporter()

	lp, shutdown, err := NewOTLPLoggerProvider(context.Background(), exporter)
	if err != nil {
		t.Fatalf("NewOTLPLoggerProvider error: %v", err)
	}
	if lp == nil {
		t.Fatalf("expected logger provider, got nil")
	}
	if shutdown == nil {
		t.Fatalf("expected shutdown function, got nil")
	}
	if err := shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown error: %v", err)
	}
}
