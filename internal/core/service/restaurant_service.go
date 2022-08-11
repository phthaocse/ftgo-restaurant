package service

import (
	"ftgo-restaurant/internal/core/event/restauarant"
	"ftgo-restaurant/internal/core/model"
	"ftgo-restaurant/internal/outbound/interface/repository"
	"ftgo-restaurant/pkg/event"
	"ftgo-restaurant/pkg/helpers"
)

type RestaurantServiceI interface {
	Create(restaurant model.Restaurant)
	FindById(id int)
}

type RestaurantService struct {
	restaurantRepo           repository.RestaurantRepo
	restaurantEventPublisher event.Producer
}

func NewRestaurantService(restaurantRepo repository.RestaurantRepo, restaurantEventPublisher event.Producer) *RestaurantService {
	return &RestaurantService{
		restaurantRepo:           restaurantRepo,
		restaurantEventPublisher: restaurantEventPublisher,
	}
}

func (rs *RestaurantService) Create(restaurant model.Restaurant) {
	restaurantCreatedEvent := restauarant.NewRestaurantCreatedEvent(restaurant.Name, restaurant.Address, restaurant.Menu)
	events := helpers.NewReadOnlySlice([]event.DomainEvent{restaurantCreatedEvent})
	rs.restaurantEventPublisher.Publish("restaurant", "restaurant", events)
}

func (rs *RestaurantService) FindById(id int) {

}
