package adapter_test

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"

	ada "github.com/brunojet/go-infra-backend/internal/observability/traces/adapter"
)

func TestOTelTracer_StartInjectExtract(t *testing.T) {
	// set up in-memory tracer provider to capture spans
	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	tr := ada.NewOTelTracer("test-tracer")
	ctx, sp := tr.Start(context.Background(), "op", map[string]interface{}{"k": "v"})
	if sp == nil {
		t.Fatal("expected non-nil span")
	}
	// Inject should produce a carrier with traceparent
	carr := tr.Inject(ctx)
	if len(carr) == 0 {
		t.Fatal("expected inject to add entries to carrier")
	}

	// Extract should return a context that yields a valid span context
	extracted := tr.Extract(context.Background(), carr)
	_, span2 := otel.Tracer("test-tracer").Start(extracted, "inner")
	// close started span
	span2.End()

	// End the original span to record it
	sp.End()

	if len(sr.Ended()) == 0 {
		t.Fatalf("expected recorded ended spans, got 0")
	}
}
