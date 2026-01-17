package httpmiddleware_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	logcore "github.com/brunojet/go-infra-backend/internal/observability/logging"
	contracts_log "github.com/brunojet/go-infra-backend/internal/observability/logging/contracts"
	metricscore "github.com/brunojet/go-infra-backend/internal/observability/metrics"
	metricsadapter "github.com/brunojet/go-infra-backend/internal/observability/metrics/adapter"
	httpmw "github.com/brunojet/go-infra-backend/internal/observability/middlewares/http"
	tracescore "github.com/brunojet/go-infra-backend/internal/observability/traces"
	contracts_traces "github.com/brunojet/go-infra-backend/internal/observability/traces/contracts"
)

// mock tracer
type mockSpan struct{ ended bool }

func (m *mockSpan) End()                                                    { m.ended = true }
func (m *mockSpan) SetAttribute(key string, value interface{})              {}
func (m *mockSpan) AddEvent(name string, attributes map[string]interface{}) {}
func (m *mockSpan) Context() context.Context                                { return context.Background() }
func (m *mockSpan) SetStatus(code int, message string)                      {}
func (m *mockSpan) RecordError(err error)                                   {}
func (m *mockSpan) SetName(name string)                                     {}

type mockTracer struct{ lastName string }

func (m *mockTracer) Start(ctx context.Context, name string, attributes map[string]interface{}) (context.Context, contracts_traces.Span) {
	m.lastName = name
	return ctx, &mockSpan{}
}
func (m *mockTracer) Inject(ctx context.Context) map[string]string { return nil }
func (m *mockTracer) Extract(ctx context.Context, carrier map[string]string) context.Context {
	return ctx
}

// mock logger
type mockLogger struct{ lastMsg string }

func (m *mockLogger) Debug(ctx context.Context, msg string, fields contracts_log.Fields) {
	m.lastMsg = msg
}
func (m *mockLogger) Info(ctx context.Context, msg string, fields contracts_log.Fields) {
	m.lastMsg = msg
}
func (m *mockLogger) Warn(ctx context.Context, msg string, fields contracts_log.Fields) {
	m.lastMsg = msg
}
func (m *mockLogger) Error(ctx context.Context, msg string, err error, fields contracts_log.Fields) {
	m.lastMsg = msg
}
func (m *mockLogger) WithFields(fields contracts_log.Fields) contracts_log.Logger { return m }
func (m *mockLogger) WithContext(ctx context.Context) contracts_log.Logger        { return m }

func TestHTTPMiddleware_TracingLoggingMetrics(t *testing.T) {
	// tracer
	mt := &mockTracer{}
	tr := tracescore.NewTraces(mt)

	// logger
	ml := &mockLogger{}
	lc := logcore.NewLogging(ml)

	// metrics
	rec := metricsadapter.NewMemRecorder()
	reg := metricsadapter.NewMemRegister(rec)
	mcore := metricscore.NewMetrics(reg, rec)

	rg := gin.New()
	rg.Use(httpmw.TracingMiddleware(tr))
	rg.Use(httpmw.LoggingMiddleware(lc))
	rg.Use(httpmw.MetricsMiddleware(mcore))

	rg.GET("/all", func(c *gin.Context) {
		// logging: retrieve derived logger
		l := httpmw.WithLoggerFromContext(c.Request.Context())
		if l == nil {
			t.Fatalf("expected logger in context")
		}
		l.Info(c.Request.Context(), "hi", nil)
		c.Status(204)
	})

	req := httptest.NewRequest("GET", "/all", nil)
	recorder := httptest.NewRecorder()
	rg.ServeHTTP(recorder, req)

	if mt.lastName == "" {
		t.Fatalf("expected tracer Start called")
	}
	if ml.lastMsg != "hi" {
		t.Fatalf("expected logger called, got %s", ml.lastMsg)
	}
	counters := rec.Counters()
	hist := rec.Histograms()
	if counters["http_requests_total"] != 1 {
		t.Fatalf("expected metrics counter incremented")
	}
	if len(hist["http_request_duration_ms"]) == 0 {
		t.Fatalf("expected histogram observed")
	}
}
