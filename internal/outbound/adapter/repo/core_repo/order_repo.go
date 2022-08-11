package core_repo

import (
	"ftgo-restaurant/internal/core/model"
	"ftgo-restaurant/internal/outbound/adapter/repo/postgres_repo"
	"ftgo-restaurant/internal/outbound/interface/repository"
)

type restaurantRepo struct {
	restaurantPostgresRepo *postgres_repo.RestaurantPostgresRepo
}

func NewRestaurantRepo(restaurantPostgresRepo *postgres_repo.RestaurantPostgresRepo) repository.RestaurantRepo {
	return &restaurantRepo{restaurantPostgresRepo: restaurantPostgresRepo}
}

func (r *restaurantRepo) Create(restaurant model.Restaurant) {

}

func (r *restaurantRepo) GetById(id int) {

}
