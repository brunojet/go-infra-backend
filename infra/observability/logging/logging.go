package logging

import (
	"context"
	"sync"

	contracts "github.com/brunojet/go-infra-backend/infra/observability/logging/contracts"
)

// Logging is a thin, thread-safe facade around a `contracts.Logger` adapter.
// It implements `contracts.Logger` and provides helpers `WithFields`,
// `WithContext` and `SetLogger` to swap the underlying adapter at runtime.
type Logging struct {
	mu     sync.RWMutex
	logger contracts.Logger
}

// NewLogging constructs a new logging core backed by the provided adapter.
func NewLogging(logger contracts.Logger) *Logging {
	return &Logging{logger: logger}
}

func (c *Logging) getLogger() contracts.Logger {
	c.mu.RLock()
	l := c.logger
	c.mu.RUnlock()
	return l
}

// SetLogger replaces the underlying logger adapter.
func (c *Logging) SetLogger(logger contracts.Logger) {
	c.mu.Lock()
	c.logger = logger
	c.mu.Unlock()
}

// Debug/Info/Warn/Error delegate to the underlying logger if present.
func (c *Logging) Debug(ctx context.Context, msg string, fields contracts.Fields) {
	if l := c.getLogger(); l != nil {
		l.Debug(ctx, msg, fields)
	}
}
func (c *Logging) Info(ctx context.Context, msg string, fields contracts.Fields) {
	if l := c.getLogger(); l != nil {
		l.Info(ctx, msg, fields)
	}
}
func (c *Logging) Warn(ctx context.Context, msg string, fields contracts.Fields) {
	if l := c.getLogger(); l != nil {
		l.Warn(ctx, msg, fields)
	}
}
func (c *Logging) Error(ctx context.Context, msg string, err error, fields contracts.Fields) {
	if l := c.getLogger(); l != nil {
		l.Error(ctx, msg, err, fields)
	}
}

// derivedLogger implements contracts.Logger and carries optional pre-set
// context and fields that are merged with call-time values.
type derivedLogger struct {
	base   *Logging
	ctx    context.Context
	fields contracts.Fields
}

func (d *derivedLogger) effectiveCtx(callCtx context.Context) context.Context {
	if d.ctx != nil {
		return d.ctx
	}
	return callCtx
}

func (d *derivedLogger) mergeFields(callFields contracts.Fields) contracts.Fields {
	if d.fields == nil && callFields == nil {
		return nil
	}
	out := make(contracts.Fields)
	for k, v := range d.fields {
		out[k] = v
	}
	for k, v := range callFields {
		out[k] = v
	}
	return out
}

func (d *derivedLogger) Debug(ctx context.Context, msg string, fields contracts.Fields) {
	if l := d.base.getLogger(); l != nil {
		l.Debug(d.effectiveCtx(ctx), msg, d.mergeFields(fields))
	}
}
func (d *derivedLogger) Info(ctx context.Context, msg string, fields contracts.Fields) {
	if l := d.base.getLogger(); l != nil {
		l.Info(d.effectiveCtx(ctx), msg, d.mergeFields(fields))
	}
}
func (d *derivedLogger) Warn(ctx context.Context, msg string, fields contracts.Fields) {
	if l := d.base.getLogger(); l != nil {
		l.Warn(d.effectiveCtx(ctx), msg, d.mergeFields(fields))
	}
}
func (d *derivedLogger) Error(ctx context.Context, msg string, err error, fields contracts.Fields) {
	if l := d.base.getLogger(); l != nil {
		l.Error(d.effectiveCtx(ctx), msg, err, d.mergeFields(fields))
	}
}

func (d *derivedLogger) WithFields(fields contracts.Fields) contracts.Logger {
	merged := make(contracts.Fields)
	for k, v := range d.fields {
		merged[k] = v
	}
	for k, v := range fields {
		merged[k] = v
	}
	return &derivedLogger{base: d.base, ctx: d.ctx, fields: merged}
}

func (d *derivedLogger) WithContext(ctx context.Context) contracts.Logger {
	return &derivedLogger{base: d.base, ctx: ctx, fields: d.fields}
}

// WithFields returns a derived logger that will include the provided fields
// on every entry.
func (c *Logging) WithFields(fields contracts.Fields) contracts.Logger {
	return &derivedLogger{base: c, fields: fields}
}

// WithContext returns a derived logger that will use the provided context for
// all log entries (useful to avoid passing ctx on every call).
func (c *Logging) WithContext(ctx context.Context) contracts.Logger {
	return &derivedLogger{base: c, ctx: ctx}
}
