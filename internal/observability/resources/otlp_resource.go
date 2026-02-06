package resources

import (
	"context"
	"os"

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

const (
	ServiceNameEnv        = "SERVICE_NAME"
	ServiceNameDefault    = "meu-servico"
	ServiceVersionEnv     = "SERVICE_VERSION"
	ServiceVersionDefault = "1.0.0"
)

// ServiceInfo holds the service identity used for telemetry resources.
type ServiceInfo struct {
	Name    string
	Version string
}

// GetServiceName returns the service name from environment or the default.
func GetServiceName() string {
	return getEnv(ServiceNameEnv, ServiceNameDefault)
}

// GetServiceVersion returns the service version from environment or the default.
func GetServiceVersion() string {
	return getEnv(ServiceVersionEnv, ServiceVersionDefault)
}

// GetServiceInfo returns a ServiceInfo containing name and version.
func GetServiceInfo() ServiceInfo {
	return ServiceInfo{
		Name:    GetServiceName(),
		Version: GetServiceVersion(),
	}
}

// NewOTLPResource builds an OpenTelemetry Resource using service name/version
// read from environment (SERVICE_NAME, SERVICE_VERSION) with sensible defaults.
func NewOTLPResource(ctx context.Context) (*resource.Resource, ServiceInfo, error) {
	info := GetServiceInfo()

	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName(info.Name),
			semconv.ServiceVersion(info.Version),
		),
	)
	if err != nil {
		return nil, ServiceInfo{}, err
	}

	return res, info, nil
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
