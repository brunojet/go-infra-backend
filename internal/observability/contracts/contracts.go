package contracts

import "context"

type ShutdownFunc func(ctx context.Context) error

type ShutdownRegister interface {
	RegisterFunc(name string, fn ShutdownFunc)
}
