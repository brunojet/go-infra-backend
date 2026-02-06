package providers

import (
	"context"
	"testing"

	"github.com/brunojet/go-infra-backend/internal/observability/exporters"
)

func TestNewOTLPTracerProvider(t *testing.T) {
	exp := exporters.NewOTLPNoopTracerExporter()

	tp, shutdown, err := NewOTLPTracerProvider(context.Background(), exp)
	if err != nil {
		t.Fatalf("NewOTLPTracerProvider error: %v", err)
	}
	if tp == nil {
		t.Fatalf("expected tracer provider, got nil")
	}
	if shutdown == nil {
		t.Fatalf("expected shutdown function, got nil")
	}
	if err := shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown error: %v", err)
	}
}
