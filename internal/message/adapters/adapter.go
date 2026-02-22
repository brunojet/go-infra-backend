package adapters

type MessageQueueAdapter interface {
	Start(onMessage func(event any)) (stop func())
}
