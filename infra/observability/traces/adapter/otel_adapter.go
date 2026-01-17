package adapter

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"

	contracts "github.com/brunojet/go-infra-backend/infra/observability/traces/contracts"
)

// OTelSpan wraps an OpenTelemetry span into contracts.Span.
type OTelSpan struct {
	span oteltrace.Span
	ctx  context.Context
}

func (s *OTelSpan) End() {
	if s.span != nil {
		s.span.End()
	}
}
func (s *OTelSpan) SetAttribute(key string, value interface{}) {
	if s.span == nil {
		return
	}
	s.span.SetAttributes(attribute.String(key, fmt.Sprint(value)))
}
func (s *OTelSpan) AddEvent(name string, attributes map[string]interface{}) {
	if s.span == nil {
		return
	}
	attrs := make([]attribute.KeyValue, 0, len(attributes))
	for k, v := range attributes {
		attrs = append(attrs, attribute.String(k, fmt.Sprint(v)))
	}
	s.span.AddEvent(name, oteltrace.WithAttributes(attrs...))
}
func (s *OTelSpan) Context() context.Context {
	if s.ctx == nil {
		return context.Background()
	}
	return s.ctx
}
func (s *OTelSpan) SetStatus(code int, message string) {
	if s.span == nil {
		return
	}
	if code == 0 {
		s.span.SetStatus(codes.Ok, message)
	} else {
		s.span.SetStatus(codes.Error, message)
	}
}
func (s *OTelSpan) RecordError(err error) {
	if s.span == nil {
		return
	}
	s.span.RecordError(err)
}
func (s *OTelSpan) SetName(name string) {
	if s.span == nil {
		return
	}
	s.span.SetName(name)
}

// OTelTracer implements contracts.Tracer using the global OpenTelemetry Tracer.
type OTelTracer struct {
	tracer     oteltrace.Tracer
	propagator propagation.TextMapPropagator
}

// NewOTelTracer returns a tracer adapter. It uses the global otel TracerProvider when tracerName empty.
func NewOTelTracer(tracerName string) *OTelTracer {
	name := tracerName
	if name == "" {
		name = "go-infra-backend"
	}
	return &OTelTracer{tracer: otel.Tracer(name), propagator: otel.GetTextMapPropagator()}
}

func (t *OTelTracer) Start(ctx context.Context, name string, attributes map[string]interface{}) (context.Context, contracts.Span) {
	attrs := make([]attribute.KeyValue, 0, len(attributes))
	for k, v := range attributes {
		attrs = append(attrs, attribute.String(k, fmt.Sprint(v)))
	}
	ctx, span := t.tracer.Start(ctx, name, oteltrace.WithAttributes(attrs...))
	if span == nil {
		return ctx, nil
	}
	return ctx, &OTelSpan{span: span}
}

func (t *OTelTracer) Inject(ctx context.Context) map[string]string {
	carrier := propagation.MapCarrier{}
	t.propagator.Inject(ctx, carrier)
	out := map[string]string{}
	for k, v := range carrier {
		out[k] = v
	}
	return out
}

func (t *OTelTracer) Extract(ctx context.Context, carrier map[string]string) context.Context {
	m := propagation.MapCarrier(carrier)
	return t.propagator.Extract(ctx, m)
}
