# pkg/eventbus

Contratos e wrapper de alto-nível para o EventBus usado pelo projeto.

Resumo
- Exporta tipos de contrato (Handler, HandlerName, EventBus) e fornece um
  construtor `NewEventBus()` que retorna o contrato.

Arquitetura

```mermaid
flowchart LR
  A[Producer] -->|publica evento| B[EventBus]
  B -->|registro| C[WorkerPool]
  C -->|executa| D[Handlers]
  note right of B: EventBus mantém map de handlers por tipo
  note right of C: WorkerPool enfileira tasks e processa em workers
```

Uso (exemplo)

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/brunojet/go-infra-backend/pkg/eventbus"
)

func main() {
    bus := eventbus.NewEventBus()
    defer bus.Stop()

    _ = bus.Register(eventbus.HandlerName("teste"), func(ctx context.Context, ev any) error {
        fmt.Println("evento recebido:", ev)
        return nil
    }, 2, 4)

    _ = bus.Publish("teste", map[string]any{"msg": "hello"})
    time.Sleep(50 * time.Millisecond)
}
```

Notas
- `Register` aceita `numWorkers` e `queueBacklog` para controlar paralelismo.
- `PublishWithContext` permite cancelar enfileiramento via `context`.