package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/brunojet/go-infra-backend/internal/bootstrap/contracts"
)

type shutdownEntry struct {
	name string
	fn   contracts.ShutdownFunc
}

// shutdownManager stores graceful shutdown handlers.
//
// Registration order matters:
// - Register: append at the end
// - Shutdown: executes in reverse order (LIFO) to respect teardown (deregister) semantics.
type shutdownManager struct {
	ctx     context.Context
	mu      sync.Mutex
	entries []shutdownEntry
	ran     bool
	logger  *slog.Logger
}

var _ contracts.ShutdownManager = (*shutdownManager)(nil)

func NewShutdownManager(ctx context.Context) *shutdownManager {
	return &shutdownManager{ctx: ctx, logger: slog.Default()}
}

// SetLogger sets the logger used by the shutdown manager.
// If logger is nil, slog.Default() is used.
//
// Panics if called after shutdown has started.
func (m *shutdownManager) SetLogger(logger *slog.Logger) {
	if logger == nil {
		logger = slog.Default()
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ran {
		panic("bootstrap: cannot set logger after shutdown started")
	}

	m.logger = logger
}

// NewShutdownManagerWithSignals creates a shutdownManager and a root context cancelled
// on OS signals, returning a stop function that also triggers graceful shutdown (LIFO)
// for all registered handlers.
//
// Typical usage:
//
//	rootCtx, sm, stop := bootstrap.NewShutdownManagerWithSignals(nil, 10*time.Second)
//	defer stop()
//	...
//	<-rootCtx.Done()
//	return // deferred stop() runs shutdown
//
// If parent is nil, context.Background() is used.
// If no signals are provided, defaults to os.Interrupt and syscall.SIGTERM.
func NewShutdownManagerWithSignals(shutdownTimeout time.Duration, signalsToWatch ...os.Signal) (sm *shutdownManager, stop func()) {
	ctx := context.Background()
	if len(signalsToWatch) == 0 {
		signalsToWatch = []os.Signal{os.Interrupt, syscall.SIGTERM}
	}

	rootCtx, cancel := signal.NotifyContext(ctx, signalsToWatch...)
	sm = NewShutdownManager(rootCtx)

	var once sync.Once
	stop = func() {
		once.Do(func() {
			cancel()
			_ = sm.ShutdownWithTimeout(shutdownTimeout)
		})
	}

	return sm, stop
}

func (m *shutdownManager) GetContext() context.Context {
	if m.ctx == nil {
		m.ctx = context.Background()
	}
	return m.ctx
}

// Register adds a shutdown target.
// Name is used only for error context/logging; it may be empty.
//
// Panics if s is nil.
func (m *shutdownManager) Register(name string, s contracts.Shutdown) {
	if s == nil {
		panic("bootstrap: shutdown target is nil")
	}

	m.RegisterFunc(name, s.Shutdown)
}

// RegisterFunc adds a shutdown handler function.
// Name is used only for error context/logging; it may be empty.
//
// Panics if fn is nil.
func (m *shutdownManager) RegisterFunc(name string, fn contracts.ShutdownFunc) {
	if fn == nil {
		panic("bootstrap: shutdown handler is nil")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ran {
		panic("bootstrap: cannot register shutdown handler after shutdown started")
	}

	m.entries = append(m.entries, shutdownEntry{name: name, fn: fn})
}

// RegisterShutdown registers an object implementing Shutdown.
func (m *shutdownManager) RegisterShutdown(name string, s contracts.Shutdown) {
	m.Register(name, s)
}

func (m *shutdownManager) snapshotForShutdown() (entries []shutdownEntry, logger *slog.Logger, alreadyRan bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ran {
		return nil, nil, true
	}
	// mark as ran before executing handlers to prevent duplicate shutdown runs
	m.ran = true

	entries = make([]shutdownEntry, len(m.entries))
	copy(entries, m.entries)
	logger = m.logger
	return entries, logger, false
}

// Shutdown executes all registered handlers in reverse registration order.
// It runs best-effort: all handlers are attempted even if earlier ones fail.
//
// The returned error aggregates all handler errors (errors.Join).
func (m *shutdownManager) shutdown(ctx context.Context) error {
	startAll := time.Now()

	entries, logger, alreadyRan := m.snapshotForShutdown()
	if alreadyRan {
		return nil
	}
	if logger == nil {
		logger = slog.Default()
	}

	if dl, ok := ctx.Deadline(); ok {
		logger.Info("shutdown starting", "handlers", len(entries), "deadline", dl, "time_left", time.Until(dl))
	} else {
		logger.Info("shutdown starting", "handlers", len(entries))
	}

	var errs []error
	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]
		name := entry.name
		if name == "" {
			name = fmt.Sprintf("handler-%d", i)
		}

		hStart := time.Now()
		logger.Info("shutdown handler start", "name", name, "order", (len(entries) - i), "remaining", i)
		if err := entry.fn(ctx); err != nil {
			if entry.name != "" {
				err = fmt.Errorf("shutdown %q: %w", entry.name, err)
			}
			logger.Error("shutdown handler failed", "name", name, "duration", time.Since(hStart), "err", err)
			errs = append(errs, err)
			continue
		}
		logger.Info("shutdown handler done", "name", name, "duration", time.Since(hStart))
	}

	joined := errors.Join(errs...)
	if joined != nil {
		logger.Error("shutdown finished with errors", "duration", time.Since(startAll), "err", joined)
		return joined
	}
	logger.Info("shutdown finished", "duration", time.Since(startAll))
	return nil
}

// ShutdownWithTimeout creates a derived context with timeout (when timeout > 0)
// and runs Shutdown using it.
func (m *shutdownManager) ShutdownWithTimeout(timeout time.Duration) error {
	// Important: do NOT derive the shutdown context from m.ctx directly.
	// m.ctx is typically a signal-cancelled root context; when stop() cancels it,
	// every shutdown handler would see ctx.Err()==context.Canceled and abort flush.
	// We intentionally detach cancellation while keeping any context values.
	parent := context.Background()
	if m.ctx != nil {
		parent = context.WithoutCancel(m.ctx)
	}
	if timeout <= 0 {
		return m.shutdown(parent)
	}
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()
	return m.shutdown(ctx)
}

func (m *shutdownManager) Shutdown() error {
	return m.ShutdownWithTimeout(0)
}
