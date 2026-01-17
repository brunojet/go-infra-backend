package adapter

import (
	"context"

	contracts "github.com/brunojet/go-infra-backend/infra/observability/traces/contracts"
)

// NoopSpan is a trivial Span implementation that does nothing.
type NoopSpan struct{}

func (NoopSpan) End()                                                    {}
func (NoopSpan) SetAttribute(key string, value interface{})              {}
func (NoopSpan) AddEvent(name string, attributes map[string]interface{}) {}
func (NoopSpan) Context() context.Context                                { return context.Background() }
func (NoopSpan) SetStatus(code int, message string)                      {}
func (NoopSpan) RecordError(err error)                                   {}
func (NoopSpan) SetName(name string)                                     {}

// NoopTracer is a lightweight tracer that satisfies contracts.Tracer but does nothing.
type NoopTracer struct{}

func NewNoopTracer() *NoopTracer { return &NoopTracer{} }

func (n *NoopTracer) Start(ctx context.Context, name string, attributes map[string]interface{}) (context.Context, contracts.Span) {
	return ctx, nil
}

func (n *NoopTracer) Inject(ctx context.Context) map[string]string { return nil }

func (n *NoopTracer) Extract(ctx context.Context, carrier map[string]string) context.Context {
	return ctx
}
