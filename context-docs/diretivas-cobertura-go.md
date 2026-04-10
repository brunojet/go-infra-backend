# Diretivas para Cobertura de Código Go

Para gerar o relatório de cobertura de código em projetos Go, utilize sempre o comando:

```
go test <Caminho> -coverprofile="cover.out"
```

- Substitua `<Caminho>` pelo caminho do(s) pacote(s) que deseja testar (exemplo: `./internal/config/ports/...`).
- O arquivo `cover.out` será gerado no diretório atual.
- Para visualizar o relatório detalhado, utilize:

```
go tool cover -func="cover.out"
```

> **Importante:**
> Execute o comando preferencialmente no terminal cmd do Windows para evitar problemas de geração/leitura do arquivo `cover.out`.
