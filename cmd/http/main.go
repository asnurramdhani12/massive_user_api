package main

import (
	"context"
	"fmt"
	v1 "user_api/app/v1"
	"user_api/config"
	"user_api/utils/logger"

	"github.com/gofiber/fiber/v2"
)

func main() {
	context := context.Background()

	if err := config.NewAppContext(context); err != nil {
		panic(err)
	}

	// init config
	cfg, err := config.NewConfig(context)
	if err != nil {
		logger.GetLogger(context).Fatalf("failed to load config: %v", err)
		panic(err)
	}

	// Start server
	app := fiber.New(fiber.Config{})

	// Start Deps
	deps := v1.NewDependencies(context, cfg)

	v1.NewRouter(app, deps)

	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	if err := app.Listen(address); err != nil {
		logger.GetLogger(context).Fatalf("failed to start server: %v", err)
	}
}
