package adapters

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	bffcts "github.com/brunojet/go-infra-backend/pkg/bffclient/contracts"
)

// tracerName is the instrumentation scope name used when acquiring a Tracer
// from the global TracerProvider (registered by the observability bootstrap).
const tracerName = "github.com/brunojet/go-infra-backend/bffclient"

type otelMiddleware struct {
	tracer trace.Tracer
}

// NewOtelMiddleware returns a BffClientMiddleware that creates a child OTel
// span for each outgoing upstream call.
//
// It uses otel.Tracer(tracerName) which resolves to whatever TracerProvider
// was registered by the observability bootstrap — no extra configuration needed.
func NewOtelMiddleware() bffcts.BffClientMiddleware {
	return &otelMiddleware{tracer: otel.Tracer(tracerName)}
}

// OnRequest starts a client-kind span whose name is "{METHOD} {path}".
// The span is stored in the returned context so OnResponse can retrieve and
// end it. The adapter will inject the W3C traceparent header separately via
// propagation.TraceContext{}.Inject.
func (m *otelMiddleware) OnRequest(ctx context.Context, req bffcts.BffRequestInfo) context.Context {
	ctx, _ = m.tracer.Start(ctx, req.Method+" "+req.Path,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("http.method", req.Method),
			attribute.String("http.url", req.Path),
		),
	)
	return ctx
}

// OnResponse records the HTTP status code on the span, marks it as error if
// applicable, and ends it. This is always called — even on transport errors.
func (m *otelMiddleware) OnResponse(ctx context.Context, _ bffcts.BffRequestInfo, resp bffcts.BffResponseInfo) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	if resp.StatusCode != 0 {
		span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))
	}

	if resp.Err != nil {
		span.RecordError(resp.Err)
		span.SetStatus(codes.Error, resp.Err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}
}
