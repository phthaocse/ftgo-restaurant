package message

type Producer interface {
	Produce(destination string, message Message)
}
