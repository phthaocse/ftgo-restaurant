package event

import (
	"encoding/json"
	"ftgo-restaurant/internal/outbound/interface/logger"
	"ftgo-restaurant/pkg/message"
)

type DomainEventDispatcher struct {
	DispatcherId    string
	Handlers        map[string]*DomainEventHandler
	MessageConsumer message.Consumer
	Logger          logger.Logger
}

func NewDomainEventDispatcher(messageConsumer message.Consumer, logger logger.Logger) *DomainEventDispatcher {
	return &DomainEventDispatcher{
		MessageConsumer: messageConsumer,
		Handlers:        map[string]*DomainEventHandler{},
		Logger:          logger,
	}
}

func (d *DomainEventDispatcher) Dispatch() message.DispatcherFn {
	return func(message message.Message) {
		var eventName string
		var header map[string][]byte
		if err := json.Unmarshal(message.Header, &header); err != nil {
			d.Logger.Errorf("can't read header: %v", err)
			return
		}
		d.Logger.Infof("Catch event %s", header["key"])
		if err := json.Unmarshal(header["key"], &eventName); err != nil {
			d.Logger.Errorf("can't read header: %v", err)
			return
		}
		if handler, ok := d.Handlers[eventName]; ok {
			handler.handler.ServeEvent(handler.event)
		}
		return
	}
}

func (d *DomainEventDispatcher) Subscribe(subscriptionId string, channels map[string]struct{}, handlers []*DomainEventHandler) {
	for _, handler := range handlers {
		d.Handlers[handler.event.GetEvent()] = handler
	}
	d.MessageConsumer.Subscribe(subscriptionId, channels, d.Dispatch())
}
