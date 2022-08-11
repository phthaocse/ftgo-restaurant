package main

import (
	"ftgo-restaurant/internal/core/service"
	"ftgo-restaurant/internal/inbound/adapter/rest"
	"ftgo-restaurant/internal/outbound/adapter/logger"
	eventProducer "ftgo-restaurant/internal/outbound/adapter/producer/event"
	messageProducer "ftgo-restaurant/internal/outbound/adapter/producer/message"
	coreRepo "ftgo-restaurant/internal/outbound/adapter/repo/core_repo"
	"ftgo-restaurant/internal/outbound/adapter/repo/postgres_repo"
	"ftgo-restaurant/pkg/producer/kafka"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	pgConn, err := postgres_repo.Init(logger.ZapLogger)
	if err != nil {
		return
	}
	restaurantPostgresRepo := postgres_repo.NewRestaurantPostgresRepo(pgConn)
	restaurantRepo := coreRepo.NewRestaurantRepo(restaurantPostgresRepo)

	kafkaProducer := kafka.NewProducer(logger.ZapLogger)
	restaurantMessageProducer := messageProducer.NewRestaurantMessageProducer(kafkaProducer, logger.ZapLogger)
	restaurantEventProducer := eventProducer.NewRestaurantEventPublisher(restaurantMessageProducer)
	restaurantService := service.NewRestaurantService(restaurantRepo, restaurantEventProducer)
	services := service.BusinessService{
		RestaurantService: restaurantService,
	}

	ginServer := rest.NewGinServer(logger.ZapLogger)
	rest.StartHTTPServer(ginServer, services)

}
