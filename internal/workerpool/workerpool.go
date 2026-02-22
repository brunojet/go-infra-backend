package workerpool

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"sync/atomic"
	"time"

	wpcontracts "github.com/brunojet/go-infra-backend/pkg/workerpool/contracts"
)

type WorkerPool struct {
	tasks   chan wpcontracts.Task
	workers int
	wg      sync.WaitGroup
	stopped int32      // flag atômico para indicar se o pool foi parado
	mu      sync.Mutex // protege operações críticas de enqueue/stop
	ctx     context.Context
	cancel  context.CancelFunc
	OnPanic func(interface{}) // callback opcional para panics em tasks

	// Métricas
	TasksProcessed  uint64 // total de tasks executadas
	TasksRejected   uint64 // total de tasks rejeitadas
	TasksEnqueued   uint64 // total de tasks aceitas
	CurrentInFlight int64  // tasks em andamento
}

// NewWorkerPool cria um novo WorkerPool com N workers e buffer de tarefas
func NewWorkerPool(workers, buffer int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	wp := &WorkerPool{
		tasks:   make(chan wpcontracts.Task, buffer),
		workers: workers,
		ctx:     ctx,
		cancel:  cancel,
	}
	// Default OnPanic to no-op to avoid propagating panics from tasks.
	wp.OnPanic = func(interface{}) {}
	return wp
}

// Start inicia os workers
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			defer func() {
				if r := recover(); r != nil {
					if wp.OnPanic != nil {
						wp.OnPanic(r)
					} else {
						// TODO: tornar comportamento configurável (ex: PropagatePanic bool)
						// Atualmente re-panica para preservação do comportamento fail-fast.
						panic(r)
					}
				}
			}()
			for task := range wp.tasks {
				atomic.AddInt64(&wp.CurrentInFlight, 1)
				func() {
					defer func() {
						atomic.AddInt64(&wp.CurrentInFlight, -1)
						atomic.AddUint64(&wp.TasksProcessed, 1)
					}()
					task(wp.ctx)
				}()
			}
		}()
	}
}

// Enqueue adiciona uma tarefa ao pool
func (wp *WorkerPool) Enqueue(task wpcontracts.Task) bool {
	return wp.EnqueueWithContext(context.Background(), task)
}

// EnqueueWithContext tenta enfileirar respeitando o contexto fornecido.
// Retorna true se a tarefa foi enfileirada com sucesso, false caso o pool
// esteja parado ou o contexto seja cancelado antes da inserção.
func (wp *WorkerPool) EnqueueWithContext(ctx context.Context, task wpcontracts.Task) (ok bool) {
	if atomic.LoadInt32(&wp.stopped) == 1 {
		atomic.AddUint64(&wp.TasksRejected, 1)
		return false
	}

	// recover para capturar panic ao enviar em canal fechado.
	defer func() {
		if r := recover(); r != nil {
			// Se o pool não estiver marcado como stopped, repropagamos o panic.
			if atomic.LoadInt32(&wp.stopped) != 1 {
				panic(r)
			}
			// Caso esteja stopped, consideramos a task rejeitada.
			atomic.AddUint64(&wp.TasksRejected, 1)
			ok = false
		}
	}()

	select {
	case <-ctx.Done():
		atomic.AddUint64(&wp.TasksRejected, 1)
		return false
	case wp.tasks <- task:
		atomic.AddUint64(&wp.TasksEnqueued, 1)
		return true
	}
}

// Stop encerra o pool, cancela o contexto e aguarda todos os workers finalizarem
func (wp *WorkerPool) Stop() {
	wp.mu.Lock()
	if atomic.CompareAndSwapInt32(&wp.stopped, 0, 1) {
		wp.cancel() // cancela contexto para tasks cooperativas
		close(wp.tasks)
	}
	wp.mu.Unlock()
	wp.wg.Wait()
}

// Wait aguarda todos os workers processarem as tasks pendentes (útil para testes)
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

// Exporta um snapshot das métricas atuais do pool
func (wp *WorkerPool) ExportMetrics() wpcontracts.WorkerPoolMetrics {
	return wpcontracts.WorkerPoolMetrics{
		TasksProcessed:  atomic.LoadUint64(&wp.TasksProcessed),
		TasksRejected:   atomic.LoadUint64(&wp.TasksRejected),
		TasksEnqueued:   atomic.LoadUint64(&wp.TasksEnqueued),
		CurrentInFlight: atomic.LoadInt64(&wp.CurrentInFlight),
	}
}

// Exemplo de exportação periódica para log estruturado
func (wp *WorkerPool) StartMetricsExporter(intervalSec int) chan struct{} {
	stop := make(chan struct{})
	go func() {
		ticker := time.NewTicker(time.Duration(intervalSec) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				metrics := wp.ExportMetrics()
				jsonData, _ := json.Marshal(metrics)
				log.Printf("WORKERPOOL_METRICS %s", string(jsonData))
			case <-stop:
				return
			}
		}
	}()
	return stop
}
