package contracts

import "context"

// Handler é o tipo de função usado para processar eventos.
type Handler func(ctx context.Context, event any) error

// HandlerName identifica o tipo de evento manejado.
type HandlerName string

// EventBus é o contrato mínimo que uma implementação de barramento de eventos
// deve expor para registro, publicação e gerenciamento de handlers.
type EventBus interface {
	Register(eventType HandlerName, handler Handler, numWorkers int, queueBacklog int) error
	Unregister(eventType HandlerName) error
	PublishWithContext(ctx context.Context, eventType HandlerName, data any) error
	Publish(eventType HandlerName, data any) error
	Stop()
	WaitForHandlers(eventTypes ...HandlerName) error
}
