# Prompt: Definição dos Contratos (Interfaces/Ports) de Traces

Descreva as interfaces necessárias para abstrair o domínio de traces na infraestrutura. Considere:
- Quais operações o core da aplicação precisa expor para rastreamento? (ex: iniciar span, finalizar span, adicionar atributos, etc)
- Como será a interface para propagação de contexto?
- Como garantir que a interface seja agnóstica de implementação (OpenTelemetry, Datadog, etc)?

Exemplo de perguntas para guiar a definição:
- Quais métodos são essenciais para o uso de tracing?
- Como será a tipagem dos dados (nomes, atributos, contexto)?
- Como facilitar mocks e testes?

---

Após definir as interfaces, siga para a implementação do core/orquestração e, por fim, dos adapters concretos.