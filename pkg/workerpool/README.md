# pkg/workerpool

Contrato e utilitários para o pool de workers usado no projeto.

Resumo
- Expõe o construtor `New(workers, buffer int)` e tipos de contrato (`WorkerPool`,
  `Task`, `WorkerPoolMetrics`).

Arquitetura

```mermaid
flowchart LR
  Producer --> Queue[Task Queue]
  Queue -->|consume| W[Worker 1]
  Queue -->|consume| X[Worker 2]
  W --> Handler
  X --> Handler
  note right of Queue: Canal buffered para backlog
```

Uso (exemplo)

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/brunojet/go-infra-backend/pkg/workerpool"
)

func main() {
    wp := workerpool.New(2, 10)
    wp.Start()
    defer wp.Stop()

    wp.Enqueue(func(ctx context.Context) {
        fmt.Println("task executada")
    })

    time.Sleep(20 * time.Millisecond)
}
```

Observações
- `EnqueueWithContext` permite controlar enfileiramento e timeout via `context`.
- `ExportMetrics` retorna métricas para observabilidade.
