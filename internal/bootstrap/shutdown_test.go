package bootstrap

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testShutdown struct{ called int }

func (t *testShutdown) Shutdown(context.Context) error {
	t.called++
	return nil
}

func TestShutdownManager_SetLogger_NilUsesDefault(t *testing.T) {
	assert := assert.New(t)
	sm := NewShutdownManager(context.Background())
	sm.SetLogger(nil)
	assert.NotNil(sm.logger, "expected logger to be set")
}

func TestShutdownManager_SetLogger_AfterShutdownPanics(t *testing.T) {
	assert := assert.New(t)
	sm := NewShutdownManager(context.Background())
	sm.RegisterFunc("noop", func(context.Context) error { return nil })
	_ = sm.ShutdownWithTimeout(0)

	assert.Panics(func() { sm.SetLogger(slog.Default()) })
}

func TestShutdownManager_Register_NilPanics(t *testing.T) {
	assert := assert.New(t)
	sm := NewShutdownManager(context.Background())
	assert.Panics(func() { sm.Register("x", nil) })
}

func TestShutdownManager_RegisterFunc_NilPanics(t *testing.T) {
	assert := assert.New(t)
	sm := NewShutdownManager(context.Background())
	assert.Panics(func() { sm.RegisterFunc("x", nil) })
}

func TestShutdownManager_RegisterFunc_AfterShutdownPanics(t *testing.T) {
	assert := assert.New(t)
	sm := NewShutdownManager(context.Background())
	sm.RegisterFunc("noop", func(context.Context) error { return nil })
	_ = sm.ShutdownWithTimeout(0)

	assert.Panics(func() { sm.RegisterFunc("late", func(context.Context) error { return nil }) })
}

func TestShutdownManager_Shutdown_LIFOAndJoinErrors(t *testing.T) {
	assert := assert.New(t)
	sm := NewShutdownManager(context.Background())

	var order []string
	errA := errors.New("A")
	errB := errors.New("B")

	sm.RegisterFunc("first", func(context.Context) error {
		order = append(order, "first")
		return errA
	})
	sm.RegisterFunc("", func(context.Context) error {
		order = append(order, "second")
		return nil
	})
	sm.RegisterFunc("third", func(context.Context) error {
		order = append(order, "third")
		return errB
	})

	err := sm.ShutdownWithTimeout(25 * time.Millisecond)
	assert.Error(err, "expected aggregated error")
	assert.ErrorIs(err, errA)
	assert.ErrorIs(err, errB)

	// LIFO execution: last registered runs first
	want := []string{"third", "second", "first"}
	assert.Equal(want, order)
}

func TestShutdownManager_Shutdown_SecondCallNoop(t *testing.T) {
	assert := assert.New(t)
	sm := NewShutdownManager(context.Background())
	called := 0
	sm.RegisterFunc("x", func(context.Context) error { called++; return nil })
	_ = sm.ShutdownWithTimeout(0)
	_ = sm.ShutdownWithTimeout(0)
	assert.Equal(1, called, "expected handler called once")
}

func TestShutdownManager_ShutdownWithTimeout_DetachesCancellation(t *testing.T) {
	assert := assert.New(t)
	root, cancel := context.WithCancel(context.Background())
	sm := NewShutdownManager(root)

	cancel()
	called := 0
	sm.RegisterFunc("check", func(ctx context.Context) error {
		called++
		assert.NoError(ctx.Err(), "shutdown ctx should not be canceled")
		_, hasDeadline := ctx.Deadline()
		assert.True(hasDeadline, "expected shutdown ctx to have a deadline")
		return nil
	})

	assert.NoError(sm.ShutdownWithTimeout(50 * time.Millisecond))
	assert.Equal(1, called, "expected handler called once")
}

func TestShutdownManager_RegisterShutdown_DelegatesToRegister(t *testing.T) {
	assert := assert.New(t)
	sm := NewShutdownManager(context.Background())
	var s testShutdown
	sm.RegisterShutdown("t", &s)
	_ = sm.Shutdown()
	assert.Equal(1, s.called, "expected shutdown called once")
}

func TestShutdownManager_Shutdown_WrapperExecutesHandlers(t *testing.T) {
	assert := assert.New(t)
	sm := NewShutdownManager(context.Background())
	called := 0
	sm.RegisterFunc("x", func(context.Context) error { called++; return nil })
	assert.NoError(sm.Shutdown())
	assert.Equal(1, called, "expected handler called once")
}

func TestShutdownManager_Shutdown_EmptyNameErrorAndNilLogger(t *testing.T) {
	assert := assert.New(t)
	sm := NewShutdownManager(context.Background())
	sm.logger = nil
	errX := errors.New("X")
	sm.RegisterFunc("", func(context.Context) error { return errX })
	err := sm.Shutdown()
	assert.Error(err)
	assert.ErrorIs(err, errX)
}

func TestNewShutdownManagerWithSignals_StopIdempotent(t *testing.T) {
	assert := assert.New(t)
	sm, stop := NewShutdownManagerWithSignals(50 * time.Millisecond)
	called := 0
	sm.RegisterFunc("x", func(context.Context) error { called++; return nil })

	stop()
	stop()

	assert.Equal(1, called, "expected handler called once")
}
