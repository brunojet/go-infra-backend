package exporters

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	otellog "go.opentelemetry.io/otel/log"
	otellogglobal "go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

const initialStdLogBufferCap = 16 * 1024

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

// RedirectStdLog redirects the package-level stdlib `log` output (log.Printf, log.Println, ...)
// to the OpenTelemetry logs pipeline.
//
// It returns a restore function that reverts the stdlib log configuration.

type consoleLoggerExporter struct {
	exporter       sdklog.Exporter
	once           sync.Once
	activated      atomic.Bool
	redirectLog    atomic.Bool
	previousWriter io.Writer
	previousFlags  int
	previousPrefix string
}

// stdoutlogNew exists to allow deterministic testing of the error branch in
// NewConsoleLoggerExporter (stdoutlog.New rarely fails in practice).
var stdoutlogNew = stdoutlog.New

func NewConsoleLoggerExporter(redirectLog bool) (sdklog.Exporter, error) {
	exporter, err := stdoutlogNew(stdoutlog.WithWriter(os.Stdout))
	if err != nil {
		return nil, err
	}

	consoleExporter := consoleLoggerExporter{
		exporter: exporter,
	}

	consoleExporter.redirectLog.Store(redirectLog)

	return &consoleExporter, nil
}

func (e *consoleLoggerExporter) Activate(_ context.Context, loggerName string, severity otellog.Severity) {
	if e.redirectLog.Load() {
		e.once.Do(func() {
			e.activated.Store(true)
			e.previousWriter = log.Writer()
			e.previousFlags = log.Flags()
			e.previousPrefix = log.Prefix()
			log.Println("redirecting stdlib log output to OpenTelemetry logger")
			log.SetOutput(NewStdLogWriter(loggerName, severity))
			log.SetFlags(0)
			log.SetPrefix("")
		})
	}
}

func (e *consoleLoggerExporter) Deactivate(_ context.Context) {
	if !e.redirectLog.Load() || !e.activated.Load() {
		return
	}
	log.SetOutput(e.previousWriter)
	log.SetFlags(e.previousFlags)
	log.SetPrefix(e.previousPrefix)
	e.once = sync.Once{}
	log.Println("restored stdlib log configuration")
	e.activated.Store(false)
}

func (e *consoleLoggerExporter) Export(ctx context.Context, logs []sdklog.Record) error {
	return e.exporter.Export(ctx, logs)
}

func (e *consoleLoggerExporter) ForceFlush(ctx context.Context) error {
	return e.exporter.ForceFlush(ctx)
}

func (e *consoleLoggerExporter) Shutdown(ctx context.Context) error {
	return e.exporter.Shutdown(ctx)
}
