package providers

import (
	"context"
	"os"
	"testing"
)

func TestGetServiceNameAndVersion_Defaults(t *testing.T) {
	os.Unsetenv(ServiceNameEnv)
	os.Unsetenv(ServiceVersionEnv)

	if got := GetServiceName(); got == "" {
		t.Fatalf("expected non-empty service name")
	}
	if got := GetServiceVersion(); got == "" {
		t.Fatalf("expected non-empty service version")
	}
}

func TestGetServiceInfo_FromEnv(t *testing.T) {
	os.Setenv(ServiceNameEnv, "my-service")
	os.Setenv(ServiceVersionEnv, "v9.9.9")
	defer func() {
		os.Unsetenv(ServiceNameEnv)
		os.Unsetenv(ServiceVersionEnv)
	}()

	info := GetServiceInfo()
	if info.Name != "my-service" {
		t.Fatalf("unexpected name: %s", info.Name)
	}
	if info.Version != "v9.9.9" {
		t.Fatalf("unexpected version: %s", info.Version)
	}
}

func TestNewOTLPResource_Builds(t *testing.T) {
	os.Setenv(ServiceNameEnv, "srv-test")
	os.Setenv(ServiceVersionEnv, "0.1.2")
	defer func() {
		os.Unsetenv(ServiceNameEnv)
		os.Unsetenv(ServiceVersionEnv)
	}()

	res, info, err := NewOTLPResource(context.Background())
	if err != nil {
		t.Fatalf("NewOTLPResource error: %v", err)
	}
	if info.Name != "srv-test" || info.Version != "0.1.2" {
		t.Fatalf("unexpected ServiceInfo: %+v", info)
	}
	if res == nil {
		t.Fatalf("resource is nil")
	}
}
