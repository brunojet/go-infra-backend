package contracts

// MessageQueueAdapter define o contrato mínimo para adaptadores de fila de
// mensagens usados pelo pacote `internal/message`.
type MessageQueueAdapter interface {
	// Start inicia o adapter e chama o callback para cada mensagem recebida.
	// Retorna uma função `stop` para encerrar o processamento.
	Start(onMessage func(event any)) (stop func())
}
