package exporters

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

func TestOTLPNoopLoggerExporter_Basics(t *testing.T) {
	exp := NewOTLPNoopLoggerExporter()
	require.NotNil(t, exp)

	// Export with empty slice should not error
	require.NoError(t, exp.Export(context.Background(), []sdklog.Record{}))
	require.NoError(t, exp.ForceFlush(context.Background()))
	require.NoError(t, exp.Shutdown(context.Background()))
}
