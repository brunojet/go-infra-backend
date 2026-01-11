package http

import (
	"errors"
	"fmt"
	stdhttp "net/http"
	"sync"

	ginadapter "github.com/brunojet/go-infra-backend/internal/http/adapters/gin"
	"github.com/brunojet/go-infra-backend/internal/http/adapters/nethttp"
	"github.com/brunojet/go-infra-backend/internal/http/contracts"
	"github.com/brunojet/go-infra-backend/internal/http/types"
	obstypes "github.com/brunojet/go-infra-backend/internal/observability/types"
	"github.com/gin-gonic/gin"
)

var ErrMissingRegistrar = errors.New("register enabled but no registrar provided")

// HTTPParams describes how to initialize an HTTP runtime.
//
// Driver selection is delegated to Select.
// Middlewares and CORS are wired through runtime options.
//
// Register is intentionally kept here as a convenience flag for higher-level bootstraps;
// the actual route registration is performed by HTTPManager via a provided callback.
type HTTPParams struct {
	Driver                   types.HTTPDriver
	CORS                     types.CORSConfig
	ObservabilityMiddlewares obstypes.Middlewares
	Register                 bool
}

// HttpRuntime wires an HTTP framework into a framework-agnostic Router contract.
//
// Router is where modules register routes. Handler is what net/http servers serve.
//
// This package is infra: it owns framework selection/initialization.
type HttpRuntime struct {
	Router  contracts.Router
	Handler stdhttp.Handler
}

// HTTPManager owns an HttpRuntime instance and can optionally run a route-registration
// callback when Params.Register is true.
//
// Typical usage:
//
//	mgr, _ := http.NewHTTPManager(http.HTTPParams{Driver: http.HTTPDriverGin, Register: true})
//	rt, _ := mgr.OpenAndRegister(func(r contracts.Router) error {
//		r.GET("/health", ...)
//		return nil
//	})
type HTTPManager struct {
	Params  HTTPParams
	mu      sync.Mutex
	runtime *HttpRuntime
}

func NewHTTPManager(params HTTPParams) (*HTTPManager, error) {
	m := &HTTPManager{Params: params}
	// Eagerly validate driver/params.
	if _, err := m.Open(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *HTTPManager) Open() (*HttpRuntime, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.runtime != nil {
		return m.runtime, nil
	}

	rt, err := m.buildRuntime()
	if err != nil {
		return nil, err
	}

	m.runtime = rt
	return rt, nil
}

func (m *HTTPManager) buildRuntime() (*HttpRuntime, error) {
	switch m.Params.Driver {
	case types.HTTPDriverGin:
		mws := m.buildGinMiddlewares()
		router, handler := ginadapter.NewRuntime(ginadapter.WithMiddlewares(mws...))
		return &HttpRuntime{Router: router, Handler: handler}, nil
	case types.HTTPDriverChi:
		mws := m.buildNetHTTPMiddlewares()
		router, handler := nethttp.NewRuntime(nethttp.WithMiddlewares(mws...))
		return &HttpRuntime{Router: router, Handler: handler}, nil
	default:
		return nil, fmt.Errorf("unsupported HTTP_DRIVER %q", m.Params.Driver)
	}
}

func (m *HTTPManager) buildGinMiddlewares() []gin.HandlerFunc {
	mws := append([]gin.HandlerFunc(nil), m.Params.ObservabilityMiddlewares.Gin...)
	if m.Params.CORS.Enabled {
		mws = append(mws, ginadapter.CORS(m.Params.CORS))
	}
	return mws
}

func (m *HTTPManager) buildNetHTTPMiddlewares() []func(stdhttp.Handler) stdhttp.Handler {
	mws := append([]func(stdhttp.Handler) stdhttp.Handler(nil), m.Params.ObservabilityMiddlewares.NetHTTP...)
	if m.Params.CORS.Enabled {
		mws = append(mws, nethttp.CORS(m.Params.CORS))
	}
	return mws
}

func (m *HTTPManager) OpenAndRegister(registrar func(contracts.Router) error) (*HttpRuntime, error) {
	rt, err := m.Open()
	if err != nil {
		return nil, err
	}
	if !m.Params.Register {
		return rt, nil
	}
	if registrar == nil {
		return nil, ErrMissingRegistrar
	}
	if err := registrar(rt.Router); err != nil {
		return nil, err
	}
	return rt, nil
}

func (m *HTTPManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.runtime = nil
	return nil
}
