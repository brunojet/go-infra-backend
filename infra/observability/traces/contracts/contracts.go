package contracts

import "context"

// Span represents an active trace span with a minimal, implementation-agnostic
// API to be used by core code and adapters.
type Span interface {
	End()
	SetAttribute(key string, value interface{})
	AddEvent(name string, attributes map[string]interface{})
	// Context returns the context containing the span for propagation.
	Context() context.Context
	// SetStatus sets the span status code and message (useful to mark errors).
	SetStatus(code int, message string)
	// RecordError records an error on the span.
	RecordError(err error)
	// SetName renames the span if needed.
	SetName(name string)
}

// Tracer abstracts starting spans and context propagation in an agnostic way.
type Tracer interface {
	// Start begins a span and returns a context containing that span plus the Span handle.
	Start(ctx context.Context, name string, attributes map[string]interface{}) (context.Context, Span)

	// Inject converts the trace context into a carrier map for propagation (HTTP headers, etc.).
	Inject(ctx context.Context) map[string]string

	// Extract obtains a context with the extracted trace data from a carrier map.
	Extract(ctx context.Context, carrier map[string]string) context.Context
}
