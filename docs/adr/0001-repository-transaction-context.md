# ADR 0001: Padrão de Transação por Contexto no Repositório

- **Status**: Aceito
- **Data**: 2026-03-03
- **Contexto**: `internal/ports/repositories`, `demoapp/services`

## Objetivo

Padronizar o uso de transações para que o fluxo de negócio no serviço execute com consistência (leitura/escrita no mesmo contexto transacional), sem acoplar serviços ao Gorm diretamente.

## Decisão

As implementações de repositório devem resolver o banco efetivo por contexto:

- se existir transação no contexto, usar essa transação;
- caso contrário, usar o `*gorm.DB` base do repositório com `WithContext(ctx)`.

Isso é centralizado em `dbFromContext` em `internal/ports/repositories/repository_impl.go`.

## Motivação

Antes deste padrão, era possível:

- abrir transação no serviço;
- atualizar dados dentro da transação;
- e fazer leitura posterior fora da mesma transação sem perceber.

Isso gerava inconsistência sutil de comportamento (principalmente em update + read-back e regras de sincronização).

## Como usar

### No serviço

1. Envolver o caso de uso com `WithTx` do repositório raiz da operação.
2. Chamar métodos de repositório normalmente dentro do callback transacional.
3. Não passar `*gorm.DB` manualmente entre camadas.

### No repositório

1. Todas as operações (`Create`, `GetByID`, `List`, `Update`, `Delete`) devem usar `dbFromContext(ctx)`.
2. Evitar acessar `g.db` diretamente em operações de domínio.

## Exemplo (fluxo)

- Serviço inicia `WithTx`.
- Repositório recebe `ctx` e detecta `tx` via `TxFromContext`.
- Todas as queries no fluxo usam a mesma transação.
- Commit/Rollback é gerenciado no `WithTx`.

## Benefícios

- Consistência transacional ponta a ponta.
- Serviços mais limpos (foco em regra de negócio).
- Menor chance de bug por leitura fora da transação.
- Testes de fluxo mais previsíveis.

## Regras de ouro

- Serviço não deve carregar plumbing de Gorm.
- Repositório não deve ignorar contexto recebido.
- Regras de orquestração ficam em service/repository, não em hook de model.
- Hook de model deve ficar restrito a integridade local da entidade.

## Responsabilidades por camada

- **Model**
	- Garantir integridade local da entidade (campos/estados válidos e hooks locais).
	- Não abrir transação e não orquestrar fluxo entre agregados.
- **Repository**
	- Executar persistência e consulta.
	- Resolver conexão efetiva via `dbFromContext(ctx)` (tx do contexto ou db base).
	- Não codificar regra de negócio de processo.
- **Service**
	- Orquestrar caso de uso e regras de negócio entre entidades/repositórios.
	- Delimitar transação com `WithTx` quando necessário.
	- Não carregar detalhes de Gorm/SQL.

## Checklist para PR

- [ ] Fluxo de escrita está dentro de `WithTx` quando necessário.
- [ ] Operações de repositório usam `dbFromContext(ctx)`.
- [ ] Não há uso novo de `g.db` direto em métodos de CRUD/lista.
- [ ] Testes de comportamento transacional cobrem o caso alterado.
