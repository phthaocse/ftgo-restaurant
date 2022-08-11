package event

type Consumer interface {
	Subscribe(subscribeId string, channels map[string]struct{}, handler Handler)
}
