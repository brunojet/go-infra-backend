package adapters

import (
	"context"
	"testing"
)

func TestOTLPLoggingExporter_Emit(t *testing.T) {
	e := NewOTLPLoggingExporter()
	if err := e.Start(context.Background()); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	if err := e.Emit(context.Background(), "info", "test log body"); err != nil {
		t.Fatalf("Emit failed: %v", err)
	}
	if err := e.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}
}
