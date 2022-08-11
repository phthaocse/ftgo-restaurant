package message

type Consumer interface {
	Subscribe(subscribeId string, channels map[string]struct{}, handler Dispatcher)
}

type Dispatcher interface {
	Dispatch(message Message)
}

type DispatcherFn func(message Message)

func (h DispatcherFn) Dispatch(message Message) {
	h(message)
}
