package adapters

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestOTLPTracingExporter_NoopSpan(t *testing.T) {
	e := NewOTLPTracingExporter("")
	ctx, end := e.StartSpan(context.Background(), "unit-test-span")
	if ctx == nil {
		t.Fatal("expected context returned")
	}
	end()
}

func TestOTLPTracingExporter_StartShutdown(t *testing.T) {
	// skip if no local OTLP gRPC server available
	conn, err := net.DialTimeout("tcp", "localhost:4317", 200*time.Millisecond)
	if err != nil {
		t.Skip("no local OTLP gRPC server on localhost:4317, skipping integration Start/Shutdown test")
	}
	_ = conn.Close()

	e := NewOTLPTracingExporter("")
	if err := e.Start(context.Background()); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	_, end := e.StartSpan(context.Background(), "integration-span")
	end()
	if err := e.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}
}
