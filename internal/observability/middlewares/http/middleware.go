package httpmiddleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"

	logcore "github.com/brunojet/go-infra-backend/internal/observability/logging"
	contracts_log "github.com/brunojet/go-infra-backend/internal/observability/logging/contracts"
	metricscore "github.com/brunojet/go-infra-backend/internal/observability/metrics"
	contracts_metrics "github.com/brunojet/go-infra-backend/internal/observability/metrics/contracts"
	tracescore "github.com/brunojet/go-infra-backend/internal/observability/traces"
)

// Logging context key
const ctxLoggerKey = "logger"

func WithLoggerFromContext(ctx context.Context) contracts_log.Logger {
	if ctx == nil {
		return nil
	}
	if v := ctx.Value(ctxLoggerKey); v != nil {
		if l, ok := v.(contracts_log.Logger); ok {
			return l
		}
	}
	return nil
}

// TracingMiddleware starts a span for each HTTP request using the provided traces core.
// It extracts incoming propagation headers and injects the resulting context into the request.
func TracingMiddleware(tr *tracescore.Traces) gin.HandlerFunc {
	return func(c *gin.Context) {
		// extract carrier from headers
		carrier := map[string]string{}
		for k, v := range c.Request.Header {
			if len(v) > 0 {
				carrier[k] = v[0]
			}
		}
		ctx := tr.Extract(c.Request.Context(), carrier)

		// If an upstream middleware (e.g. otelgin) already created a span,
		// avoid starting a new one here to prevent duplicate/overlapping spans.
		existing := trace.SpanFromContext(ctx)
		if existing != nil && existing.SpanContext().IsValid() {
			// attach ctx and continue — enrichment/closing is handled by the starter
			c.Request = c.Request.WithContext(ctx)
			c.Next()
			return
		}

		// no existing span: start one using traces core
		name := c.Request.Method + " " + c.FullPath()
		if name == " " { // fallback if route not matched
			name = c.Request.Method + " " + c.Request.URL.Path
		}
		ctx, span := tr.Start(ctx, name, map[string]interface{}{
			"http.method": c.Request.Method,
			"http.route":  c.FullPath(),
			"http.path":   c.Request.URL.Path,
		})
		// attach ctx to request
		c.Request = c.Request.WithContext(ctx)
		c.Next()
		status := c.Writer.Status()
		if span != nil {
			span.SetAttribute("http.status_code", status)
			if status >= 500 {
				span.SetStatus(1, "server error")
			}
			span.End()
		}
	}
}

// NewMiddlewareChain builds the ordered middleware slice for the application.
//   - adapter: optional tracing adapter (e.g. adapters.OtelGinMiddleware). When
//     provided it is placed first so it can act as the span starter.
//   - tr: tracing core used as fallback starter when no adapter created the span.
//   - l: logging core; when non-nil LoggingMiddleware is added.
//   - m: metrics core; when non-nil MetricsMiddleware is added.
func NewMiddlewareChain(adapter gin.HandlerFunc, tr *tracescore.Traces, l *logcore.Logging, m *metricscore.Metrics) []gin.HandlerFunc {
	chain := make([]gin.HandlerFunc, 0, 6)
	if adapter != nil {
		chain = append(chain, adapter)
	}
	// tracing (adapter may have started span; TracingMiddleware will detect it)
	if tr != nil {
		chain = append(chain, TracingMiddleware(tr))
	}
	// enrich span with application attributes
	chain = append(chain, EnrichAppSpan())
	// logging
	if l != nil {
		chain = append(chain, LoggingMiddleware(l))
	}
	// metrics
	if m != nil {
		chain = append(chain, MetricsMiddleware(m))
	}
	return chain
}

// MiddlewaresWithTracing returns an ordered middleware slice that includes the
// provided tracing adapter (e.g. `adapters.OtelGinMiddleware`) followed by the
// package `TracingMiddleware`. When adapter is nil only the custom
// `TracingMiddleware(tr)` is returned. The `TracingMiddleware` itself will
// detect an already-started span and avoid creating a duplicate.
func MiddlewaresWithTracing(adapter gin.HandlerFunc, tr *tracescore.Traces) []gin.HandlerFunc {
	if adapter == nil {
		return []gin.HandlerFunc{TracingMiddleware(tr)}
	}
	return []gin.HandlerFunc{adapter, TracingMiddleware(tr)}
}

// LoggingMiddleware derives a logger with request fields and attaches it to context via WithContext/WithFields.
// The derived logger is stored in context value under key "logger" to be retrieved by handlers.
func LoggingMiddleware(l *logcore.Logging) gin.HandlerFunc {
	return func(c *gin.Context) {
		fields := contracts_log.Fields{
			"http.method": c.Request.Method,
			"http.path":   c.Request.URL.Path,
		}
		derived := l.WithFields(fields).WithContext(c.Request.Context())
		// store in context
		ctx := context.WithValue(c.Request.Context(), ctxLoggerKey, derived)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// MetricsMiddleware records request count and latency.
func MetricsMiddleware(m *metricscore.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		status := c.Writer.Status()
		m.CounterAdd(c.Request.Context(), "http_requests_total", 1, contracts_metrics.Labels{"method": c.Request.Method, "path": c.FullPath(), "status": http.StatusText(status)})
		elapsed := float64(time.Since(start).Milliseconds())
		m.HistogramObserve(c.Request.Context(), "http_request_duration_ms", elapsed, contracts_metrics.Labels{"method": c.Request.Method, "path": c.FullPath()})
	}
}
