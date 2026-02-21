package events

import (
	"context"
	"sync"
	"time"

	contracts "github.com/brunojet/go-infra-backend/pkg/events/contracts"
)

// EventManager coordinates an Adapter's lifecycle and dispatches events to
// registered handlers. It encapsulates goroutine management, event loop
// and graceful shutdown so adapters remain minimal.
type adapterEntry struct {
	adapter  contracts.MessageAdapter
	handlers []contracts.Handler
}

type EventManager struct {
	adapters   []adapterEntry
	mu         sync.Mutex
	wg         sync.WaitGroup // workers + handler goroutines
	producerWg sync.WaitGroup // producers (WaitForEvents) goroutines
	ctx        context.Context
	cancel     context.CancelFunc
	started    bool

	// flow control
	concurrency      int
	queueSize        int
	handlersParallel bool
	eventQueues      map[contracts.MessageAdapter]chan contracts.Message
	// enqueue behavior when queues are full
	enqueueTimeout time.Duration
	nackDelay      time.Duration
}

// NewManager creates an event manager over the provided adapter.
// NewManager creates a manager with default flow control: concurrency=1, queueSize=100, handlersParallel=false
func NewManager(adapters ...contracts.MessageAdapter) *EventManager {
	return NewManagerWithOpts(1, 100, false, adapters...)
}

// NewManagerWithOpts creates a manager with configurable concurrency (worker count),
// queueSize (per-adapter buffer) and whether handlers for the same event run in parallel.
func NewManagerWithOpts(concurrency, queueSize int, handlersParallel bool, adapters ...contracts.MessageAdapter) *EventManager {
	entries := make([]adapterEntry, 0, len(adapters))
	for _, a := range adapters {
		entries = append(entries, adapterEntry{adapter: a, handlers: nil})
	}
	if concurrency <= 0 {
		concurrency = 1
	}
	if queueSize <= 0 {
		queueSize = 100
	}
	return &EventManager{
		adapters:         entries,
		concurrency:      concurrency,
		queueSize:        queueSize,
		handlersParallel: handlersParallel,
		eventQueues:      make(map[contracts.MessageAdapter]chan contracts.Message),
		// sensible defaults: try enqueue for 2s, then ask adapter to requeue with 5s delay
		enqueueTimeout: 2 * time.Second,
		nackDelay:      5 * time.Second,
	}
}

// RegisterHandler registers a handler to be invoked for every incoming event.
// Handlers are executed concurrently and the manager waits for them on Stop.
// RegisterHandlerForAdapter registers a handler for a specific adapter. If the
// adapter was not previously registered it will be added.
func (m *EventManager) RegisterHandlerForAdapter(a contracts.MessageAdapter, h contracts.Handler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i := range m.adapters {
		if m.adapters[i].adapter == a {
			m.adapters[i].handlers = append(m.adapters[i].handlers, h)
			return
		}
	}
	// adapter not found, add it
	m.adapters = append(m.adapters, adapterEntry{adapter: a, handlers: []contracts.Handler{h}})
}

// Start begins the adapter and the internal event loop. It's safe to call
// Start multiple times; subsequent calls are no-ops.
func (m *EventManager) Start(ctx context.Context) error {
	m.mu.Lock()
	if m.started {
		m.mu.Unlock()
		return nil
	}
	m.ctx, m.cancel = context.WithCancel(ctx)
	m.started = true
	m.mu.Unlock()
	// start all registered adapters
	m.mu.Lock()
	entries := append([]adapterEntry(nil), m.adapters...)
	m.mu.Unlock()

	for _, e := range entries {
		if err := m.startAdapter(e.adapter); err != nil {
			return err
		}
	}
	return nil
}

// loop removed: manager now expects adapters to implement WaitForMessages
// and processes Message envelopes via per-adapter queues and workers.

// RegisterAdapter registers an adapter with the manager. If the manager is
// already started the adapter will be started immediately and its goroutines
// created. Otherwise it will be started when `Start` is called.
func (m *EventManager) RegisterAdapter(a contracts.MessageAdapter) error {
	m.mu.Lock()
	m.adapters = append(m.adapters, adapterEntry{adapter: a, handlers: nil})
	started := m.started
	ctx := m.ctx
	m.mu.Unlock()

	if started {
		return m.startAdapterWithCtx(a, ctx)
	}
	return nil
}

// startAdapter starts an adapter using the manager's context and spawns its
// poller + loop goroutines.
func (m *EventManager) startAdapter(a contracts.MessageAdapter) error {
	return m.startAdapterWithCtx(a, m.ctx)
}

func (m *EventManager) startAdapterWithCtx(a contracts.MessageAdapter, ctx context.Context) error {
	if err := a.Start(ctx); err != nil {
		return err
	}
	// create per-adapter queue
	m.mu.Lock()
	if m.eventQueues == nil {
		m.eventQueues = make(map[contracts.MessageAdapter]chan contracts.Message)
	}
	queue := make(chan contracts.Message, m.queueSize)
	m.eventQueues[a] = queue
	m.mu.Unlock()

	// start producer goroutine using the adapter's WaitForMessages implementation
	m.producerWg.Add(1)
	go func() {
		defer m.producerWg.Done()
		for {
			select {
			case <-m.ctx.Done():
				return
			default:
			}

			msgs, err := a.WaitForMessages(m.ctx)
			if err != nil {
				if m.ctx.Err() != nil {
					return
				}
				time.Sleep(time.Second)
				continue
			}
			for _, msg := range msgs {
				// try to enqueue respecting manager context; if queue stays full
				// for longer than enqueueTimeout, ask adapter to requeue (Nack)
				select {
				case queue <- msg:
					// enqueued successfully
				case <-time.After(m.enqueueTimeout):
					// unable to enqueue in time: request adapter to delay/requeue
					_ = msg.Nack(m.ctx, m.nackDelay)
				case <-m.ctx.Done():
					return
				}
			}
		}
	}()

	// start workers that will process messages from the queue
	for i := 0; i < m.concurrency; i++ {
		m.wg.Add(1)
		go m.adapterWorker(a, queue)
	}
	return nil
}

// adapterWorker consumes events from the adapter's queue and processes
// them by invoking registered handlers. If handlersParallel is true each
// handler for the event will run concurrently, otherwise handlers run
// sequentially in this worker goroutine.
func (m *EventManager) adapterWorker(a contracts.MessageAdapter, q <-chan contracts.Message) {
	defer m.wg.Done()
	for msg := range q {
		ev := msg.Event()

		m.mu.Lock()
		var handlers []contracts.Handler
		for _, e := range m.adapters {
			if e.adapter == a {
				handlers = append([]contracts.Handler(nil), e.handlers...)
				break
			}
		}
		m.mu.Unlock()

		var handlerErr error
		if m.handlersParallel {
			// run handlers concurrently and collect errors
			errCh := make(chan error, len(handlers))
			for _, h := range handlers {
				m.wg.Add(1)
				go func(fn contracts.Handler, event contracts.Event) {
					defer m.wg.Done()
					errCh <- fn(m.ctx, event)
				}(h, ev)
			}
			// wait for handler results
			for i := 0; i < len(handlers); i++ {
				if err := <-errCh; err != nil && handlerErr == nil {
					handlerErr = err
				}
			}
		} else {
			// run handlers sequentially in this worker
			for _, h := range handlers {
				if err := h(m.ctx, ev); err != nil {
					handlerErr = err
					break
				}
			}
		}

		// ack/nack based on handlers result
		if handlerErr == nil {
			_ = msg.Ack(m.ctx)
		} else {
			// default nack: immediate requeue (adapter may implement delay)
			_ = msg.Nack(m.ctx, 0)
		}
	}
}

// Shutdown requests shutdown and waits for pending handlers to finish or until
// the provided context expires.
func (m *EventManager) Shutdown(ctx context.Context) error {
	m.mu.Lock()
	if !m.started {
		m.mu.Unlock()
		return nil
	}
	m.cancel()
	m.mu.Unlock()

	// wait for producers to stop producing
	m.producerWg.Wait()

	// close all event queues so workers can drain and exit
	m.mu.Lock()
	for _, q := range m.eventQueues {
		close(q)
	}
	m.mu.Unlock()

	done := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
