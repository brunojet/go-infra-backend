package contracts_test

import (
	"context"
	"testing"

	mm "github.com/brunojet/go-infra-backend/internal/observability/logging/contracts"
)

type mockLogger struct{}

func (l *mockLogger) Debug(ctx context.Context, msg string, fields mm.Fields)            {}
func (l *mockLogger) Info(ctx context.Context, msg string, fields mm.Fields)             {}
func (l *mockLogger) Warn(ctx context.Context, msg string, fields mm.Fields)             {}
func (l *mockLogger) Error(ctx context.Context, msg string, err error, fields mm.Fields) {}
func (l *mockLogger) WithFields(fields mm.Fields) mm.Logger                              { return l }
func (l *mockLogger) WithContext(ctx context.Context) mm.Logger                          { return l }

func TestLoggerInterface(t *testing.T) {
	var _ mm.Logger = (*mockLogger)(nil)
	// simple call to ensure methods exist
	l := &mockLogger{}
	l.Info(context.Background(), "ok", mm.Fields{"k": "v"})
}
