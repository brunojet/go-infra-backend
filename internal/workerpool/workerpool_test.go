package workerpool

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorkerPool_Basic(t *testing.T) {
	var count int32
	pool := New(3, 10)
	pool.Start()

	tasks := 20
	for i := 0; i < tasks; i++ {
		ok := pool.Enqueue(func(ctx context.Context) {
			atomic.AddInt32(&count, 1)
		})
		assert.True(t, ok)
	}

	pool.Stop()
	assert.Equal(t, int32(tasks), atomic.LoadInt32(&count))
	metrics := pool.ExportMetrics()
	assert.Equal(t, uint64(tasks), metrics.TasksProcessed)
	assert.Equal(t, uint64(tasks), metrics.TasksEnqueued)
	assert.Equal(t, uint64(0), metrics.TasksRejected)
	assert.Equal(t, int64(0), metrics.CurrentInFlight)
}

func TestWorkerPool_StopEarly(t *testing.T) {
	var count int32
	pool := New(2, 2)
	pool.Start()

	ok := pool.Enqueue(func(ctx context.Context) {
		atomic.AddInt32(&count, 1)
	})
	assert.True(t, ok)

	pool.Stop()
	// Após Stop, Enqueue deve falhar
	ok = pool.Enqueue(func(ctx context.Context) {})
	assert.False(t, ok)
	metrics := pool.ExportMetrics()
	assert.Equal(t, uint64(1), metrics.TasksProcessed)
	assert.Equal(t, uint64(1), metrics.TasksEnqueued)
	assert.Equal(t, uint64(1), metrics.TasksRejected)
	assert.Equal(t, int64(0), metrics.CurrentInFlight)
}

func TestWorkerPool_Parallelism(t *testing.T) {
	var count int32
	pool := New(5, 10)
	pool.Start()

	tasks := 5
	start := time.Now()
	for i := 0; i < tasks; i++ {
		pool.Enqueue(func(ctx context.Context) {
			time.Sleep(50 * time.Millisecond)
			atomic.AddInt32(&count, 1)
		})
	}
	pool.Stop()
	dur := time.Since(start)
	assert.Less(t, int(dur.Milliseconds()), 200) // Deve rodar em paralelo
	assert.Equal(t, int32(tasks), atomic.LoadInt32(&count))
	metrics := pool.ExportMetrics()
	assert.Equal(t, uint64(tasks), metrics.TasksProcessed)
	assert.Equal(t, uint64(tasks), metrics.TasksEnqueued)
	assert.Equal(t, uint64(0), metrics.TasksRejected)
	assert.Equal(t, int64(0), metrics.CurrentInFlight)
}

func TestWorkerPool_CooperativeCancel(t *testing.T) {
	var cancelled int32
	pool := New(2, 2)
	pool.Start()
	done := make(chan struct{})
	pool.Enqueue(func(ctx context.Context) {
		select {
		case <-ctx.Done():
			atomic.AddInt32(&cancelled, 1)
		case <-time.After(200 * time.Millisecond):
		}
		close(done)
	})
	time.Sleep(50 * time.Millisecond)
	pool.Stop()
	<-done
	assert.Equal(t, int32(1), atomic.LoadInt32(&cancelled))
	metrics := pool.ExportMetrics()
	assert.Equal(t, uint64(1), metrics.TasksProcessed)
	assert.Equal(t, uint64(1), metrics.TasksEnqueued)
	assert.Equal(t, uint64(0), metrics.TasksRejected)
	assert.Equal(t, int64(0), metrics.CurrentInFlight)
}

func TestWorkerPool_Enqueue_ChannelClosedSuppressPanic(t *testing.T) {
	pool := New(1, 1)
	pool.Start()
	pool.Stop() // fecha o canal
	// Deve suprimir o panic pois o pool está parado
	ok := pool.Enqueue(func(ctx context.Context) {})
	assert.False(t, ok)
	metrics := pool.ExportMetrics()
	assert.Equal(t, uint64(0), metrics.TasksProcessed)
	assert.Equal(t, uint64(0), metrics.TasksEnqueued)
	assert.Equal(t, uint64(1), metrics.TasksRejected)
	assert.Equal(t, int64(0), metrics.CurrentInFlight)
}

func TestWorkerPool_Enqueue_PanicPropagates(t *testing.T) {
	pool := New(1, 1)
	panicCaught := make(chan interface{}, 1)
	pool.OnPanic = func(r interface{}) {
		panicCaught <- r
	}
	pool.Start()
	panicTask := func(ctx context.Context) {
		panic("erro inesperado")
	}
	pool.Enqueue(panicTask)
	pool.Stop()
	select {
	case r := <-panicCaught:
		assert.Equal(t, "erro inesperado", r)
	case <-time.After(1 * time.Second):
		t.Errorf("panic esperado não foi capturado")
	}
	metrics := pool.ExportMetrics()
	assert.Equal(t, uint64(1), metrics.TasksProcessed)
	assert.Equal(t, uint64(1), metrics.TasksEnqueued)
	assert.Equal(t, uint64(0), metrics.TasksRejected)
	assert.Equal(t, int64(0), metrics.CurrentInFlight)
}
