package bootstrap

import (
	"context"
	"testing"
	"time"

	bootcontracts "github.com/brunojet/go-infra-backend/pkg/bootstrap/contracts"
	"github.com/stretchr/testify/assert"
)

func TestInitObservability_All_NoErrors(t *testing.T) {
	oldTracer := newOTLPTracerFromEnv
	oldMetric := newOTLPMetricFromEnv
	oldLogger := newOTLPLoggerFromEnv
	t.Cleanup(func() {
		newOTLPTracerFromEnv = oldTracer
		newOTLPMetricFromEnv = oldMetric
		newOTLPLoggerFromEnv = oldLogger
	})

	newOTLPTracerFromEnv = func(ctx context.Context) (bootcontracts.Shutdown, error) { return nil, nil }
	newOTLPMetricFromEnv = func(ctx context.Context) (bootcontracts.Shutdown, error) { return nil, nil }
	newOTLPLoggerFromEnv = func(ctx context.Context) (bootcontracts.Shutdown, error) { return nil, nil }

	sm := NewShutdownManager(context.Background())

	assert := assert.New(t)
	assert.NoError(InitObservability(sm))
	assert.NoError(sm.ShutdownWithTimeout(2 * time.Second))
}
