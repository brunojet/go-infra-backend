# Prompt: Definição dos Contratos (Interfaces/Ports) de Logging

Descreva as interfaces necessárias para abstrair o domínio de logging na infraestrutura. Considere:
- Quais operações o core da aplicação precisa expor para logs? (ex: Info, Warn, Error, Debug, WithFields, etc)
- Como será a interface para logging estruturado?
- Como garantir que a interface seja agnóstica de implementação (zerolog, logrus, etc)?

Exemplo de perguntas para guiar a definição:
- Quais métodos são essenciais para o uso de logging?
- Como será a tipagem dos dados (mensagem, campos, nível)?
- Como facilitar mocks e testes?

---

Após definir as interfaces, siga para a implementação do core/orquestração e, por fim, dos adapters concretos.