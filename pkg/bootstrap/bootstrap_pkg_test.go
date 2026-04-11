package bootstrap

import (
	"io"
	"log/slog"
	"path/filepath"
	"testing"
	"time"

	dbcontracts "github.com/brunojet/go-infra-backend/pkg/infra/database/contracts"
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

func TestNewDatabaseWithObservability_UsesEnvConfig(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv(dbcontracts.DB_DRIVER_ENV, string(dbcontracts.DbDriverSQLite))
	t.Setenv(dbcontracts.DB_MODE_ENV, string(dbcontracts.DatabaseModeDisk))
	t.Setenv(dbcontracts.DB_NAME_ENV, filepath.Join(tmpDir, "app.db"))

	sm, stop := NewShutdownManagerWithSignals(2 * time.Second)
	defer stop()

	db, err := NewDatabaseWithObservability(sm)
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
		require.FailNow(t, "timeout waiting for http server termination")
	}
}
