package adapters

import (
	"bytes"
	"context"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	otellog "go.opentelemetry.io/otel/log"
	otellogglobal "go.opentelemetry.io/otel/log/global"
)

const (
	initialStdLogBufferCap = 16 * 1024
	maxStdLogBufferCap     = 64 * 1024
	truncatedLogSuffix     = " [truncated]"
)

type otelStdLogWriter struct {
	loggerName string
	severity   otellog.Severity
	logger     otellog.Logger
	ctx        context.Context

	mu  sync.Mutex
	buf bytes.Buffer
}

func (w *otelStdLogWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Buffer writes because the stdlib logger may write in chunks.
	_, _ = w.buf.Write(p)

	for {
		b := w.buf.Bytes()
		idx := bytes.IndexByte(b, '\n')
		if idx < 0 {
			w.flushOversizedBufferLocked()
			return len(p), nil
		}

		line := string(b[:idx])
		// Drop the processed line + newline.
		w.buf.Next(idx + 1)

		line = strings.TrimSuffix(line, "\r")
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		w.emitLine(line)
	}
}

func (w *otelStdLogWriter) flushOversizedBufferLocked() {
	for w.buf.Len() > maxStdLogBufferCap {
		chunk := string(w.buf.Bytes()[:maxStdLogBufferCap])
		w.buf.Next(maxStdLogBufferCap)

		chunk = strings.TrimSuffix(chunk, "\r")
		chunk = strings.TrimSpace(chunk)
		if chunk == "" {
			continue
		}

		w.emitLine(chunk + truncatedLogSuffix)
	}
}

func (w *otelStdLogWriter) emitLine(line string) {
	now := time.Now()
	var record otellog.Record
	record.SetTimestamp(now)
	record.SetObservedTimestamp(now)
	record.SetSeverity(w.severity)
	record.SetBody(otellog.StringValue(line))
	w.logger.Emit(w.ctx, record)
}

func NewStdLogWriter(loggerName string, severity otellog.Severity) io.Writer {
	w := &otelStdLogWriter{
		loggerName: loggerName,
		severity:   severity,
		logger:     otellogglobal.Logger(loggerName),
		ctx:        context.Background(),
	}
	w.buf.Grow(initialStdLogBufferCap)
	return w
}

type stdlogRedirector struct {
	previousWriter io.Writer
	previousFlags  int
	previousPrefix string
}

func NewStdLogRedirector(loggerName string, severity otellog.Severity) *stdlogRedirector {
	writer := NewStdLogWriter(loggerName, severity)
	previousWriter := log.Writer()
	previousFlags := log.Flags()
	previousPrefix := log.Prefix()
	log.SetOutput(writer)
	log.SetFlags(0)
	log.SetPrefix("")

	return &stdlogRedirector{
		previousWriter: previousWriter,
		previousFlags:  previousFlags,
		previousPrefix: previousPrefix,
	}
}

func (r *stdlogRedirector) Shutdown(ctx context.Context) error {
	log.SetOutput(r.previousWriter)
	log.SetFlags(r.previousFlags)
	log.SetPrefix(r.previousPrefix)
	return nil
}
