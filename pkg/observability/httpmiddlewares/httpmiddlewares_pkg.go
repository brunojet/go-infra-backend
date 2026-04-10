package httpmiddlewares

import (
	internalmiddlewares "github.com/brunojet/go-infra-backend/internal/observability/http_middlewares"
	"github.com/gin-gonic/gin"
)

type ErrorLogOption = internalmiddlewares.ErrorLogOption

var WithMinStatus = internalmiddlewares.WithMinStatus

func OtelGinMiddleware() gin.HandlerFunc {
	return internalmiddlewares.OtelGinMiddleware()
}

func OTLPErrorLogMiddleware(opts ...ErrorLogOption) gin.HandlerFunc {
	return internalmiddlewares.OTLPErrorLogMiddleware(opts...)
}
