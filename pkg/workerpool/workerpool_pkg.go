package workerpool

// Package workerpool provides the public contracts for worker pool
// implementations used across the repository.

import (
	iwp "github.com/brunojet/go-infra-backend/internal/workerpool"
	"github.com/brunojet/go-infra-backend/pkg/workerpool/contracts"
)

type (
	Task              = contracts.Task
	WorkerPool        = contracts.WorkerPool
	WorkerPoolMetrics = contracts.WorkerPoolMetrics
)

// NewWorkerPool creates a new WorkerPool and returns it as the exported contract.
func NewWorkerPool(workers, buffer int) WorkerPool {
	return iwp.NewWorkerPool(workers, buffer)
}
