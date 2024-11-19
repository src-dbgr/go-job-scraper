package app

import (
	"context"
	"fmt"
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

func (a *App) Run(ctx context.Context) {
	log.Info().Msg("Starting application...")

	go func() {
		log.Info().Int("port", a.cfg.Prometheus.Port).Msg("Starting Prometheus metrics server")
		startPrometheusServer(a.cfg.Prometheus.Port)
	}()

	go a.scheduler.Start(ctx)

	addr := fmt.Sprintf(":%d", a.cfg.API.Port)
	log.Info().Str("address", addr).Msg("Starting API server")
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error().Err(err).Msg("API server failed")
	}
}

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
