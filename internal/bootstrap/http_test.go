package bootstrap

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewHttpServerWithObservability_StartAndShutdown(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	sm := NewShutdownManager(ctx)
	hs := NewHttpServerWithObservability(sm)
	assert.NotNil(hs)
	assert.NotNil(hs.Router)

	// register a simple handler to verify router is usable
	hs.Router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// allow OS to pick an available port
	hs.srv.Addr = ":0"

	done := make(chan error, 1)
	go func() {
		done <- hs.StartAndWaitTermination()
	}()

	// brief pause to let server start
	time.Sleep(50 * time.Millisecond)

	// simulate signal: cancel root context then run shutdown (same behavior as NewShutdownManagerWithSignals.stop)
	cancel()
	assert.NoError(sm.ShutdownWithTimeout(2 * time.Second))

	select {
	case err := <-done:
		assert.NoError(err)
	case <-time.After(3 * time.Second):
		assert.Fail("timeout waiting for server termination")
	}
}

func TestStartAndWaitTermination_ListenError(t *testing.T) {
	assert := assert.New(t)

	sm := NewShutdownManager(context.Background())
	hs := NewHttpServerWithObservability(sm)
	assert.NotNil(hs)

	// invalid port to force ListenAndServe immediate error
	hs.srv.Addr = ":99999"

	done := make(chan error, 1)
	go func() { done <- hs.StartAndWaitTermination() }()

	select {
	case err := <-done:
		// StartAndWaitTermination always returns nil, but we ensure it completes
		assert.NoError(err)
	case <-time.After(2 * time.Second):
		assert.Fail("timeout waiting for server to report error")
	}
}
