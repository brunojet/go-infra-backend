package observability

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
func RedirectStdLog(loggerName string, severity otellog.Severity) (restore func()) {
	previousWriter := log.Writer()
	previousFlags := log.Flags()
	previousPrefix := log.Prefix()

	log.SetOutput(NewStdLogWriter(loggerName, severity))
	log.SetFlags(0)
	log.SetPrefix("")

	return func() {
		log.SetOutput(previousWriter)
		log.SetFlags(previousFlags)
		log.SetPrefix(previousPrefix)
	}
}
