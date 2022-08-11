package restauarant

import (
	"encoding/json"
	"ftgo-restaurant/internal/core/model"
	"ftgo-restaurant/pkg/message"
)

type Created struct {
	name    string            `json:"name"`
	address string            `json:"address"`
	menu    []*model.MenuItem `json:"menu"`
}

func NewRestaurantCreatedEvent(name, address string, menu []*model.MenuItem) *Created {
	return &Created{
		name:    name,
		address: address,
		menu:    menu,
	}
}

func (e Created) GetEvent() string {
	return "create_restaurant"
}

func (e Created) GetMenu() []*model.MenuItem {
	return e.menu
}

func (e *Created) SetMenu(menu []*model.MenuItem) {
	e.menu = menu
}

func (e Created) GetAggregateId() string {
	return "restaurant"
}

func (e Created) GetMessage() message.Message {
	body, _ := json.Marshal(e)
	header := map[string]string{
		"event": "restaurant_created",
	}
	headerByte, _ := json.Marshal(header)
	return message.Message{Header: headerByte, Payload: body}
}

func (e Created) GetAggregateType() string {
	return "restaurant"
}

func (e Created) GetEventId() string {
	return "restaurant"
}
