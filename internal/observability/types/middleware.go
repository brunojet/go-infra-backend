package types

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MiddlewareConfig defines which observability features are enabled.
//
// Defaults are applied by the app launcher.
type MiddlewareConfig struct {
	RequestID bool
	AccessLog bool
	Telemetry bool
	Recovery  bool
}

// GinMiddleware is a Gin middleware.
type GinMiddleware = gin.HandlerFunc

// NetHTTPMiddleware is a net/http middleware.
type NetHTTPMiddleware = func(http.Handler) http.Handler

// Middlewares is a single "package" of middlewares for the selected HTTP driver.
// Exactly one of Gin or NetHTTP should be populated.
type Middlewares struct {
	Gin     []GinMiddleware
	NetHTTP []NetHTTPMiddleware
}
