package bootstrap

import (
	"github.com/brunojet/go-infra-backend/internal/infra/observability/adapters"
	bootcontracts "github.com/brunojet/go-infra-backend/pkg/bootstrap/contracts"
)

// InitObservability initializes observability adapters from environment and
// registers their Shutdown handlers on the provided ShutdownManager.
func InitObservability(sm bootcontracts.ShutdownManager) error {
	ctx := sm.GetContext()

	if tracer, err := adapters.NewOTLPTracerFromEnv(ctx); err != nil {
		return err
	} else if tracer != nil {
		sm.Register("otel-tracer", tracer)
	}

	if metrics, err := adapters.NewOTLPMetricFromEnv(ctx); err != nil {
		return err
	} else if metrics != nil {
		sm.Register("otel-metric", metrics)
	}

	if logger, err := adapters.NewOTLPLoggerFromEnv(ctx); err != nil {
		return err
	} else if logger != nil {
		sm.Register("otel-logger", logger)
	}

	return nil
}
