package postgres_repo

import (
	"ftgo-restaurant/internal/core/model"
	"github.com/jackc/pgconn"
)

type RestaurantPostgresRepo struct {
	pgConn *pgconn.PgConn
}

func NewRestaurantPostgresRepo(pgConn *pgconn.PgConn) *RestaurantPostgresRepo {
	return &RestaurantPostgresRepo{
		pgConn: pgConn,
	}
}

func (r *RestaurantPostgresRepo) Create(restaurant model.Restaurant) {

}

func (r *RestaurantPostgresRepo) GetById(id int) {

}
