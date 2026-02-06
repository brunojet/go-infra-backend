package providers

import (
	"context"
	"os"
	"testing"

	"github.com/brunojet/go-infra-backend/internal/observability/exporters"
)

func TestNewOTLPMetricProvider(t *testing.T) {
	os.Unsetenv(exporters.OTLPEndpointEnv)
	ctx := context.Background()
	exporter, err := exporters.NewOTLPMetricExporter(ctx)
	if err != nil {
		t.Fatalf("NewOTLPMetricExporter error: %v", err)
	}
	mp, shutdown, err := NewOTLPMetricProvider(ctx, exporter)
	if err != nil {
		t.Fatalf("NewOTLPMetricProvider error: %v", err)
	}
	if mp == nil {
		t.Fatalf("expected meter provider, got nil")
	}
	if shutdown == nil {
		t.Fatalf("expected shutdown function, got nil")
	}
	if err := shutdown(ctx); err != nil {
		t.Fatalf("shutdown error: %v", err)
	}
}
