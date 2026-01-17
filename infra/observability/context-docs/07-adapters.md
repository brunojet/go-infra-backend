# Prompt: Implementação dos Adapters Concretos

Implemente os adapters para cada domínio (metrics, traces, logging) conectando as interfaces definidas aos frameworks/bibliotecas reais (ex: OpenTelemetry, zerolog, etc).

Considere:
- Como mapear cada método do contrato para a implementação concreta?
- Como garantir desacoplamento e facilidade de troca?
- Como lidar com configuração e inicialização dos adapters?

Exemplo de perguntas para guiar a implementação:
- Como tratar erros e fallback?
- Como garantir testabilidade dos adapters?

---

Após os adapters, siga para a integração nas camadas externas (middlewares, handlers, etc).