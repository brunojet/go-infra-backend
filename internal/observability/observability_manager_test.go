package observability

import (
	"testing"

	httptypes "github.com/brunojet/go-infra-backend/internal/http/types"
	"github.com/brunojet/go-infra-backend/internal/observability/types"
)

func TestObservabilityManager_OpenAndClose(t *testing.T) {
	params := ObservabilityParams{
		Driver: httptypes.HTTPDriverGin,
		Config: types.MiddlewareConfig{RequestID: true, AccessLog: true, Telemetry: true, Recovery: true},
	}
	mgr := NewObservabilityManager(params)
	mws, err := mgr.Open()
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	if len(mws.Gin) == 0 {
		t.Error("expected Gin middlewares")
	}
	// Test cache
	mws2, err := mgr.Open()
	if err != nil || len(mws2.Gin) == 0 {
		t.Error("expected cached middlewares")
	}
	if &mws == &mws2 {
		t.Error("should return a copy, not the same pointer")
	}
	if err := mgr.Close(); err != nil {
		t.Errorf("Close failed: %v", err)
	}
}

func TestObservabilityManager_OpenUnsupported(t *testing.T) {
	params := ObservabilityParams{
		Driver: "unsupported",
		Config: types.MiddlewareConfig{},
	}
	mgr := NewObservabilityManager(params)
	_, err := mgr.Open()
	if err == nil {
		t.Error("expected error for unsupported driver")
	}
}

func TestBuildMiddlewares_Chi(t *testing.T) {
	cfg := types.MiddlewareConfig{RequestID: true, AccessLog: true, Telemetry: true, Recovery: true}
	mws, err := BuildMiddlewares(httptypes.HTTPDriverChi, cfg)
	if err != nil {
		t.Fatalf("BuildMiddlewares failed: %v", err)
	}
	if len(mws.NetHTTP) == 0 {
		t.Error("expected NetHTTP middlewares")
	}
}

func TestBuildMiddlewares_Unsupported(t *testing.T) {
	_, err := BuildMiddlewares("invalid", types.MiddlewareConfig{})
	if err == nil {
		t.Error("expected error for unsupported driver")
	}
}
