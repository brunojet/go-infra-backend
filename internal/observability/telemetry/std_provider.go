package telemetry

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/brunojet/go-infra-backend/internal/observability/contracts"
	"github.com/brunojet/go-infra-backend/internal/observability/requestid"
	"github.com/brunojet/go-infra-backend/internal/observability/types"
)

// StdProvider is a dependency-free telemetry implementation.
//
// It turns spans/metrics into log lines. This is a good "first step" until
// you adopt OpenTelemetry exporters.
//
// It lives in infra and can read request correlation data.
//
// Typical output:
//
//	span start name=apps.Create rid=... fields=...
//	span end   name=apps.Create rid=... latency=... err=...
//	metric inc name=apps.create.ok value=1 rid=...
//
// NOTE: This is intentionally simple and synchronous.
// For high-volume metrics, replace with proper metrics backend.
type StdProvider struct{}

func NewStdProvider() StdProvider { return StdProvider{} }

func (StdProvider) Tracer() contracts.Tracer   { return stdTracer{} }
func (StdProvider) Metrics() contracts.Metrics { return stdMetrics{} }
func (StdProvider) Logger() contracts.Logger   { return stdLogger{} }

type stdTracer struct{}

type stdSpan struct {
	ctx   context.Context
	name  string
	start time.Time
	base  []types.Field
}

type stdMetrics struct{}

type stdLogger struct{}

func (stdTracer) Start(ctx context.Context, name string, fields ...types.Field) (context.Context, contracts.Span) {
	log.Printf("span start name=%s%s", name, formatFields(ctx, fields...))
	return ctx, &stdSpan{ctx: ctx, name: name, start: time.Now(), base: fields}
}

func (s *stdSpan) End(err error, fields ...types.Field) {
	all := append(append([]types.Field{}, s.base...), fields...)
	lat := time.Since(s.start)
	if err != nil {
		log.Printf("span end name=%s latency=%s err=%q%s", s.name, lat, err.Error(), formatFields(s.ctx, all...))
		return
	}
	log.Printf("span end name=%s latency=%s%s", s.name, lat, formatFields(s.ctx, all...))
}

func (stdMetrics) Inc(name string, value int64, fields ...types.Field) {
	log.Printf("metric inc name=%s value=%d%s", name, value, formatFields(context.Background(), fields...))
}

func (stdMetrics) ObserveDuration(name string, d time.Duration, fields ...types.Field) {
	log.Printf("metric observe_duration name=%s value=%s%s", name, d, formatFields(context.Background(), fields...))
}

func (stdLogger) Info(ctx context.Context, msg string, fields ...types.Field) {
	log.Printf("app info msg=%q%s", msg, formatFields(ctx, fields...))
}

func (stdLogger) Error(ctx context.Context, msg string, err error, fields ...types.Field) {
	if err != nil {
		log.Printf("app error msg=%q err=%q%s", msg, err.Error(), formatFields(ctx, fields...))
		return
	}
	log.Printf("app error msg=%q%s", msg, formatFields(ctx, fields...))
}

func formatFields(ctx context.Context, fields ...types.Field) string {
	var b strings.Builder
	if rid, ok := requestid.From(ctx); ok {
		b.WriteString(" rid=")
		b.WriteString(rid)
	}
	if len(fields) == 0 {
		return b.String()
	}
	b.WriteString(" fields=")
	for i, f := range fields {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(f.Key)
		b.WriteString(":")
		b.WriteString(fmt.Sprint(f.Value))
	}
	return b.String()
}
