package logging_test

import (
	"context"
	"errors"
	"testing"

	"github.com/brunojet/go-infra-backend/internal/observability/logging"
	"github.com/brunojet/go-infra-backend/internal/observability/logging/contracts"
)

type mockLogger struct {
	lastMsg    string
	lastFields contracts.Fields
	lastErr    error
	lastCtx    context.Context
}

func (m *mockLogger) Debug(ctx context.Context, msg string, fields contracts.Fields) {
	m.lastMsg = msg
	m.lastFields = fields
	m.lastCtx = ctx
}
func (m *mockLogger) Info(ctx context.Context, msg string, fields contracts.Fields) {
	m.lastMsg = msg
	m.lastFields = fields
	m.lastCtx = ctx
}
func (m *mockLogger) Warn(ctx context.Context, msg string, fields contracts.Fields) {
	m.lastMsg = msg
	m.lastFields = fields
	m.lastCtx = ctx
}
func (m *mockLogger) Error(ctx context.Context, msg string, err error, fields contracts.Fields) {
	m.lastMsg = msg
	m.lastFields = fields
	m.lastErr = err
	m.lastCtx = ctx
}
func (m *mockLogger) WithFields(fields contracts.Fields) contracts.Logger { return m }
func (m *mockLogger) WithContext(ctx context.Context) contracts.Logger    { return m }

func TestLoggingCore_DelegationAndSwap(t *testing.T) {
	m1 := &mockLogger{}
	c := logging.NewLogging(m1)

	c.Info(context.Background(), "hello", contracts.Fields{"k": "v"})
	if m1.lastMsg != "hello" {
		t.Fatalf("expected hello got %s", m1.lastMsg)
	}

	// swap logger
	m2 := &mockLogger{}
	c.SetLogger(m2)
	c.Debug(context.Background(), "x", nil)
	if m2.lastMsg != "x" {
		t.Fatalf("expected x got %s", m2.lastMsg)
	}
}

func TestDerivedLogger_WithFieldsAndContext(t *testing.T) {
	m := &mockLogger{}
	c := logging.NewLogging(m)

	ctx := context.WithValue(context.Background(), "trace", "t1")
	d := c.WithFields(contracts.Fields{"a": "1"}).WithContext(ctx)
	d.Info(context.Background(), "msg", contracts.Fields{"b": "2"})

	if m.lastMsg != "msg" {
		t.Fatalf("expected msg got %s", m.lastMsg)
	}
	if m.lastCtx.Value("trace") != "t1" {
		t.Fatalf("expected t1 trace")
	}
	if m.lastFields["a"] != "1" || m.lastFields["b"] != "2" {
		t.Fatalf("fields not merged: %v", m.lastFields)
	}
}

func TestErrorDelegation(t *testing.T) {
	m := &mockLogger{}
	c := logging.NewLogging(m)
	err := errors.New("boom")
	c.Error(context.Background(), "err", err, contracts.Fields{"x": "y"})
	if m.lastErr == nil {
		t.Fatalf("expected error forwarded")
	}
}
