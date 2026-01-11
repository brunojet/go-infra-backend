package observability

import (
	"github.com/brunojet/go-infra-backend/internal/observability"
	observabilitytypes "github.com/brunojet/go-infra-backend/internal/observability/types"
)

type ObservabilityParams = observability.ObservabilityParams

type ObservabilityManager = observability.ObservabilityManager

var NewObservabilityManager = observability.NewObservabilityManager

// Export MiddlewareConfig type
type MiddlewareConfig = observabilitytypes.MiddlewareConfig
