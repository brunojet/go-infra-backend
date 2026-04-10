package eventbus

import (
	"context"
	"fmt"
	"regexp"
	"sync"

	"github.com/brunojet/go-infra-backend/internal/workerpool"
	ebcontracts "github.com/brunojet/go-infra-backend/pkg/eventbus/contracts"
	wpcontracts "github.com/brunojet/go-infra-backend/pkg/workerpool/contracts"
)

var (
	handlerNameRegex = regexp.MustCompile(`^[a-z][a-z0-9_]{0,99}$`)
)

type WorkerMap map[ebcontracts.HandlerName]*eventWorker

type eventWorker struct {
	pool    wpcontracts.WorkerPool
	handler ebcontracts.Handler
}

type EventBus struct {
	mu      sync.RWMutex
	workers WorkerMap
}

func NewEventBus() *EventBus {
	return &EventBus{
		workers: make(WorkerMap),
	}
}

func (b *EventBus) getWorkerUnsafe(eventType ebcontracts.HandlerName) (*eventWorker, bool) {
	w, ok := b.workers[eventType]
	return w, ok
}

func (b *EventBus) getOrErrorWorkerUnsafe(caller string, eventType ebcontracts.HandlerName) (*eventWorker, error) {
	if err := b.isValidHandlerName(caller, eventType); err != nil {
		return nil, err
	}
	w, ok := b.getWorkerUnsafe(eventType)
	if !ok {
		err := fmt.Errorf("handler não registrado para o tipo de evento: %s", eventType)
		return nil, err
	}
	return w, nil
}

func (b *EventBus) getOrErrorWorker(caller string, eventType ebcontracts.HandlerName) (*eventWorker, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.getOrErrorWorkerUnsafe(caller, eventType)
}

func (b *EventBus) getIfExistsWorker(caller string, eventType ebcontracts.HandlerName) error {
	if err := b.isValidHandlerName(caller, eventType); err != nil {
		return err
	}
	_, ok := b.getWorkerUnsafe(eventType)
	if ok {
		err := fmt.Errorf("handler já registrado para o tipo de evento: %s", eventType)
		return err
	}
	return nil
}

func (b *EventBus) isValidHandlerName(caller string, eventType ebcontracts.HandlerName) error {
	if !handlerNameRegex.MatchString(string(eventType)) {
		return fmt.Errorf("caller: %s eventType inválido: deve começar com letra minúscula, conter apenas letras minúsculas, números ou sublinhado, e ter até 100 caracteres", caller)
	}
	return nil
}

func (b *EventBus) isValidWorkerParams(caller string, handler ebcontracts.Handler, numWorkers int, queueBacklog int) error {
	if handler == nil {
		return fmt.Errorf("caller: %s handler não pode ser nil", caller)
	}

	if numWorkers < 1 || queueBacklog < numWorkers {
		return fmt.Errorf("caller: %s numWorkers deve ser > 0 e queueBacklog >= numWorkers", caller)
	}

	return nil
}

func (b *EventBus) Register(eventType ebcontracts.HandlerName, handler ebcontracts.Handler, numWorkers int, queueBacklog int) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := b.isValidWorkerParams("eventbus.register", handler, numWorkers, queueBacklog); err != nil {
		return err
	}
	if err := b.getIfExistsWorker("eventbus.register", eventType); err != nil {
		return err
	}
	pool := workerpool.NewWorkerPool(numWorkers, queueBacklog)
	pool.Start()
	b.workers[eventType] = &eventWorker{pool: pool, handler: handler}
	return nil
}

func (b *EventBus) Unregister(eventType ebcontracts.HandlerName) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	w, err := b.getOrErrorWorkerUnsafe("eventbus.unregister", eventType)
	if err != nil {
		return err
	}
	w.pool.Stop()
	delete(b.workers, eventType)
	return nil
}

func (b *EventBus) PublishWithContext(ctx context.Context, eventType ebcontracts.HandlerName, data any) error {
	w, err := b.getOrErrorWorker("eventbus.publish_with_context", eventType)
	if err != nil {
		return err
	}

	// enfileirar respeitando contexto; evita manter lock durante enqueue
	enqueued := w.pool.EnqueueWithContext(ctx, func(poolCtx context.Context) {
		realCtx := ctx
		if ctx == context.TODO() {
			realCtx = poolCtx
		}
		// manter comportamento anterior: ignorar erro retornado pelo handler
		_ = w.handler(realCtx, data)
	})
	if !enqueued {
		return fmt.Errorf("evento rejeitado (fila cheia ou pool parado) para o tipo de evento: %s", eventType)
	}
	return nil
}

// Publish mantém compatibilidade, usando contexto nil (sem cancelamento externo)
func (b *EventBus) Publish(eventType ebcontracts.HandlerName, data any) error {
	return b.PublishWithContext(context.TODO(), eventType, data)
}

// Stop encerra todos os workers de todos os tipos de evento e limpa o map.
func (b *EventBus) Stop() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, w := range b.workers {
		w.pool.Stop()
	}
	for k := range b.workers {
		delete(b.workers, k)
	}
}

func (b *EventBus) WaitForHandlers(eventTypes ...ebcontracts.HandlerName) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if len(eventTypes) == 0 {
		for _, w := range b.workers {
			w.pool.Wait()
		}
		return nil
	}
	for _, eventType := range eventTypes {
		w, err := b.getOrErrorWorkerUnsafe("eventbus.wait.error", eventType)
		if err != nil {
			return err
		}
		w.pool.Wait()
	}
	return nil
}
