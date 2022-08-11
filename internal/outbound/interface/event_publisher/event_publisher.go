package producer

import "ftgo-restaurant/internal/core/event"

type EventPublisher interface {
	Publish(aggregateType string, id interface{}, events []event.DomainEvent)
}
