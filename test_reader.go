package main

import (
	"context"
	"job-scraper/internal/config"
	"job-scraper/internal/logging"
	"job-scraper/internal/metrics"
	"job-scraper/internal/storage"
	"job-scraper/internal/storage/mongodb"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg("Error loading .env file")
	}

	// Log environment variables (be careful not to log sensitive information in production)
	log.Info().
		Str("MONGODB_URI", os.Getenv("MONGODB_URI")).
		Str("MONGODB_DATABASE", os.Getenv("MONGODB_DATABASE")).
		Msg("Environment variables")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := run(ctx, cancel); err != nil {
		log.Fatal().Err(err).Msg("Application failed")
	}
}

func run(ctx context.Context, cancel context.CancelFunc) error {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	// Log the loaded configuration
	log.Info().
		Str("mongodb_uri", cfg.MongoDB.URI).
		Str("mongodb_database", cfg.MongoDB.Database).
		Str("log_level", cfg.Logging.Level).
		Msg("Loaded configuration")

	// Initialize logger
	logging.InitLogger(cfg.Logging.Level)

	// Initialize MongoDB client
	mongoClient, err := mongodb.NewClient(ctx, cfg.MongoDB.URI, cfg.MongoDB.Database)
	if err != nil {
		return err
	}
	defer mongoClient.Close(ctx)

	// Start Prometheus metrics server
	srv := startPrometheusServer()

	// Process jobs
	go processJobs(ctx, mongoClient)

	// Wait for termination signal
	waitForShutdown(cancel, srv)

	return nil
}

func startPrometheusServer() *http.Server {
	srv := &http.Server{Addr: ":2112"}
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		log.Info().Msg("Starting Prometheus metrics server on :2112")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Prometheus metrics server failed")
		}
	}()

	return srv
}

func processJobs(ctx context.Context, storage storage.Storage) {
	jobs, err := storage.GetJobs(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read jobs from MongoDB")
		return
	}

	log.Info().Msgf("Read %d jobs from MongoDB", len(jobs))

	metrics.ScrapedJobs.Add(float64(len(jobs)))

	for _, job := range jobs {
		select {
		case <-ctx.Done():
			return
		default:
			startTime := time.Now()
			// Simulate some processing...
			time.Sleep(100 * time.Millisecond)
			metrics.ProcessedJobs.Inc()
			metrics.ScraperDuration.Observe(time.Since(startTime).Seconds())

			log.Info().Str("title", job.Title).Msg("Processed job")
		}
	}

	log.Info().Msg("Job processing completed successfully")
}

func waitForShutdown(cancel context.CancelFunc, srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down gracefully...")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exiting")
}
