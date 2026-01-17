package contracts_test

import (
	"context"
	"errors"
	"testing"

	mm "github.com/brunojet/go-infra-backend/infra/observability/traces/contracts"
)

type mockSpan struct{}

func (s *mockSpan) End()                                                    {}
func (s *mockSpan) SetAttribute(key string, value interface{})              {}
func (s *mockSpan) AddEvent(name string, attributes map[string]interface{}) {}
func (s *mockSpan) Context() context.Context                                { return context.Background() }
func (s *mockSpan) SetStatus(code int, message string)                      {}
func (s *mockSpan) RecordError(err error)                                   {}
func (s *mockSpan) SetName(name string)                                     {}

type mockTracer struct{}

func (t *mockTracer) Start(ctx context.Context, name string, attributes map[string]interface{}) (context.Context, mm.Span) {
	return ctx, &mockSpan{}
}
func (t *mockTracer) Inject(ctx context.Context) map[string]string {
	return map[string]string{"x": "1"}
}
func (t *mockTracer) Extract(ctx context.Context, carrier map[string]string) context.Context {
	return ctx
}

func TestTracerInterface(t *testing.T) {
	var _ mm.Tracer = (*mockTracer)(nil)
	// exercise RecordError to keep linter happy
	s := &mockSpan{}
	s.RecordError(errors.New("x"))
	s.SetStatus(1, "ok")
}
