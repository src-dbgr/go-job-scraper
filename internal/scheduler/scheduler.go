package scheduler

import (
	"context"
	"job-scraper/internal/scraper"
	"job-scraper/internal/storage"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

type ScraperConfig struct {
	Type     string
	Schedule cron.Schedule
}

type Scheduler struct {
	storage  storage.Storage
	scrapers map[string]scraper.Scraper
	cron     *cron.Cron
}

func NewScheduler(storage storage.Storage, scrapers map[string]scraper.Scraper, configs []ScraperConfig) *Scheduler {
	s := &Scheduler{
		storage:  storage,
		scrapers: scrapers,
		cron:     cron.New(),
	}

	for _, config := range configs {
		scraper := scrapers[config.Type]
		s.cron.Schedule(config.Schedule, cron.FuncJob(func() {
			s.runScraper(context.Background(), scraper)
		}))
	}

	return s
}

func (s *Scheduler) Start(ctx context.Context) {
	s.cron.Start()
	<-ctx.Done()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}

func (s *Scheduler) runScraper(ctx context.Context, scraper scraper.Scraper) {
	jobs, err := scraper.Scrape(ctx)
	if err != nil {
		log.Error().Err(err).Str("scraper", scraper.Name()).Msg("Scrape failed")
		return
	}

	for _, job := range jobs {
		if err := s.storage.SaveJob(ctx, job); err != nil {
			log.Error().Err(err).Msg("Failed to save job")
		}
	}

	log.Info().Str("scraper", scraper.Name()).Int("jobs", len(jobs)).Msg("Scrape completed")
}
