package bootstrap

import (
	"context"
	"log/slog"
	"os"
	"time"

	internalbootstrap "github.com/brunojet/go-infra-backend/internal/bootstrap"
	"github.com/brunojet/go-infra-backend/pkg/bootstrap/contracts"
)

// NewShutdownManager delegates to the internal implementation.
func NewShutdownManager(ctx context.Context) contracts.ShutdownManager {
	return internalbootstrap.NewShutdownManager(ctx)
}

// NewShutdownManagerWithSignals delegates to the internal implementation.
func NewShutdownManagerWithSignals(shutdownTimeout time.Duration, signalsToWatch ...os.Signal) (sm contracts.ShutdownManager, stop func()) {
	return internalbootstrap.NewShutdownManagerWithSignals(shutdownTimeout, signalsToWatch...)
}

// SetShutdownLogger configures the logger if the underlying implementation supports it.
// This keeps pkg minimal while still allowing optional configuration.
func SetShutdownLogger(sm contracts.ShutdownManager, logger *slog.Logger) {
	if sm == nil {
		return
	}
	if s, ok := sm.(interface{ SetLogger(*slog.Logger) }); ok {
		s.SetLogger(logger)
	}
}
