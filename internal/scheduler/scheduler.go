package scheduler

import (
	"context"
	"job-scraper/internal/scraper"
	"job-scraper/internal/services"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

type ScraperConfig struct {
	Type     string
	Schedule cron.Schedule
	Pages    int
}

type Scheduler struct {
	scraperService *services.ScraperService
	scrapers       map[string]scraper.Scraper
	configs        map[string]ScraperConfig
	cron           *cron.Cron
}

func NewScheduler(
	scraperService *services.ScraperService,
	scrapers map[string]scraper.Scraper,
	configs []ScraperConfig,
) *Scheduler {
	s := &Scheduler{
		scraperService: scraperService,
		scrapers:       scrapers,
		configs:        make(map[string]ScraperConfig),
		cron:           cron.New(),
	}

	for _, config := range configs {
		s.configs[config.Type] = config
		scraper := scrapers[config.Type]
		s.cron.Schedule(config.Schedule, cron.FuncJob(func() {
			s.runScraper(context.Background(), scraper, config.Type)
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

func (s *Scheduler) runScraper(ctx context.Context, scraper scraper.Scraper, scraperType string) {
	config := s.configs[scraperType]
	log.Info().
		Str("scraper", scraper.Name()).
		Int("pages", config.Pages).
		Msg("Starting scheduled scraping")

	result, err := s.scraperService.ExecuteScraping(ctx, scraper, config.Pages)
	if err != nil {
		log.Error().
			Err(err).
			Str("scraper", scraper.Name()).
			Msg("Scheduled scraping failed")
		return
	}

	log.Info().
		Str("scraper", scraper.Name()).
		Int("total_jobs", result.TotalJobs).
		Int("processed_jobs", result.ProcessedJobs).
		Str("status", result.Status).
		Msg("Scheduled scraping completed")
}
