package main

import (
	"context"
	"job-scraper/internal/app"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg("Error loading .env file")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application, err := app.New(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize application")
	}
	defer application.Shutdown(ctx)

	go application.Run(ctx)

	waitForShutdown(cancel)
}

func waitForShutdown(cancel context.CancelFunc) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down gracefully...")
	cancel()
}
