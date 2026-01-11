package observability

import (
	"sync"

	httptypes "github.com/brunojet/go-infra-backend/internal/http/types"
	"github.com/brunojet/go-infra-backend/internal/observability/types"
)

// ObservabilityParams describes how to build observability middlewares
// for a selected HTTP driver.
type ObservabilityParams struct {
	// Driver is the HTTP driver (e.g. "gin", "chi").
	Driver httptypes.HTTPDriver
	// Config enables/disables features.
	Config types.MiddlewareConfig
}

// ObservabilityManager owns the computed middleware set for a given driver+config.
//
// It caches the computed slice(s) since middleware construction can allocate.
type ObservabilityManager struct {
	Params      ObservabilityParams
	mu          sync.Mutex
	middlewares *types.Middlewares
}

func NewObservabilityManager(params ObservabilityParams) *ObservabilityManager {
	return &ObservabilityManager{Params: params}
}

func (m *ObservabilityManager) Open() (types.Middlewares, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.middlewares != nil {
		return *m.middlewares, nil
	}

	mw, err := BuildMiddlewares(m.Params.Driver, m.Params.Config)
	if err != nil {
		return types.Middlewares{}, err
	}
	m.middlewares = &mw
	return mw, nil
}

func (m *ObservabilityManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.middlewares = nil
	return nil
}
