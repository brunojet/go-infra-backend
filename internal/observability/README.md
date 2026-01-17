# Observabilidade - Arquitetura e Estrutura de Diretórios

## Visão Geral

Esta camada segue Clean Architecture/Hexagonal, promovendo desacoplamento, testabilidade e evolução independente. Cada domínio de observabilidade (métricas, traces, logging) é modular, com suas próprias abstrações e implementações.

## Estrutura Recomendada

```
infra/observability/
│
├── README.md                # Este documento
├── exports.go               # Fachada para inicialização e integração
│
├── metrics/
│   ├── contracts/           # Interfaces (ports) para métricas
│   ├── adapters/            # Adapters concretos (ex: otelmetrics)
│   └── metrics.go           # Orquestração/Helpers
│
├── traces/
│   ├── contracts/           # Interfaces (ports) para traces
│   ├── adapters/            # Adapters concretos (ex: oteltraces)
│   └── traces.go            # Orquestração/Helpers
│
├── logging/
│   ├── contracts/           # Interfaces (ports) para logging
│   ├── adapters/            # Adapters concretos (ex: zerolog)
│   └── logging.go           # Orquestração/Helpers
│
└── (middlewares fora daqui)
```

## Descrição dos Componentes

- **metrics/**, **traces/**, **logging/**: Cada domínio tem sua própria pasta, com:
    - **contracts/**: Interfaces/ports para abstração.
    - **adapters/**: Implementações concretas (ex: OpenTelemetry, zerolog).
    - Arquivo principal para helpers/orquestração.
- **exports.go**: Ponto de entrada para inicialização/configuração global.
- **middlewares**: Devem ficar fora de observability, consumindo apenas as interfaces abstratas de métricas, traces e logging.

## Benefícios
- Evolução e manutenção isoladas por domínio.
- Troca fácil de ferramentas (ex: só logging, só traces).
- Testabilidade (mocks nos contracts).
- Middlewares desacoplados das implementações.

---

Adapte conforme as necessidades do projeto. Essa estrutura facilita a escalabilidade e a clareza arquitetural.