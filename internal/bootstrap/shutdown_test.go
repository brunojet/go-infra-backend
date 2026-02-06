package bootstrap

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"
)

type testShutdown struct{ called int }

func (t *testShutdown) Shutdown(context.Context) error {
	t.called++
	return nil
}

func TestShutdownManager_GetContext_DefaultsToBackground(t *testing.T) {
	sm := NewShutdownManager(nil)
	ctx := sm.GetContext()
	if ctx == nil {
		t.Fatalf("expected non-nil context")
	}
	if err := ctx.Err(); err != nil {
		t.Fatalf("expected ctx not canceled, got %v", err)
	}
}

func TestShutdownManager_SetLogger_NilUsesDefault(t *testing.T) {
	sm := NewShutdownManager(context.Background())
	sm.SetLogger(nil)
	if sm.logger == nil {
		t.Fatalf("expected logger to be set")
	}
}

func TestShutdownManager_SetLogger_AfterShutdownPanics(t *testing.T) {
	sm := NewShutdownManager(context.Background())
	sm.RegisterFunc("noop", func(context.Context) error { return nil })
	_ = sm.ShutdownWithTimeout(0)

	defer func() {
		rec := recover()
		if rec == nil {
			t.Fatalf("expected panic")
		}
	}()
	sm.SetLogger(slog.Default())
}

func TestShutdownManager_Register_NilPanics(t *testing.T) {
	sm := NewShutdownManager(context.Background())
	defer func() {
		rec := recover()
		if rec == nil {
			t.Fatalf("expected panic")
		}
	}()
	sm.Register("x", nil)
}

func TestShutdownManager_RegisterFunc_NilPanics(t *testing.T) {
	sm := NewShutdownManager(context.Background())
	defer func() {
		rec := recover()
		if rec == nil {
			t.Fatalf("expected panic")
		}
	}()
	sm.RegisterFunc("x", nil)
}

func TestShutdownManager_RegisterFunc_AfterShutdownPanics(t *testing.T) {
	sm := NewShutdownManager(context.Background())
	sm.RegisterFunc("noop", func(context.Context) error { return nil })
	_ = sm.ShutdownWithTimeout(0)

	defer func() {
		rec := recover()
		if rec == nil {
			t.Fatalf("expected panic")
		}
	}()
	sm.RegisterFunc("late", func(context.Context) error { return nil })
}

func TestShutdownManager_Shutdown_LIFOAndJoinErrors(t *testing.T) {
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
	if err == nil {
		t.Fatalf("expected aggregated error")
	}
	if !errors.Is(err, errA) {
		t.Fatalf("expected errors.Is(err, errA) true")
	}
	if !errors.Is(err, errB) {
		t.Fatalf("expected errors.Is(err, errB) true")
	}

	// LIFO execution: last registered runs first
	want := []string{"third", "second", "first"}
	if len(order) != len(want) {
		t.Fatalf("expected %d handlers executed, got %d", len(want), len(order))
	}
	for i := range want {
		if order[i] != want[i] {
			t.Fatalf("order[%d]=%q, want %q", i, order[i], want[i])
		}
	}
}

func TestShutdownManager_Shutdown_SecondCallNoop(t *testing.T) {
	sm := NewShutdownManager(context.Background())
	called := 0
	sm.RegisterFunc("x", func(context.Context) error { called++; return nil })
	_ = sm.ShutdownWithTimeout(0)
	_ = sm.ShutdownWithTimeout(0)
	if called != 1 {
		t.Fatalf("expected handler called once, got %d", called)
	}
}

func TestShutdownManager_ShutdownWithTimeout_DetachesCancellation(t *testing.T) {
	root, cancel := context.WithCancel(context.Background())
	sm := NewShutdownManager(root)

	cancel()
	called := 0
	sm.RegisterFunc("check", func(ctx context.Context) error {
		called++
		if err := ctx.Err(); err != nil {
			t.Fatalf("shutdown ctx should not be canceled, got %v", err)
		}
		_, hasDeadline := ctx.Deadline()
		if !hasDeadline {
			t.Fatalf("expected shutdown ctx to have a deadline")
		}
		return nil
	})

	if err := sm.ShutdownWithTimeout(50 * time.Millisecond); err != nil {
		t.Fatalf("unexpected shutdown error: %v", err)
	}
	if called != 1 {
		t.Fatalf("expected handler called once, got %d", called)
	}
}

func TestShutdownManager_RegisterShutdown_DelegatesToRegister(t *testing.T) {
	sm := NewShutdownManager(context.Background())
	var s testShutdown
	sm.RegisterShutdown("t", &s)
	_ = sm.Shutdown()
	if s.called != 1 {
		t.Fatalf("expected shutdown called once, got %d", s.called)
	}
}

func TestShutdownManager_Shutdown_WrapperExecutesHandlers(t *testing.T) {
	sm := NewShutdownManager(context.Background())
	called := 0
	sm.RegisterFunc("x", func(context.Context) error { called++; return nil })
	if err := sm.Shutdown(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called != 1 {
		t.Fatalf("expected handler called once, got %d", called)
	}
}

func TestShutdownManager_Shutdown_EmptyNameErrorAndNilLogger(t *testing.T) {
	sm := NewShutdownManager(context.Background())
	sm.logger = nil
	errX := errors.New("X")
	sm.RegisterFunc("", func(context.Context) error { return errX })
	err := sm.Shutdown()
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, errX) {
		t.Fatalf("expected errors.Is(err, errX) true")
	}
}

func TestNewShutdownManagerWithSignals_StopIdempotent(t *testing.T) {
	sm, stop := NewShutdownManagerWithSignals(50 * time.Millisecond)
	called := 0
	sm.RegisterFunc("x", func(context.Context) error { called++; return nil })

	stop()
	stop()

	if called != 1 {
		t.Fatalf("expected handler called once, got %d", called)
	}
}
