package event

import "ftgo-restaurant/pkg/helpers"

type Producer interface {
	Publish(aggregateType string, aggregateId any, events *helpers.ReadOnlySlice[DomainEvent])
}
