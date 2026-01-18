package adapters

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestOTLPMetricsExporter_RecordNoop(t *testing.T) {
	e := NewOTLPMetricsExporter()
	// Recording without Start should use global meter (no panic)
	if err := e.Record(context.Background(), "unit_test_metric", 1.23); err != nil {
		t.Fatalf("Record failed: %v", err)
	}
}

func TestOTLPMetricsExporter_StartShutdown(t *testing.T) {
	conn, err := net.DialTimeout("tcp", "localhost:4317", 200*time.Millisecond)
	if err != nil {
		t.Skip("no local OTLP gRPC server on localhost:4317, skipping integration Start/Shutdown test")
	}
	_ = conn.Close()

	e := NewOTLPMetricsExporter()
	if err := e.Start(context.Background()); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	if err := e.Record(context.Background(), "integration_metric", 42.0); err != nil {
		t.Fatalf("Record failed: %v", err)
	}
	if err := e.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}
}
