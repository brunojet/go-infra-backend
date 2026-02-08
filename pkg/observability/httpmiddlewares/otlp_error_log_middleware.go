package httpmiddlewares

import (
	internalmiddlewares "github.com/brunojet/go-infra-backend/internal/observability/http_middlewares"
	"github.com/gin-gonic/gin"
)

type ErrorLogOption = internalmiddlewares.ErrorLogOption

var WithMinStatus = internalmiddlewares.WithMinStatus

// OTLPErrorLogMiddleware emits an OpenTelemetry LogRecord when a request finishes
// with a status code >= minStatus (default: 500) OR when gin collected errors.
//
// Important: register this middleware AFTER OtelGinMiddleware(), so it runs
// while the request context still contains the active span.
func OTLPErrorLogMiddleware(opts ...ErrorLogOption) gin.HandlerFunc {
	return internalmiddlewares.OTLPErrorLogMiddleware(opts...)
}
