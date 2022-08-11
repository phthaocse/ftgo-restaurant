package service

type BusinessServiceFn func()

type BusinessService struct {
	RestaurantService RestaurantServiceI
}
