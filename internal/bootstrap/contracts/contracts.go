package contracts

import (
	"context"
	"time"
)

type Shutdown interface {
	Shutdown(ctx context.Context) error
}

type ShutdownFunc func(ctx context.Context) error

type ShutdownManager interface {
	GetContext() context.Context
	Register(name string, s Shutdown)
	RegisterFunc(name string, fn ShutdownFunc)
	Shutdown() error
	ShutdownWithTimeout(timeout time.Duration) error
}
