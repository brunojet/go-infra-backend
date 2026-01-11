package observability

import (
	"fmt"

	httptypes "github.com/brunojet/go-infra-backend/internal/http/types"
	"github.com/brunojet/go-infra-backend/internal/observability/adapters/ginmw"
	"github.com/brunojet/go-infra-backend/internal/observability/adapters/httpmw"
	infrotel "github.com/brunojet/go-infra-backend/internal/observability/telemetry"
	"github.com/brunojet/go-infra-backend/internal/observability/types"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5/middleware"
)

// ginMiddlewares returns the configured Gin middleware chain.
func ginMiddlewares(cfg types.MiddlewareConfig) []types.GinMiddleware {
	tel := infrotel.NewStdProvider()

	var mws []types.GinMiddleware
	if cfg.RequestID {
		mws = append(mws, ginmw.RequestID())
	}
	if cfg.AccessLog {
		mws = append(mws, ginmw.AccessLog())
	}
	if cfg.Telemetry {
		mws = append(mws, ginmw.Telemetry(tel))
	}
	if cfg.Recovery {
		mws = append(mws, gin.Recovery())
	}
	return mws
}

// netHTTPMiddlewares returns the configured net/http middleware chain.
func netHTTPMiddlewares(cfg types.MiddlewareConfig) []types.NetHTTPMiddleware {
	tel := infrotel.NewStdProvider()

	var mws []types.NetHTTPMiddleware
	if cfg.RequestID {
		mws = append(mws, httpmw.RequestID)
	}
	if cfg.AccessLog {
		mws = append(mws, httpmw.AccessLog)
	}
	if cfg.Telemetry {
		mws = append(mws, httpmw.Telemetry(tel))
	}
	if cfg.Recovery {
		mws = append(mws, middleware.Recoverer)
	}
	return mws
}

func BuildMiddlewares(driver httptypes.HTTPDriver, cfg types.MiddlewareConfig) (types.Middlewares, error) {
	if !driver.IsSupported() {
		return types.Middlewares{}, fmt.Errorf("unsupported HTTP_DRIVER %q", driver)
	}

	switch driver {
	case httptypes.HTTPDriverChi:
		return types.Middlewares{NetHTTP: netHTTPMiddlewares(cfg)}, nil
	default:
		return types.Middlewares{Gin: ginMiddlewares(cfg)}, nil
	}
}
