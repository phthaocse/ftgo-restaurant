package repository

import (
	"ftgo-restaurant/internal/core/model"
)

type RestaurantRepo interface {
	Create(restaurant model.Restaurant)
	GetById(id int)
}
