package adapter

import (
	"context"

	"github.com/rs/zerolog"

	contracts "github.com/brunojet/go-infra-backend/internal/observability/logging/contracts"
)

// ZerologAdapter adapts zerolog to the contracts.Logger interface.
type ZerologAdapter struct {
	logger  zerolog.Logger
	enabled bool
}

// NewZerologAdapter constructs a new ZerologAdapter.
func NewZerologAdapter(l zerolog.Logger) *ZerologAdapter {
	return &ZerologAdapter{logger: l, enabled: true}
}

func (z *ZerologAdapter) extractTraceFields(ctx context.Context) contracts.Fields {
	if ctx == nil {
		return nil
	}
	out := make(contracts.Fields)
	if v := ctx.Value("trace"); v != nil {
		out["trace"] = v
	}
	if v := ctx.Value("span"); v != nil {
		out["span"] = v
	}
	if v := ctx.Value("trace_id"); v != nil {
		out["trace_id"] = v
	}
	if v := ctx.Value("span_id"); v != nil {
		out["span_id"] = v
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func mergeFields(a, b contracts.Fields) contracts.Fields {
	if a == nil && b == nil {
		return nil
	}
	out := make(contracts.Fields)
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		out[k] = v
	}
	return out
}

func (z *ZerologAdapter) logWithLevel(ctx context.Context, level string, msg string, err error, fields contracts.Fields) {
	if !z.enabled {
		return
	}
	ev := z.logger.With().Logger()
	var e *zerolog.Event
	switch level {
	case "debug":
		e = ev.Debug()
	case "info":
		e = ev.Info()
	case "warn":
		e = ev.Warn()
	default:
		e = ev.Error()
	}

	// attach trace/span fields from ctx
	if tf := z.extractTraceFields(ctx); tf != nil {
		for k, v := range tf {
			e.Interface(k, v)
		}
	}

	// attach provided fields
	for k, v := range fields {
		e.Interface(k, v)
	}

	if err != nil {
		e.Err(err).Msg(msg)
		return
	}
	e.Msg(msg)
}

func (z *ZerologAdapter) Debug(ctx context.Context, msg string, fields contracts.Fields) {
	z.logWithLevel(ctx, "debug", msg, nil, fields)
}
func (z *ZerologAdapter) Info(ctx context.Context, msg string, fields contracts.Fields) {
	z.logWithLevel(ctx, "info", msg, nil, fields)
}
func (z *ZerologAdapter) Warn(ctx context.Context, msg string, fields contracts.Fields) {
	z.logWithLevel(ctx, "warn", msg, nil, fields)
}
func (z *ZerologAdapter) Error(ctx context.Context, msg string, err error, fields contracts.Fields) {
	z.logWithLevel(ctx, "error", msg, err, fields)
}

// derived wrapper to satisfy WithFields/WithContext returning contracts.Logger
type zerologDerived struct {
	base   *ZerologAdapter
	ctx    context.Context
	fields contracts.Fields
}

func (d *zerologDerived) effectiveCtx(callCtx context.Context) context.Context {
	if d.ctx != nil {
		return d.ctx
	}
	return callCtx
}

func (d *zerologDerived) mergedFields(callFields contracts.Fields) contracts.Fields {
	return mergeFields(d.fields, callFields)
}

func (d *zerologDerived) Debug(ctx context.Context, msg string, fields contracts.Fields) {
	d.base.Debug(d.effectiveCtx(ctx), msg, d.mergedFields(fields))
}
func (d *zerologDerived) Info(ctx context.Context, msg string, fields contracts.Fields) {
	d.base.Info(d.effectiveCtx(ctx), msg, d.mergedFields(fields))
}
func (d *zerologDerived) Warn(ctx context.Context, msg string, fields contracts.Fields) {
	d.base.Warn(d.effectiveCtx(ctx), msg, d.mergedFields(fields))
}
func (d *zerologDerived) Error(ctx context.Context, msg string, err error, fields contracts.Fields) {
	d.base.Error(d.effectiveCtx(ctx), msg, err, d.mergedFields(fields))
}
func (d *zerologDerived) WithFields(fields contracts.Fields) contracts.Logger {
	return &zerologDerived{base: d.base, ctx: d.ctx, fields: mergeFields(d.fields, fields)}
}
func (d *zerologDerived) WithContext(ctx context.Context) contracts.Logger {
	return &zerologDerived{base: d.base, ctx: ctx, fields: d.fields}
}

func (z *ZerologAdapter) WithFields(fields contracts.Fields) contracts.Logger {
	return &zerologDerived{base: z, fields: fields}
}
func (z *ZerologAdapter) WithContext(ctx context.Context) contracts.Logger {
	return &zerologDerived{base: z, ctx: ctx}
}
