package exporters

import (
	"context"
	"io"
	"os"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	otellog "go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

// RedirectStdLog redirects the package-level stdlib `log` output (log.Printf, log.Println, ...)
// to the OpenTelemetry logs pipeline.
//
// It returns a restore function that reverts the stdlib log configuration.

type consoleLoggerExporter struct {
	exporter       sdklog.Exporter
	redirector     *stdlogRedirector
	redirectLog    bool
	previousWriter io.Writer
	previousFlags  int
	previousPrefix string
}

var stdoutlogNew = stdoutlog.New

func NewConsoleLoggerExporter(redirectLog bool) (sdklog.Exporter, error) {
	exporter, err := stdoutlogNew(stdoutlog.WithWriter(os.Stdout))
	if err != nil {
		return nil, err
	}

	consoleExporter := consoleLoggerExporter{
		exporter:    exporter,
		redirectLog: redirectLog,
	}

	return &consoleExporter, nil
}

func (e *consoleLoggerExporter) Activate(_ context.Context, loggerName string, severity otellog.Severity) {
	if e.redirectLog && e.redirector == nil {
		e.redirector = NewStdLogRedirector(loggerName, severity)
	}
}

func (e *consoleLoggerExporter) Deactivate(ctx context.Context) {
	if e.redirectLog && e.redirector != nil {
		e.redirector.Shutdown(ctx)
	}
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
