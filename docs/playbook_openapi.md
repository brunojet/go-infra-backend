
# Playbook OpenAPI – Projeto go-infra-backend

Este playbook documenta as melhores práticas, padrões e correções aplicadas nos arquivos OpenAPI do projeto. Ele serve como guia para a criação, manutenção e expansão de qualquer endpoint ou módulo OpenAPI, independentemente do domínio (ex: filtros, usuários, produtos, etc). Use este guia como referência para garantir conformidade com linters, ferramentas de documentação e clientes como Insomnia.

---

## 0. Estrutura Recomendada de Diretórios e Modularização

- **Cada endpoint/domínio deve ser autocontido em seu próprio diretório** dentro de `openapi/` (ex: `openapi/usuarios/usuarios.yaml`, `openapi/produtos/produtos.yaml`).
- **Schemas, exemplos e arquivos auxiliares** desse domínio devem ficar juntos no mesmo diretório (ex: `openapi/usuarios/schemas.yaml`).
- O arquivo OpenAPI principal (ex: `openapi/openapi.yaml`) **NÃO deve conter os endpoints diretamente**. Ele deve importar (via `$ref`) os paths, components e demais blocos de cada domínio/endpoints de forma genérica.

Exemplo de estrutura:
```
openapi/
  openapi.yaml           # Arquivo principal, importa os domínios
  usuarios/
    usuarios.yaml        # Endpoints de usuários
    schemas.yaml         # Schemas de usuários
  produtos/
    produtos.yaml        # Endpoints de produtos
    schemas.yaml         # Schemas de produtos
  ...
```

Exemplo de importação genérica no openapi.yaml:
```yaml
paths:
  /usuarios:
    $ref: './usuarios/usuarios.yaml#/paths/~1usuarios'
  /produtos:
    $ref: './produtos/produtos.yaml#/paths/~1produtos'
components:
  schemas:
    Usuario:
      $ref: './usuarios/schemas.yaml#/components/schemas/Usuario'
    Produto:
      $ref: './produtos/schemas.yaml#/components/schemas/Produto'
```

**Vantagens:**
- Facilita manutenção, reuso e expansão.
- Permite times diferentes trabalharem em domínios distintos sem conflitos.
- O OpenAPI principal sempre permanece enxuto e genérico.

---

---

## 1. Estrutura Básica do OpenAPI

- Sempre inicie o arquivo com a versão do OpenAPI:
  ```yaml
  openapi: 3.0.0
  ```
- Defina o servidor padrão:
  ```yaml
  servers:
    - url: http://localhost:8080
  ```
- Preencha o bloco `info` completo:
  ```yaml
  info:
    title: Nome da API
    version: 1.0.0
    description: Descrição da API.
    contact:
      name: Suporte
      email: suporte@example.com
  ```

## 2. Definição Global de Tags

- Sempre defina as tags globais usadas nas operações. Cada domínio (ex: filtros, usuários, produtos) deve ter sua tag:
  ```yaml
  tags:
    - name: filters
      description: Operações relacionadas a filtros
    - name: usuarios
      description: Operações relacionadas a usuários
    - name: produtos
      description: Operações relacionadas a produtos
    # ...adicione novas tags conforme novos domínios/endpoints
  ```

## 3. Operações (paths)

Para cada operação (get, post, put, delete):
- Adicione sempre:
  - `summary`: Resumo curto da operação
  - `description`: Descrição detalhada
  - `operationId`: Identificador único da operação (ex: listUsuarios, createProduto, updatePedido)
  - `tags`: Array com pelo menos uma tag definida globalmente (ex: ["usuarios"])

Exemplo genérico:
```yaml
get:
  summary: Listar entidades
  description: Retorna a lista de entidades do domínio.
  operationId: listEntidades
  tags: ["dominio"]
```

## 4. Parâmetros e RequestBody

- Sempre defina os parâmetros de path, query, etc., explicitamente.
- Para requestBody, use `$ref` para schemas, preferencialmente com caminho relativo ao arquivo principal.

Exemplo genérico:
```yaml
parameters:
  - in: path
    name: entidadeId
    required: true
    schema:
      type: integer
requestBody:
  required: true
  content:
    application/json:
      schema:
        $ref: './schemas.yaml#/components/schemas/Entidade'
```

## 5. Respostas

- Sempre inclua o bloco `responses` com os códigos HTTP esperados e uma descrição.

Exemplo:
```yaml
responses:
  '200':
    description: OK
```

## 6. Correções de Linter

- **Campos obrigatórios:**
  - `servers` deve ser um array não vazio.
  - `info.contact` e `info.description` são obrigatórios.
  - Todas as operações devem ter `description`, `operationId` e `tags`.
  - As tags usadas nas operações devem estar definidas globalmente.
- **Referências ($ref):**
  - Use caminhos relativos ao arquivo principal sempre que possível.
  - Se o Insomnia apresentar erro de path, verifique o diretório de abertura do projeto ou ajuste o path conforme necessário.

## 7. Observações sobre Insomnia

- O Insomnia pode resolver `$ref` de forma diferente do linter/VS Code. Se ocorrer erro de path, tente:
  - Ajustar o path relativo.
  - Usar caminho absoluto ou formato `file://` (não recomendado para versionamento).
  - Garantir que o projeto seja aberto a partir da raiz correta.

## 8. Exemplo Completo de Cabeçalho Genérico

```yaml
openapi: 3.0.0
servers:
  - url: http://localhost:8080
info:
  title: Nome da API
  version: 1.0.0
  description: Descrição da API.
  contact:
    name: Suporte
    email: suporte@example.com
tags:
  - name: dominio
    description: Operações relacionadas ao domínio
```

---

**Mantenha este playbook atualizado conforme novas práticas, domínios e padrões forem adotados no projeto.**
