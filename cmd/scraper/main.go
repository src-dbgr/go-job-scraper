package main

import (
	"context"
	"job-scraper/internal/app"
	"job-scraper/internal/apperrors"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	// Only try to load .env file if we're not in a container
	if os.Getenv("JOBSCRAPER_IN_CONTAINER") != "true" {
		if err := godotenv.Load(); err != nil {
			log.Warn().Err(err).Msg("Error loading .env file")
		}
	}

	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg("Error loading .env file")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application, err := app.New(ctx)
	if err != nil {
		log.Fatal().Err(apperrors.NewBaseError(
			apperrors.ErrCodeInitialization,
			"Failed to initialize application",
			err,
		)).Msg("Application startup failed")
	}

	log.Info().Msg("Application initialized, starting...")
	go application.Run(ctx)

	waitForShutdown(cancel)
	application.Shutdown(ctx)
}

func waitForShutdown(cancel context.CancelFunc) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down gracefully...")
	cancel()
}
