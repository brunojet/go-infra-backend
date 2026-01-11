package nethttp

import (
	"net/http"

	"github.com/brunojet/go-infra-backend/internal/http/contracts"
	"github.com/go-chi/chi/v5"
)

type Option func(*runtimeOptions)

type runtimeOptions struct {
	middlewares []func(http.Handler) http.Handler
}

func defaultRuntimeOptions() runtimeOptions {
	return runtimeOptions{}
}

// WithMiddlewares appends custom net/http middlewares.
func WithMiddlewares(mw ...func(http.Handler) http.Handler) Option {
	return func(o *runtimeOptions) {
		o.middlewares = append(o.middlewares, mw...)
	}
}

// NewRuntime builds a net/http (chi) runtime and returns framework-agnostic Router + net/http handler.
//
// Middlewares are configured via options from the host/bootstrap layer.
func NewRuntime(opts ...Option) (contracts.Router, http.Handler) {
	o := defaultRuntimeOptions()
	for _, opt := range opts {
		opt(&o)
	}

	r := chi.NewRouter()
	for _, mw := range o.middlewares {
		r.Use(mw)
	}
	return NewChiRouter(r), r
}
