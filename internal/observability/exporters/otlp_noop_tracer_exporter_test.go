package exporters

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestOTLPNoopTracerExporter_Basics(t *testing.T) {
	exp := NewOTLPNoopTracerExporter()
	require.NotNil(t, exp)

	require.NoError(t, exp.ExportSpans(context.Background(), []sdktrace.ReadOnlySpan{}))
	require.NoError(t, exp.Shutdown(context.Background()))
}
