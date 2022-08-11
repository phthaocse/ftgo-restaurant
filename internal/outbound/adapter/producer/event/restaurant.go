package event

import (
	"ftgo-restaurant/pkg/event"
	"ftgo-restaurant/pkg/helpers"
	"ftgo-restaurant/pkg/message"
)

type RestaurantEventPublisher struct {
	MessageProducer message.Producer
}

func NewRestaurantEventPublisher(messageProducer message.Producer) *RestaurantEventPublisher {
	return &RestaurantEventPublisher{
		MessageProducer: messageProducer,
	}
}

func (p *RestaurantEventPublisher) Publish(aggregateType string, aggregateId any, events *helpers.ReadOnlySlice[event.DomainEvent]) {
	for _, event := range events.Get() {
		p.MessageProducer.Produce("restaurant", event.GetMessage())
	}
}
