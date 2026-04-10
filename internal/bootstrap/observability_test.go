package bootstrap

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitObservability_All_NoErrors(t *testing.T) {
	// For these unit tests, force exporters to use noop behavior.
	// internal/config.GetEnv treats empty values as unset.
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "")

	sm := NewShutdownManager(context.Background())

	assert := assert.New(t)
	assert.NoError(InitObservability(sm))
	assert.NoError(sm.ShutdownWithTimeout(2 * time.Second))
}
