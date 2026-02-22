package contracts

import (
	"context"
)

// Task representa a unidade de trabalho enfileirada no workerpool.
type Task func(ctx context.Context)

// WorkerPoolMetrics espelha métricas expostas pelo pool.
type WorkerPoolMetrics struct {
	TasksProcessed  uint64 `json:"tasks_processed"`
	TasksRejected   uint64 `json:"tasks_rejected"`
	TasksEnqueued   uint64 `json:"tasks_enqueued"`
	CurrentInFlight int64  `json:"current_in_flight"`
}

// WorkerPool é o contrato mínimo que a implementação do pool deve prover.
type WorkerPool interface {
	Start()
	Enqueue(task Task) bool
	EnqueueWithContext(ctx context.Context, task Task) bool
	Stop()
	Wait()
	ExportMetrics() WorkerPoolMetrics
	StartMetricsExporter(intervalSec int) chan struct{}
}
