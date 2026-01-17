package traces

import (
	"context"
	"sync"

	contracts "github.com/brunojet/go-infra-backend/infra/observability/traces/contracts"
)

// Traces is a thin, thread-safe facade around a `contracts.Tracer` adapter.
// It provides helpers to start spans and to swap the underlying tracer at runtime.
type Traces struct {
	mu     sync.RWMutex
	tracer contracts.Tracer
}

// NewTraces constructs a new traces core backed by the provided tracer.
func NewTraces(tracer contracts.Tracer) *Traces {
	return &Traces{tracer: tracer}
}

func (c *Traces) getTracer() contracts.Tracer {
	c.mu.RLock()
	t := c.tracer
	c.mu.RUnlock()
	return t
}

// SetTracer replaces the underlying tracer adapter.
func (c *Traces) SetTracer(tracer contracts.Tracer) {
	c.mu.Lock()
	c.tracer = tracer
	c.mu.Unlock()
}

// Start delegates to the underlying tracer if present. If no tracer is set
// it returns the provided context and a nil Span.
func (c *Traces) Start(ctx context.Context, name string, attributes map[string]interface{}) (context.Context, contracts.Span) {
	if t := c.getTracer(); t != nil {
		return t.Start(ctx, name, attributes)
	}
	return ctx, nil
}

// Inject delegates to the underlying tracer or returns nil when no tracer is set.
func (c *Traces) Inject(ctx context.Context) map[string]string {
	if t := c.getTracer(); t != nil {
		return t.Inject(ctx)
	}
	return nil
}

// Extract delegates to the underlying tracer or returns the provided context
// unchanged when no tracer is set.
func (c *Traces) Extract(ctx context.Context, carrier map[string]string) context.Context {
	if t := c.getTracer(); t != nil {
		return t.Extract(ctx, carrier)
	}
	return ctx
}

// derivedTracer implements contracts.Tracer and carries an optional preset
// context that will be used as the effective context for Start calls.
type derivedTracer struct {
	base *Traces
	ctx  context.Context
}

func (d *derivedTracer) Start(ctx context.Context, name string, attributes map[string]interface{}) (context.Context, contracts.Span) {
	callCtx := ctx
	if d.ctx != nil {
		callCtx = d.ctx
	}
	if t := d.base.getTracer(); t != nil {
		return t.Start(callCtx, name, attributes)
	}
	return callCtx, nil
}

func (d *derivedTracer) Inject(ctx context.Context) map[string]string {
	if t := d.base.getTracer(); t != nil {
		return t.Inject(ctx)
	}
	return nil
}

func (d *derivedTracer) Extract(ctx context.Context, carrier map[string]string) context.Context {
	if t := d.base.getTracer(); t != nil {
		return t.Extract(ctx, carrier)
	}
	return ctx
}

// WithContext returns a derived tracer that will use the provided context for
// all Start calls.
func (c *Traces) WithContext(ctx context.Context) contracts.Tracer {
	return &derivedTracer{base: c, ctx: ctx}
}
