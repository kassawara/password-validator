package main

import (
	"context"
	"fmt"
	"os/signal"
	"password-validator/infrastructure/config"
	httpServer "password-validator/infrastructure/http/server"
	"sync"
	"syscall"

	_ "password-validator/infrastructure/http/docs" // swagger docs

	gotel "github.com/itau-corp/itau-jw1-dep-golibs-gotel"
	_ "github.com/joho/godotenv/autoload" // AutoLoad do .env
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var wg sync.WaitGroup
	if err := config.Load(); err != nil {
		panic(fmt.Sprintf("Failed to load config from environment variables: %v", err.Error()))
	}

	gotel.Start(ctx, &wg,
		gotel.WithLoggingConfig(config.C.AppName, config.C.AppVersion, config.C.LoggingLevel, config.C.Environment),
	)

	httpServer.Init().Start(ctx, &wg)

	wg.Wait()
}
