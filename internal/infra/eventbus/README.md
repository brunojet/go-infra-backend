# EventBus

O `EventBus` é um barramento de eventos thread-safe, observável e robusto para aplicações Go, com suporte a múltiplos workers por tipo de evento, contexto para cancelamento, shutdown gracioso e integração com métricas/logs customizados.

## Principais Recursos
- **Registro de handlers** com validação de nome e parâmetros.
- **Worker pool** por tipo de evento, com controle de concorrência e backlog.
- **Publicação de eventos** com suporte a contexto externo (cancelamento/timeout).
- **Observabilidade**: integração com sinks customizados para métricas e logs estruturados.
- **Tratamento de panics** em handlers, sem afetar o fluxo global.
- **API robusta**: Register, Unregister, Publish, PublishWithContext, Stop.
- **Validação de nomes** via regex global.
- **Shutdown seguro** de todos os workers.

## Exemplo de Uso
```go
bus := eventbus.NewEventBus()
err := bus.Register("meu_evento", func(ctx context.Context, evt any) {
    // processa evento
}, 4, 10)
if err != nil { log.Fatal(err) }

// Publica evento
_ = bus.Publish("meu_evento", "payload")

// Com contexto (cancelamento/timeout)
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()
_ = bus.PublishWithContext(ctx, "meu_evento", "payload")

// Remover handler
_ = bus.Unregister("meu_evento")

// Shutdown global
bus.Stop()
```
