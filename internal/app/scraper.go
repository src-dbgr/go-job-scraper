package app

import (
	"job-scraper/internal/config"
	"job-scraper/internal/models"
	"job-scraper/internal/scraper"
	"job-scraper/internal/scraper/jobsch"
	"net/http"

	"github.com/rs/zerolog/log"
)

func initScrapers(cfg *config.Config) map[string]scraper.Scraper {
	client := &http.Client{}

	scraperCfg := cfg.Scrapers[scraper.JobsChScraperName]
	if scraperCfg == nil {
		log.Fatal().Msg("Jobs.ch scraper configuration not found")
	}

	baseURL := scraperCfg.BaseURL
	fetcher := jobsch.NewJobsChFetcher(client, baseURL)

	jobsChConfig := jobsch.Config{
		BaseURL:    baseURL, // Verwende die bereinigte URL
		MaxPages:   scraperCfg.MaxPages,
		PageSize:   20,
		JobFetcher: fetcher,
		ParseFunc: func(data []byte) (*models.Job, error) {
			return &models.Job{}, nil
		},
	}

	baseScraper := jobsch.NewJobsChScraper(jobsChConfig)
	decoratedScraper := scraper.NewPaginatedMetricsDecorator(baseScraper)

	return map[string]scraper.Scraper{
		scraper.JobsChScraperName: decoratedScraper,
	}
}
