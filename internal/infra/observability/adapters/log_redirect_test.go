package adapters

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	otellog "go.opentelemetry.io/otel/log"
)

func Test_NewStdLogWriter_PartialWrites(t *testing.T) {
	w := NewStdLogWriter("testlogger", otellog.SeverityInfo)
	// Write partial line
	n, err := w.Write([]byte("partial "))
	require.NoError(t, err)
	require.Equal(t, len("partial "), n)

	// Write remainder with newline
	n2, err := w.Write([]byte("line\n"))
	require.NoError(t, err)
	require.Equal(t, len("line\n"), n2)
}

func Test_NewStdLogRedirector_SetAndRestore(t *testing.T) {
	// Save current log output
	prev := log.Writer()
	defer log.SetOutput(prev)

	r := NewStdLogRedirector("testredirect", otellog.SeverityInfo)
	defer r.Shutdown(context.TODO())

	// After redirect, log.Writer should not be nil and allow writes.
	require.NotNil(t, log.Writer())

	// Exercise writing via log package (should be handled by redirector writer)
	log.Println("hello from stdlog")
	// small sleep to allow background emit (if any)
	time.Sleep(10 * time.Millisecond)
}
