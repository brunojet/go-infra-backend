package contracts

import "context"

type ObservabilityAdapter interface {
	Shutdown(ctx context.Context) error
}

type ObservabilityManager interface {
	RegisterAdapter(adapter ObservabilityAdapter)
	Shutdown(ctx context.Context) error
}
