package bootstrap

import (
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewShutdownManagerWithSignals_StopIdempotent(t *testing.T) {
	sm, stop := NewShutdownManagerWithSignals(50 * time.Millisecond)
	require.NotNil(t, sm)

	stop()
	// stop must be safe to call multiple times
	stop()

	// Shutdown should be safe after stop (no-op)
	require.NoError(t, sm.Shutdown())
}

func TestSetShutdownLogger_DoesNotPanic(t *testing.T) {
	sm, stop := NewShutdownManagerWithSignals(50 * time.Millisecond)
	defer stop()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	require.NotPanics(t, func() {
		SetShutdownLogger(sm, logger)
	})
}

func TestInitObservability_HappyPath(t *testing.T) {
	// Keep OTLP disabled (uses console/noop exporters depending on env defaults)
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "")
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT_LOGGER", "")
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT_METRIC", "")
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT_TRACE", "")

	sm, stop := NewShutdownManagerWithSignals(2 * time.Second)
	defer stop()

	require.NoError(t, InitObservability(sm))
}

func TestNewSQLiteDatabaseWithObservability_Memory(t *testing.T) {
	sm, stop := NewShutdownManagerWithSignals(2 * time.Second)
	defer stop()

	db, err := NewSQLiteDatabaseWithObservability("memory", sm)
	require.NoError(t, err)
	require.NotNil(t, db)

	gormDB, err := db.GormDB()
	require.NoError(t, err)
	require.NotNil(t, gormDB)
}

func TestNewHttpServerWithObservability_Start_ReturnsErrorOnInvalidAddr(t *testing.T) {
	// invalid port should fail fast and return an error
	t.Setenv("HTTP_ADDR", ":99999")

	sm, stop := NewShutdownManagerWithSignals(500 * time.Millisecond)
	defer stop()

	hs := NewHttpServerWithObservability(sm)
	require.NotNil(t, hs)
	require.NotNil(t, hs.Router)

	err := hs.StartAndWaitTermination()
	require.Error(t, err)
}

func TestNewHttpServerWithObservability_StartAndShutdown_HappyPath(t *testing.T) {
	// use :0 to let OS pick an available port
	t.Setenv("HTTP_ADDR", ":0")

	sm, stop := NewShutdownManagerWithSignals(2 * time.Second)
	defer stop()

	hs := NewHttpServerWithObservability(sm)
	require.NotNil(t, hs)

	done := make(chan error, 1)
	go func() {
		done <- hs.StartAndWaitTermination()
	}()

	// Let the server start.
	time.Sleep(50 * time.Millisecond)

	// Trigger graceful shutdown.
	stop()

	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(3 * time.Second):
		t.Fatal("timeout waiting for http server termination")
	}
}

func TestInitLoggerMetricsTracing_IndividualCalls(t *testing.T) {
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "")
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT_LOGGER", "")
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT_METRIC", "")
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT_TRACE", "")

	sm, stop := NewShutdownManagerWithSignals(2 * time.Second)
	defer stop()

	require.NoError(t, InitLogger(sm))
	require.NoError(t, InitMetrics(sm))
	require.NoError(t, InitTracing(sm))

	// Ensure shutdown hooks can run successfully.
	require.NoError(t, sm.ShutdownWithTimeout(2*time.Second))
}
