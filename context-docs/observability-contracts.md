# Contratos de Observabilidade (Ports) — métricas, tracing e logging

Este documento descreve os contratos (ports) que a camada de observabilidade deve expor para integrar providers e exporters em um projeto Go seguindo o padrão Ports & Adapters.

Objetivo
- Fornecer interfaces simples e estáveis para que a aplicação (por exemplo `internal/database`) envie métricas, traces e logs sem depender de implementações concretas.
- Permitir que adapters (OTLP, console, in-memory, etc.) implementem esses contratos.

Diretrizes gerais
- Separar: Provider (TracerProvider / MeterProvider / Logger provider) != Exporter (transportador que envia dados ao backend).
- Interfaces devem ser pequenas e focadas em operações usadas pela aplicação (start/shutdown, record, start span, emit log).
- Fornecer helpers de teste (in-memory exporters) que exponham o buffer para asserts em testes.

1. Contratos (exemplos)

1.1 TracingExporter

Descrição: port simples que permite iniciar o envio de spans e oferecer uma conveniência para spans manuais.

Assinatura mínima (Go):

type TracingExporter interface {
  Start(ctx context.Context) error            // inicializa exporter/provider
  Shutdown(ctx context.Context) error         // encerra/flush
  StartSpan(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, func())
}

Notas:
- `Start` deve configurar/instalar (opcional) um `trace.TracerProvider` ou retornar-ready provider; documentar o comportamento.
- `StartSpan` é conveniência; instrumentações (ex.: GORM plugin) normalmente usam `otel.Tracer` global.

1.2 MetricsExporter

Descrição: port para enviar/registrar métricas simples.

Assinatura mínima (Go):

type MetricsExporter interface {
  Start(ctx context.Context) error
  Shutdown(ctx context.Context) error
  Record(ctx context.Context, name string, value float64, labels ...Label) error
}

Notas:
- `Record` é para métricas ad-hoc; implementações podem mapear para contador/histogram conforme `name`/labels.

1.3 LoggingExporter

Descrição: port para emitir logs estruturados que podem incluir trace/span ids para correlação.

Assinatura mínima (Go):

type LoggingExporter interface {
  Start(ctx context.Context) error
  Shutdown(ctx context.Context) error
  Emit(ctx context.Context, level string, msg string, fields map[string]interface{}) error
}

Notas:
- Quando `ctx` contiver contexto de trace, o exporter deve adicionar `trace_id`/`span_id` se configurado.

2. Providers (helpers)

2.1 Propósito
- Construir `trace.TracerProvider` / `metric.MeterProvider` / `Logger` com exporters e configurações (sampler, batcher, retries).

2.2 Assinaturas sugeridas

func NewTracingProvider(cfg TracingProviderConfig, exporter TraceExporter) (trace.TracerProvider, error)
func NewMetricsProvider(cfg MetricsProviderConfig, exporter MetricsExporter) (metric.MeterProvider, error)
func NewLoggingProvider(cfg LoggingProviderConfig, exporter LoggingExporter) (Logger, error)

Notas:
- Providers devem instalar os providers globalmente por padrão (`setGlobal=true`) para facilitar a integração com instrumentações que usam providers globais (ex.: plugin de GORM). Deve existir uma opção para desabilitar esse comportamento (`setGlobal=false`) quando o caller preferir gerenciar explicitamente o provider.

3. Configurações importantes
- Endpoint OTLP (URL), headers, TLS/insecure, timeout
- Batch settings: max queue, batch size, timeout
- Sampler: always/on-sample/percentage
- Redaction: política para `db.statement` (none, truncate, redact-params)
- Limits: max attribute length, max attributes per span

4. Naming e atributos esperados
- O exporter/provedor não deve alterar os nomes dos atributos OTel padrão. Deve aceitar e transportar:
  - `db.system`, `db.name`, `db.statement` (ou versão redacted), `db.operation`, `db.sql.table`
  - `net.peer.name`, `net.peer.ip`, `net.peer.port`
  - `otel.status_code`, `error` ou `exception.*`

5. Testes e helpers
- Fornecer `in-memory` exporters para traces e métricas (baseado em `tracetest.InMemoryExporter` e equivalente para métricas) e funções utilitárias para inspecionar spans/metric points.
- Exemplos de testes:
  - Unit: mock `TracingExporter` e assert que `StartSpan` é chamado com nome/attrs corretos.
  - Integration: com `in-memory` exporter, executar fluxo (create/read DB) e assert spans exportados contém `db.*` attributes.

6. Critérios de aceitação
- Provider+Exporter configurados exportam spans para um OTLP collector local ou in-memory exporter.
- Spans gerados pelas instrumentações (ex.: GORM plugin) aparecem no exporter com atributos DB esperados.
- Logs incluem `trace_id` quando habilitado.

7. Boas práticas / Segurança
- Não enviar SQL completo por padrão. Aplicar redaction/truncation configurável.
- Permitir toggle para desabilitar envio em ambientes restritos.

8. Documentação e exemplos
- Incluir um README com instruções de bootstrap:
  - Inicializar exporter
  - Opcionalmente obter provider e `otel.SetTracerProvider(tp)`
  - Registrar plugin GORM com `NewOtelGormPlugin(tp, mp)`

9. Perguntas de design
-- `Start()` do exporter deve configurar o provider globalmente por padrão para facilitar a integração com instrumentações que esperam providers globais (ex.: GORM). Também deve permitir que o caller opte por receber o `TracerProvider`/`MeterProvider` sem instalá-lo globalmente (`SetGlobal=false`).

---
Arquivo gerado por automação — editar conforme políticas de deploy do projeto.

## Estrutura de diretórios e arquivos (implementação proposta)

Abaixo a estrutura recomendada e os arquivos já criados nesta primeira implementação OTel:

- internal/observability/
  - exporters/
    - contracts/contract.go         # interfaces para Exporters (TracingExporter, MetricsExporter, LoggingExporter)
    - tracing_port.go               # port wrapper para tracing
    - metrics_port.go               # port wrapper para metrics
    - logging_port.go               # port wrapper para logging
    - adapters/
      - otlp_tracing_exporter.go    # adapter OTLP tracing (gRPC)
      - otlp_metrics_exporter.go    # adapter OTLP metrics (gRPC)
      - otlp_logging_exporter.go    # adapter OTLP logging (gRPC/console fallback)
      - inmemory_tracing_exporter.go# adapter in-memory para testes (tracetest)
      - otlp_client.go              # helper de dial/endpoint
  - providers/
    - tracing_provider.go           # helper para construir TracerProvider (OTLP or local)
    - metrics_provider.go           # helper para MeterProvider
    - logging_provider.go           # helper para Logger provider (zerolog)
  - bootstrap.go                    # helper opcional para inicializar providers/exporters com config

Esses arquivos representam a implementação inicial que já existe no repositório e podem ser estendidos.

## Implementação inicial (observações)

- `adapters/otlp_tracing_exporter.go` e `providers/tracing_provider.go`: fornecem um caminho prático para criar um `TracerProvider` conectado a um collector OTLP por gRPC. O `TracingExporter.Start()` já configura um provider local/fallback quando o collector não está disponível.
- `adapters/otlp_metrics_exporter.go` e `providers/metrics_provider.go`: expõem um `MeterProvider` com `PeriodicReader` apontando para OTLP quando possível, com fallback local.
- `adapters/otlp_logging_exporter.go`: escreve em `zerolog` por padrão e tenta usar OTLP logs se a conexão estiver disponível.
- `adapters/inmemory_tracing_exporter.go`: utilitário de testes baseado em `tracetest.InMemoryExporter` para inspecionar spans em testes automatizados.

## Exemplo de bootstrap (feijão com arroz)

Exemplo mínimo para testes (in-memory) — garantir `Start()` antes de criar a DB/GORM plugin:

```go
trExp, mem := adapters.NewInMemoryTracingExporter()
if err := trExp.Start(ctx); err != nil {
    // handle
}
defer trExp.Shutdown(ctx)

// plugin usa o provider global (plugin foi configurado para usar provider global)
p, _ := plugins.NewOtelGormPlugin(nil, nil)
db, _ := database.NewSQLiteDatabase(p)
// ... run DB ops ...
spans := mem.GetSpans()
```

Exemplo para produção (OTLP):

```go
tp, shutdown, _ := providers.NewOTLPTracerProvider(ctx, "localhost:4317", true)
defer shutdown(ctx)

// opcional: criar exporter adapter que dependa do provider
trExp := adapters.NewOTLPTracingExporter("localhost:4317")
trExp.Start(ctx) // se o exporter configurar global, não necessário; documentar comport.

p, _ := plugins.NewOtelGormPlugin(tp, nil)
db, _ := database.NewSQLiteDatabase(p)
defer db.Close()
```

## Observações finais
- Os arquivos e helpers adicionados no repositório já implementam a base descrita acima. Recomenda-se revisar `providers/*` para ajustar políticas de sampler, batch e redaction conforme necessidade.
- Documentar claramente no README se `Start()` do exporter instala o provider globalmente ou se o caller deve chamar `otel.SetTracerProvider(tp)` (a implementação atual permite ambos, por opção).

## Decisões de execução (sugestões)
Para remover ambiguidade e facilitar execução/integração, sugerimos as decisões operacionais abaixo. Registre-as ou ajuste conforme políticas do time.

- **Instalação global do provider:** por segurança, `Start()` do exporter NÃO deve instalar o provider global por padrão. Expor opção `setGlobal=true` para habilitar instalação global quando explicitamente necessária. Documentar o comportamento em cada adapter/provider.
- **Redaction de SQL:** `db.statement` deve ser redacted por padrão usando política `redact-params`. Permitir override via configuração (ex.: `redaction: none|truncate|redact-params`).
- **Ferramenta de planejamento (`manage_todo_list`):** esta é uma ferramenta do agente/automation para planejar tarefas multi-step; não é algo a ser chamado a partir do código do projeto. Use-a apenas para coordenação de trabalho humana/automação.
- **Formatação/KaTeX:** usar KaTeX somente em documentos que contenham fórmulas matemáticas; para documentação geral usar Markdown padrão do repositório.
- **Ambientes restritos:** envio de telemetria deve estar desabilitado por padrão em ambientes sensíveis; ativação via configuração explícita (ex.: `enableExport=false` por padrão, `enableExport=true` em ambientes autorizados).
- **Nomes e atributos OTel:** não altere nomes de atributos OTel padrão. Qualquer mapeamento ou substituição deve estar documentado e opt-in.
