package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	otellog "go.opentelemetry.io/otel/log"
	otellogglobal "go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/trace"
)

type ErrorLogOption func(*errorLogConfig)

type errorLogConfig struct {
	minStatus int
}

func WithMinStatus(minStatus int) ErrorLogOption {
	return func(cfg *errorLogConfig) {
		cfg.minStatus = minStatus
	}
}

func severityForHTTPStatus(status int) otellog.Severity {
	if status >= http.StatusInternalServerError {
		return otellog.SeverityError
	}
	if status >= http.StatusBadRequest {
		return otellog.SeverityWarn
	}
	return otellog.SeverityInfo
}

// OTLPErrorLogMiddleware emits an OpenTelemetry LogRecord when a request finishes
// with a status code >= minStatus (default: 500) OR when gin collected errors.
//
// Important: register this middleware AFTER OtelGinMiddleware(), so it runs
// while the request context still contains the active span.
func OTLPErrorLogMiddleware(opts ...ErrorLogOption) gin.HandlerFunc {
	cfg := errorLogConfig{minStatus: http.StatusInternalServerError}
	for _, opt := range opts {
		opt(&cfg)
	}

	logger := otellogglobal.Logger("internal/observability/http_middlewares/OTLPErrorLogMiddleware")

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := c.Writer.Status()
		if len(c.Errors) == 0 && status < cfg.minStatus {
			return
		}

		ctx := c.Request.Context()
		now := time.Now()

		var record otellog.Record
		record.SetTimestamp(now)
		record.SetObservedTimestamp(now)
		record.SetSeverity(severityForHTTPStatus(status))

		// Body: prefer gin errors when available.
		if len(c.Errors) > 0 {
			record.SetBody(otellog.StringValue(c.Errors.String()))
		} else {
			record.SetBody(otellog.StringValue("http request failed"))
		}

		attrs := []otellog.KeyValue{
			otellog.Int("http.status_code", status),
			otellog.String("http.method", c.Request.Method),
			otellog.String("http.route", c.FullPath()),
			otellog.String("http.target", c.Request.URL.Path),
			otellog.String("http.user_agent", c.Request.UserAgent()),
			otellog.String("net.peer.ip", c.ClientIP()),
			otellog.Int64("http.server.duration_ms", time.Since(start).Milliseconds()),
			otellog.Int("gin.errors.count", len(c.Errors)),
		}

		// Optional explicit IDs for backends that don't surface the SDK trace linkage.
		if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
			attrs = append(attrs,
				otellog.String("trace_id", sc.TraceID().String()),
				otellog.String("span_id", sc.SpanID().String()),
			)
		}

		record.AddAttributes(attrs...)
		logger.Emit(ctx, record)
	}
}
