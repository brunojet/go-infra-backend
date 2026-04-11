package bootstrap

import (
	"log/slog"
	"os"
	"time"

	internalbootstrap "github.com/brunojet/go-infra-backend/internal/bootstrap"
	bootcontracts "github.com/brunojet/go-infra-backend/pkg/bootstrap/contracts"
	db "github.com/brunojet/go-infra-backend/pkg/infra/database"
	"github.com/gin-gonic/gin"
)

// ---- Shutdown ----

type Shutdown = bootcontracts.Shutdown

type ShutdownFunc = bootcontracts.ShutdownFunc

type ShutdownManager = bootcontracts.ShutdownManager

// NewShutdownManagerWithSignals delegates to the internal implementation.
func NewShutdownManagerWithSignals(shutdownTimeout time.Duration, signalsToWatch ...os.Signal) (sm ShutdownManager, stop func()) {
	return internalbootstrap.NewShutdownManagerWithSignals(shutdownTimeout, signalsToWatch...)
}

// SetShutdownLogger configures the logger if the underlying implementation supports it.
// This keeps pkg minimal while still allowing optional configuration.
func SetShutdownLogger(sm ShutdownManager, logger *slog.Logger) {
	if sm == nil {
		return
	}
	if s, ok := sm.(interface{ SetLogger(*slog.Logger) }); ok {
		s.SetLogger(logger)
	}
}

// ---- HTTP ----

type HttpServer = internalbootstrap.HttpServer

func NewHttpServerWithObservability(sm ShutdownManager, middlewareFuncs ...gin.HandlerFunc) *HttpServer {
	return internalbootstrap.NewHttpServerWithObservability(sm, middlewareFuncs...)
}

// ---- Database ----
// NewDatabaseWithObservability delegates to the internal bootstrap helper.
func NewDatabaseWithObservability(sm ShutdownManager) (db.DatabaseAdapter, error) {
	return internalbootstrap.NewDatabaseWithObservability(sm)
}
