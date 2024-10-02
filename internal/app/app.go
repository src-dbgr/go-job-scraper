package app

import (
	"context"
	"job-scraper/internal/config"
	"job-scraper/internal/processor/openai"
	"job-scraper/internal/scheduler"
	"job-scraper/internal/storage"

	"github.com/rs/zerolog/log"
)

type App struct {
	cfg             *config.Config
	storage         storage.Storage
	scheduler       *scheduler.Scheduler
	openaiProcessor *openai.Processor
}

func New(ctx context.Context) (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	// Log the loaded configuration
	log.Info().
		Str("mongodb_uri", cfg.MongoDB.URI).
		Str("mongodb_database", cfg.MongoDB.Database).
		Str("log_level", cfg.Logging.Level).
		Int("prometheus_port", cfg.Prometheus.Port).
		Msg("Loaded configuration")

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

	return &App{
		cfg:             cfg,
		storage:         storage,
		scheduler:       sched,
		openaiProcessor: openaiProcessor,
	}, nil
}

func (a *App) Run(ctx context.Context) {
	go startPrometheusServer(a.cfg.Prometheus.Port)
	a.scheduler.Start(ctx)
}

func (a *App) Shutdown(ctx context.Context) {
	if err := a.storage.Close(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to close storage")
	}
	a.scheduler.Stop()
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
