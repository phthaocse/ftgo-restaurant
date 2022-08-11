package event

type Handler interface {
	ServeEvent(event DomainEvent)
}

type HandlerFn func(event DomainEvent)

func (f HandlerFn) ServeEvent(event DomainEvent) {
	f(event)
}

type DomainEventHandler struct {
	aggregateType string
	event         DomainEvent
	handler       Handler
}

func NewDomainEventHandler(aggregateType string, event DomainEvent, handler Handler) *DomainEventHandler {
	return &DomainEventHandler{
		aggregateType: aggregateType,
		event:         event,
		handler:       handler,
	}
}
