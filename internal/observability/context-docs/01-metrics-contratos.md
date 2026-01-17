# Prompt: Definição dos Contratos (Interfaces/Ports) de Metrics

Descreva as interfaces necessárias para abstrair o domínio de métricas na infraestrutura. Considere:
- Quais operações o core da aplicação precisa expor para métricas? (ex: registrar contador, gauge, histogram, etc)
- Como será a interface para registro e coleta?
- Como garantir que a interface seja agnóstica de implementação (OpenTelemetry, Prometheus, etc)?

Exemplo de perguntas para guiar a definição:
- Quais métodos são essenciais para o uso de métricas?
- Como será a tipagem dos dados (nomes, labels, valores)?
- Como facilitar mocks e testes?

---

Após definir as interfaces, siga para a implementação do core/orquestração e, por fim, dos adapters concretos.