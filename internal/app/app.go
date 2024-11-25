package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"job-scraper/internal/api"
	"job-scraper/internal/apperrors"
	"job-scraper/internal/config"
	"job-scraper/internal/processor"
	"job-scraper/internal/scheduler"
	"job-scraper/internal/services"
	"job-scraper/internal/storage"

	"github.com/rs/zerolog/log"
)

type App struct {
	cfg            *config.Config
	storage        storage.Storage
	scheduler      *scheduler.Scheduler
	processor      processor.JobProcessor
	scraperService *services.ScraperService
	api            *api.API
	server         *http.Server
}

func New(ctx context.Context) (*App, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, apperrors.NewBaseError(apperrors.ErrCodeConfig, "Failed to load configuration", err)
	}

	initLogger(cfg)

	storage, err := initStorage(ctx, cfg)
	if err != nil {
		return nil, apperrors.NewBaseError(apperrors.ErrCodeStorage, "Failed to initialize storage", err)
	}

	scrapers := initScrapers(cfg)
	initMetrics(storage)

	// Initialisiere den Prozessor basierend auf der Konfiguration
	processor, err := initProcessor(cfg)
	if err != nil {
		return nil, apperrors.NewBaseError(apperrors.ErrCodeProcessing, "Failed to initialize processor", err)
	}

	scraperService := services.NewScraperService(storage, processor)

	sched, err := initScheduler(ctx, scraperService, scrapers, cfg)
	if err != nil {
		return nil, apperrors.NewBaseError(apperrors.ErrCodeInitialization, "Failed to initialize scheduler", err)
	}

	jobStatsService := services.NewJobStatisticsService(storage)

	apiHandler := api.NewAPI(scrapers, storage, processor, scraperService, jobStatsService)

	return &App{
		cfg:            cfg,
		storage:        storage,
		scheduler:      sched,
		processor:      processor,
		scraperService: scraperService,
		api:            apiHandler,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.API.Port),
			Handler: apiHandler,
		},
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	log.Info().Msg("Starting application...")

	// Start Prometheus server
	go func() {
		if err := startPrometheusServer(a.cfg.Prometheus.Port); err != nil {
			log.Error().Err(err).Msg("Failed to start Prometheus server")
		}
	}()

	// Start and verify API server
	if err := a.startAPIServer(ctx); err != nil {
		return fmt.Errorf("application startup failed: %w", err)
	}

	// Start scheduler
	go a.scheduler.Start(ctx)

	// Wait for context cancellation
	<-ctx.Done()
	return nil
}

func (a *App) startAPIServer(ctx context.Context) error {
	addr := fmt.Sprintf(":%d", a.cfg.API.Port)
	log.Info().Msgf("Starting API server on port %d", a.cfg.API.Port)

	// First check if port is available
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to start API server: port %d is not available", a.cfg.API.Port)
		return fmt.Errorf("port %d is not available: %w", a.cfg.API.Port, err)
	}

	// Start server with created listener
	serverErrChan := make(chan error, 1)
	go func() {
		if err := a.server.Serve(listener); err != nil && err != http.ErrServerClosed {
			serverErrChan <- err
			log.Error().Err(err).Msg("API server failed")
		}
	}()

	// Wait for either server error or successful startup
	select {
	case err := <-serverErrChan:
		return fmt.Errorf("API server failed to start: %w", err)
	case <-time.After(100 * time.Millisecond):
		// Verify server is running
		if err := a.verifyAPIServer(ctx); err != nil {
			return fmt.Errorf("API server verification failed: %w", err)
		}
	}

	log.Info().Msgf("Started API server on port %d", a.cfg.API.Port)
	return nil
}

// checks if the API server is responding
func (a *App) verifyAPIServer(ctx context.Context) error {
	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: time.Second,
	}

	// Try to connect to health endpoint
	req, err := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("http://localhost:%d/health", a.cfg.API.Port), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status: %d", resp.StatusCode)
	}

	return nil
}

// Shutdown gracefully
func (a *App) Shutdown(ctx context.Context) {
	if err := a.storage.Close(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to close storage")
	}
	a.scheduler.Stop()

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Failed to shutdown API server gracefully")
	}
}
