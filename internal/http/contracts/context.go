package contracts

import "context"

// Context is a minimal, framework-agnostic HTTP request context.
// It intentionally models only what our handlers need today.
//
// The goal is to isolate any web-framework dependency in adapters.
type Context interface {
	RequestContext() context.Context

	Param(name string) string
	BindJSON(dst any) error

	JSON(status int, v any)
	String(status int, v string)
	Status(status int)
}

type HandlerFunc func(Context)
