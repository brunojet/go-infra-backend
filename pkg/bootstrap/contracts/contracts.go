package contracts

import (
	"context"
	"time"
)

// Shutdown is a generic lifecycle hook for graceful shutdown.
// It matches common components like servers, DB managers, exporters, etc.
type Shutdown interface {
	Shutdown(ctx context.Context) error
}

// ShutdownFunc adapts a function to the Shutdown contract.
type ShutdownFunc func(ctx context.Context) error

// ShutdownManager orchestrates graceful shutdown handlers.
type ShutdownManager interface {
	GetContext() context.Context
	Register(name string, s Shutdown)
	RegisterFunc(name string, fn ShutdownFunc)
	Shutdown() error
	ShutdownWithTimeout(timeout time.Duration) error
}
