package exporters

import (
	"context"
	"errors"
	"log"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	otellog "go.opentelemetry.io/otel/log"
	otellogglobal "go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

type capturingLogExporter struct {
	mu      sync.Mutex
	records []sdklog.Record
}

func (e *capturingLogExporter) Export(ctx context.Context, records []sdklog.Record) error {
	_ = ctx
	e.mu.Lock()
	defer e.mu.Unlock()
	for i := range records {
		e.records = append(e.records, records[i].Clone())
	}
	return nil
}

func (e *capturingLogExporter) Shutdown(ctx context.Context) error {
	_ = ctx
	return nil
}

func (e *capturingLogExporter) ForceFlush(ctx context.Context) error {
	_ = ctx
	return nil
}

func (e *capturingLogExporter) snapshot() []sdklog.Record {
	e.mu.Lock()
	defer e.mu.Unlock()
	out := make([]sdklog.Record, len(e.records))
	copy(out, e.records)
	return out
}

func TestStdLogWriter_EmitsOnNewlineAndBuffersChunks(t *testing.T) {
	prevProvider := otellogglobal.GetLoggerProvider()
	defer otellogglobal.SetLoggerProvider(prevProvider)

	exp := &capturingLogExporter{}
	lp := sdklog.NewLoggerProvider(sdklog.WithProcessor(sdklog.NewSimpleProcessor(exp)))
	otellogglobal.SetLoggerProvider(lp)
	defer func() { _ = lp.Shutdown(context.Background()) }()

	w := NewStdLogWriter("test", otellog.SeverityInfo)

	_, err := w.Write([]byte("hello"))
	assert.NoError(t, err)
	assert.Len(t, exp.snapshot(), 0)

	_, err = w.Write([]byte(" world\n"))
	assert.NoError(t, err)

	recs := exp.snapshot()
	assert.Len(t, recs, 1)
	assert.Equal(t, "hello world", recs[0].Body().AsString())
	assert.Equal(t, otellog.SeverityInfo, recs[0].Severity())
}

func TestConsoleLoggerExporter_ActivateDeactivate_RoutesStdlibLogToOTel(t *testing.T) {
	prevProvider := otellogglobal.GetLoggerProvider()
	defer otellogglobal.SetLoggerProvider(prevProvider)
	prevWriter := log.Writer()
	prevFlags := log.Flags()
	prevPrefix := log.Prefix()
	defer func() {
		log.SetOutput(prevWriter)
		log.SetFlags(prevFlags)
		log.SetPrefix(prevPrefix)
	}()

	exp := &capturingLogExporter{}
	lp := sdklog.NewLoggerProvider(sdklog.WithProcessor(sdklog.NewSimpleProcessor(exp)))
	otellogglobal.SetLoggerProvider(lp)
	defer func() { _ = lp.Shutdown(context.Background()) }()

	otelExp, err := NewConsoleLoggerExporter(true)
	assert.NoError(t, err)
	assert.NotNil(t, otelExp)

	// Set known stdlib log config so we can assert it's restored.
	log.SetFlags(123)
	log.SetPrefix("pfx-")
	savedWriter := log.Writer()
	savedFlags := log.Flags()
	savedPrefix := log.Prefix()

	a, ok := otelExp.(interface {
		Activate(context.Context, string, otellog.Severity)
		Deactivate(context.Context)
	})
	assert.True(t, ok)

	a.Activate(context.Background(), "stdlib-test", otellog.SeverityWarn)
	defer a.Deactivate(context.Background())

	log.Print("warn")

	recs := exp.snapshot()
	assert.Len(t, recs, 1)
	assert.Equal(t, "warn", recs[0].Body().AsString())
	assert.Equal(t, otellog.SeverityWarn, recs[0].Severity())

	a.Deactivate(context.Background())
	assert.Equal(t, savedWriter, log.Writer())
	assert.Equal(t, savedFlags, log.Flags())
	assert.Equal(t, savedPrefix, log.Prefix())
}

func TestConsoleLoggerExporter_Activate_WithRedirectDisabled_DoesNotChangeStdlibLog(t *testing.T) {
	prevWriter := log.Writer()
	prevFlags := log.Flags()
	prevPrefix := log.Prefix()

	otelExp, err := NewConsoleLoggerExporter(false)
	assert.NoError(t, err)
	assert.NotNil(t, otelExp)

	a, ok := otelExp.(interface {
		Activate(context.Context, string, otellog.Severity)
		Deactivate(context.Context)
	})
	assert.True(t, ok)

	a.Activate(context.Background(), "stdlib-test", otellog.SeverityInfo)
	assert.Equal(t, prevWriter, log.Writer())
	assert.Equal(t, prevFlags, log.Flags())
	assert.Equal(t, prevPrefix, log.Prefix())

	a.Deactivate(context.Background())
	assert.Equal(t, prevWriter, log.Writer())
	assert.Equal(t, prevFlags, log.Flags())
	assert.Equal(t, prevPrefix, log.Prefix())
}

func TestConsoleLoggerExporter_DeactivateWithoutActivate_NoOp(t *testing.T) {
	otelExp, err := NewConsoleLoggerExporter(true)
	assert.NoError(t, err)
	assert.NotNil(t, otelExp)

	d, ok := otelExp.(interface{ Deactivate(context.Context) })
	assert.True(t, ok)

	prevWriter := log.Writer()
	prevFlags := log.Flags()
	prevPrefix := log.Prefix()

	d.Deactivate(context.Background())
	assert.Equal(t, prevWriter, log.Writer())
	assert.Equal(t, prevFlags, log.Flags())
	assert.Equal(t, prevPrefix, log.Prefix())
}

func TestNewConsoleLoggerExporter_WhenStdoutlogNewFails_ReturnsError(t *testing.T) {
	prev := stdoutlogNew
	defer func() { stdoutlogNew = prev }()

	stdoutlogNew = func(_ ...stdoutlog.Option) (*stdoutlog.Exporter, error) {
		return nil, errors.New("boom")
	}

	exp, err := NewConsoleLoggerExporter(false)
	assert.Error(t, err)
	assert.Nil(t, exp)
}

type fakeSDKLogExporter struct {
	exportCalls     atomic.Int32
	forceFlushCalls atomic.Int32
	shutdownCalls   atomic.Int32

	exportErr     error
	forceFlushErr error
	shutdownErr   error
}

func (f *fakeSDKLogExporter) Export(ctx context.Context, records []sdklog.Record) error {
	_ = ctx
	_ = records
	f.exportCalls.Add(1)
	return f.exportErr
}

func (f *fakeSDKLogExporter) Shutdown(ctx context.Context) error {
	_ = ctx
	f.shutdownCalls.Add(1)
	return f.shutdownErr
}

func (f *fakeSDKLogExporter) ForceFlush(ctx context.Context) error {
	_ = ctx
	f.forceFlushCalls.Add(1)
	return f.forceFlushErr
}

func TestConsoleLoggerExporter_WrapperMethods_DelegateToInnerExporter(t *testing.T) {
	inner := &fakeSDKLogExporter{}
	e := &consoleLoggerExporter{exporter: inner}

	err := e.Export(context.Background(), nil)
	assert.NoError(t, err)
	assert.Equal(t, int32(1), inner.exportCalls.Load())

	err = e.ForceFlush(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, int32(1), inner.forceFlushCalls.Load())

	err = e.Shutdown(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, int32(1), inner.shutdownCalls.Load())
}
