package server

import (
	"context"
	appConfig "password-validator/infrastructure/config"
	"password-validator/infrastructure/http/router"

	"strconv"
	"sync"
	"time"

	"github.com/itau-corp/itau-jw1-dep-golibs-gotel/logger"
)

type config struct {
	webServer router.Server
}

func NewConfig() *config {
	return &config{}
}

func (config *config) WithWebServer() *config {
	log := logger.FromContext(context.Background())
	intPort, err := strconv.ParseInt(appConfig.C.HttpServerPort, 10, 64)
	if err != nil {
		log.Fatal("error parsing port to int", err)
	}
	duration, err := time.ParseDuration(appConfig.C.ServerTimeout + "s")
	if err != nil {
		log.Fatal("error parsing duration to time duration", err)
	}
	server := router.
		NewGinServer().
		WithPort(intPort).
		WithDuration(duration).
		WithControllers()

	log.Info("Router server has been successfully configured.")

	config.webServer = server
	return config
}

func Init() *config {
	return NewConfig().
		WithWebServer()
}

func (config *config) Start(ctx context.Context, wg *sync.WaitGroup) {
	config.webServer.Listen(ctx, wg)
}
