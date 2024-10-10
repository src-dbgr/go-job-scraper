package app

import (
	"context"
	"net/http"
	"time"

	"job-scraper/internal/api"
	"job-scraper/internal/config"
	"job-scraper/internal/processor/openai"
	"job-scraper/internal/scheduler"
	"job-scraper/internal/services"
	"job-scraper/internal/storage"

	"github.com/rs/zerolog/log"
)

type App struct {
	cfg             *config.Config
	storage         storage.Storage
	scheduler       *scheduler.Scheduler
	openaiProcessor *openai.Processor
	api             *api.API
	server          *http.Server
}

func New(ctx context.Context) (*App, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}

	initLogger(cfg)

	storage, err := initStorage(ctx, cfg)
	if err != nil {
		return nil, err
	}

	scrapers := initScrapers(cfg)
	initJobCollector(storage)

	sched, err := initScheduler(ctx, storage, scrapers, cfg)
	if err != nil {
		return nil, err
	}

	openaiProcessor, err := initOpenAIProcessor(cfg)
	if err != nil {
		return nil, err
	}

	jobStatsService := services.NewJobStatisticsService(storage)

	apiHandler := api.NewAPI(scrapers, storage, openaiProcessor, jobStatsService)

	return &App{
		cfg:             cfg,
		storage:         storage,
		scheduler:       sched,
		openaiProcessor: openaiProcessor,
		api:             apiHandler,
		server: &http.Server{
			Addr:    ":8080",
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

	log.Info().Str("address", a.server.Addr).Msg("Starting API server")
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

func initOpenAIProcessor(cfg *config.Config) (*openai.Processor, error) {
	openaiConfig := openai.Config{
		APIURL:      cfg.OpenAI.APIURL,
		APIKey:      cfg.OpenAI.APIKey,
		Model:       cfg.OpenAI.Model,
		Temperature: cfg.OpenAI.Temperature,
		MaxTokens:   cfg.OpenAI.MaxTokens,
		TopP:        cfg.OpenAI.TopP,
		FreqPenalty: cfg.OpenAI.FreqPenalty,
		PresPenalty: cfg.OpenAI.PresPenalty,
	}
	promptRepo := openai.NewFilePromptRepository()
	return openai.NewProcessor(openaiConfig, promptRepo), nil
}
