# Prompt: Implementação do Core/Orquestração de Metrics

Com base nos contratos definidos, implemente o core/orquestração do domínio de métricas. Considere:
- Como centralizar o registro e uso de métricas na aplicação?
- Como expor helpers para facilitar o uso dos contratos?
- Como garantir inicialização e shutdown corretos?

Exemplo de perguntas para guiar a implementação:
- Como registrar e recuperar métricas de forma thread-safe?
- Como facilitar a instrumentação de middlewares e handlers?

---

Após o core, siga para os adapters concretos.