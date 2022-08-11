package rest

import (
	"context"
	"ftgo-restaurant/internal/core/service"
	"ftgo-restaurant/internal/outbound/interface/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type ginServer struct {
	engine *gin.Engine
	logger logger.Logger
	service.BusinessService
}

func NewGinServer(logger logger.Logger) *ginServer {
	server := &ginServer{}
	server.engine = gin.Default()
	server.logger = logger
	return server
}

func (gs *ginServer) HandlerFnWrapper(fn service.BusinessServiceFn) gin.HandlerFunc {
	return func(c *gin.Context) {
		fn()
	}
}

func (gs *ginServer) InitOrderRoute() {
	orderGroup := gs.engine.Group("restaurant")
	{
		orderGroup.POST("", gs.createRestaurant)
	}
}

func (gs *ginServer) InitRoute() {
	gs.InitOrderRoute()
}

func (gs *ginServer) InitBusinessService(services service.BusinessService) {
	gs.BusinessService = services
}

func (gs *ginServer) InitMiddleware() {
	gs.engine.Use(authMiddleware)
}

func (gs *ginServer) Run() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := &http.Server{
		Addr:    ":8080",
		Handler: gs.engine,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			gs.logger.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	gs.logger.Info("Shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		gs.logger.Fatal("Server forced to shutdown: ", err)
	}

	gs.logger.Info("Server exiting")
}
