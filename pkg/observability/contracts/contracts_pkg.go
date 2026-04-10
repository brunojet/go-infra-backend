package contracts

import (
	"context"
	"time"
)

const (
	minErrorsForDegraded  = 3
	minErrorsForUnhealthy = 10
	maxErrorsInArray      = 10
)

type HealthState string

const (
	HealthStatePending   HealthState = "pending"
	HealthStateHealthy   HealthState = "healthy"
	HealthStateDegraded  HealthState = "degraded"
	HealthStateUnhealthy HealthState = "unhealthy"
	HealthStateDefault   HealthState = HealthStatePending
)

type ObservabilityAdapterError struct {
	Timestamp time.Time `json:"timestamp"`
	Error     string    `json:"error"`
}

type ObservabilityAdapterHealth struct {
	HealthState HealthState                 `json:"health_state"`
	LastErrors  []ObservabilityAdapterError `json:"last_errors"`
}

type ObservabilityAdapterStats struct {
	Endpoint   string `json:"endpoint"`
	IsInsecure bool   `json:"is_insecure"`
	Onqueue    int    `json:"onqueue"`
	Sent       int    `json:"sent"`
	Failed     int    `json:"failed"`
}

type ObservabilityAdapter interface {
	Shutdown(ctx context.Context) error
	GetAdapterStats() []ObservabilityAdapterStats
}

type ObservabilityManager interface {
	RegisterAdapter(adapter ObservabilityAdapter)
	GetAdapterStats() []ObservabilityAdapterStats
	Shutdown(ctx context.Context) error
}

// MaxErrorsInArray returns the configured default capacity for recent errors buffer.
func MaxErrorsInArray() int {
	return maxErrorsInArray
}
