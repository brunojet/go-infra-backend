package observability

import (
	"context"
	"errors"
)

type ObservabilityAdapter interface {
	Shutdown(ctx context.Context) error
}

type ObservabilityManager interface {
	RegisterAdapter(adapter ObservabilityAdapter)
	Shutdown(ctx context.Context) error
}

type observabilityManager struct {
	adapters []ObservabilityAdapter
}

func NewObservabilityManager() *observabilityManager {
	return &observabilityManager{}
}

func (m *observabilityManager) RegisterAdapter(adapter ObservabilityAdapter) {
	if nil == adapter {
		return
	}
	m.adapters = append(m.adapters, adapter)
}

func (m *observabilityManager) Shutdown(ctx context.Context) error {
	var errs error
	for i := len(m.adapters) - 1; i >= 0; i-- {
		err := m.adapters[i].Shutdown(ctx)
		errs = errors.Join(errs, err)
	}
	return errs
}
