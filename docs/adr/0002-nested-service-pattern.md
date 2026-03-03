# ADR 0001: Padrão de Nested Service

- **Status**: Aceito
- **Data**: 2026-03-03
- **Contexto**: `demoapp/services`, `pkg/ports/services/contracts`, `internal/ports/services`

## Contexto

Temos rotas e casos de uso aninhados (nested), por exemplo:

- perfil dentro de aplicação;
- versão dentro de aplicação + terminal model configuration.

Sem um padrão comum, cada serviço nested tende a repetir:

- aplicação de `parentScopes` no model;
- combinação de escopos de pai e query;
- mapeamento DTO ↔ model;
- paginação/listagem com comportamento inconsistente entre recursos.

## Decisão

Adotamos um `NestedService` genérico como base para operações nested de leitura/listagem e composição de escopos, com especialização local nos serviços de domínio quando necessário.

### Contrato

Usar `NestedService[D, M]` com foco em:

- `CreateNested(ctx, parentScopes, dto)`
- `ListNested(ctx, parentScopes, params)`

### Implementação

- Base genérica em `internal/ports/services/nested_service_impl.go`.
- Contratos em `pkg/ports/services/contracts/nested_contracts.go`.
- Serviços de domínio (ex.: profile/version nested) podem:
  - reutilizar `ListNested` da implementação genérica;
  - sobrescrever `CreateNested` para regras transacionais específicas (ex.: create + archive duplicates).

## Fronteiras de responsabilidade

- **NestedService (genérico)**: composição de escopos, listagem nested, mapeamento DTO/model.
- **Serviço de domínio nested**: regras de negócio do recurso nested (transição, create one-shot, regras de integridade).
- **Model**: integridade local da entidade (sem orquestração de fluxo nested).

## Consequências

### Positivas

- Menos duplicação entre serviços nested.
- Semântica uniforme para endpoints nested.
- Mais fácil evoluir paginação/filtros nested em um ponto central.
- Regras de domínio continuam explícitas nos serviços especializados.

### Negativas / Trade-offs

- Camada extra de abstração pode dificultar onboarding inicial.
- Mapper de parent scope exige disciplina para evitar acoplamento implícito.
- Nem todo caso nested deve ser 100% genérico (create costuma exigir especialização).

## Regras práticas

- Use o genérico para `ListNested` por padrão.
- Especialize `CreateNested` quando houver regra de negócio além de “persistir”.
- Não mover regra de domínio para o mapper.
- Não mover orquestração de nested para hooks de model.

## Alternativas consideradas

1. **Sem NestedService genérico** (cada service implementa tudo)
   - Rejeitada por duplicação e divergência de comportamento.
2. **NestedService totalmente genérico inclusive create complexo**
   - Rejeitada por esconder regra de negócio e reduzir clareza.

## Critérios de revisão em PR

- [ ] `parentScopes` estão aplicados de forma explícita e validada.
- [ ] `ListNested` mantém semântica consistente (itens + total).
- [ ] Regras de create nested complexas estão no serviço especializado.
- [ ] Não há lógica de orquestração de domínio no model.
