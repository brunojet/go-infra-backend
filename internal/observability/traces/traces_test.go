package traces_test

import (
	"context"
	"testing"

	"github.com/brunojet/go-infra-backend/internal/observability/traces"
	"github.com/brunojet/go-infra-backend/internal/observability/traces/contracts"
)

type mockSpan struct {
	ended bool
}

func (s *mockSpan) End()                                                    { s.ended = true }
func (s *mockSpan) SetAttribute(key string, value interface{})              {}
func (s *mockSpan) AddEvent(name string, attributes map[string]interface{}) {}
func (s *mockSpan) Context() context.Context                                { return context.Background() }
func (s *mockSpan) SetStatus(code int, message string)                      {}
func (s *mockSpan) RecordError(err error)                                   {}
func (s *mockSpan) SetName(name string)                                     {}

type mockTracer struct {
	lastCtx   context.Context
	lastName  string
	lastAttrs map[string]interface{}
	inject    map[string]string
}

func (m *mockTracer) Start(ctx context.Context, name string, attributes map[string]interface{}) (context.Context, contracts.Span) {
	m.lastCtx = ctx
	m.lastName = name
	m.lastAttrs = attributes
	return ctx, &mockSpan{}
}
func (m *mockTracer) Inject(ctx context.Context) map[string]string {
	return m.inject
}
func (m *mockTracer) Extract(ctx context.Context, carrier map[string]string) context.Context {
	return ctx
}

func TestTracesCore_DelegationAndSwap(t *testing.T) {
	m1 := &mockTracer{inject: map[string]string{"x": "1"}}
	c := traces.NewTraces(m1)

	ctx, sp := c.Start(context.Background(), "op", map[string]interface{}{"k": "v"})
	if sp == nil {
		t.Fatalf("expected non-nil span")
	}
	if m1.lastName != "op" {
		t.Fatalf("expected start delegated to m1")
	}

	if got := c.Inject(ctx); got["x"] != "1" {
		t.Fatalf("expected inject forwarded")
	}

	// swap tracer
	m2 := &mockTracer{inject: map[string]string{"y": "2"}}
	c.SetTracer(m2)
	ctx2, _ := c.Start(context.Background(), "op2", nil)
	if m2.lastName != "op2" {
		t.Fatalf("expected start delegated to m2")
	}
	if got := c.Inject(ctx2); got["y"] != "2" {
		t.Fatalf("expected inject forwarded to m2")
	}
}

func TestDerivedTracer_WithContext(t *testing.T) {
	m := &mockTracer{}
	c := traces.NewTraces(m)

	preset := context.WithValue(context.Background(), "trace", "t1")
	d := c.WithContext(preset)
	d.Start(context.Background(), "nm", nil)
	if m.lastCtx == nil || m.lastCtx.Value("trace") != "t1" {
		t.Fatalf("expected preset context to be used by derived tracer")
	}
}

func TestNilTracer_Behavior(t *testing.T) {
	c := traces.NewTraces(nil)
	ctx, sp := c.Start(context.Background(), "no", nil)
	if sp != nil {
		t.Fatalf("expected nil span when tracer is nil")
	}
	if got := c.Inject(ctx); got != nil {
		t.Fatalf("expected nil inject when tracer is nil")
	}
}
