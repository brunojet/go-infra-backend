package ginadapter

import (
	"net/http"

	"github.com/brunojet/go-infra-backend/internal/http/contracts"
	"github.com/gin-gonic/gin"
)

type Option func(*runtimeOptions)

type runtimeOptions struct {
	middlewares []gin.HandlerFunc
}

func defaultRuntimeOptions() runtimeOptions {
	return runtimeOptions{
		middlewares: nil,
	}
}

// WithMiddlewares appends custom Gin middlewares.
func WithMiddlewares(mw ...gin.HandlerFunc) Option {
	return func(o *runtimeOptions) {
		o.middlewares = append(o.middlewares, mw...)
	}
}

// NewRuntime builds a Gin HTTP runtime and returns framework-agnostic Router + net/http handler.
//
// Middlewares are configured via options from the host/bootstrap layer.
func NewRuntime(opts ...Option) (contracts.Router, http.Handler) {
	o := defaultRuntimeOptions()
	for _, opt := range opts {
		opt(&o)
	}

	engine := gin.New()
	if len(o.middlewares) > 0 {
		engine.Use(o.middlewares...)
	}

	rootGroup := engine.Group("")
	return NewGinRouter(rootGroup), engine
}
