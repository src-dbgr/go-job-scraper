package app

import (
	"job-scraper/internal/config"
	"job-scraper/internal/models"
	"job-scraper/internal/scraper"
	"job-scraper/internal/scraper/jobsch"
	"net/http"
)

func initScrapers(cfg *config.Config) map[string]scraper.Scraper {
	client := &http.Client{}
	baseURL := cfg.Scrapers[scraper.JobsChScraperName]["base_url"]
	fetcher := jobsch.NewJobsChFetcher(client, baseURL)

	jobsChConfig := jobsch.Config{
		BaseURL:    baseURL,
		MaxPages:   5,
		PageSize:   20,
		JobFetcher: fetcher,
		ParseFunc: func(data []byte) (*models.Job, error) {
			return &models.Job{}, nil
		},
	}

	baseScraper := jobsch.NewJobsChScraper(jobsChConfig)

	// Wrap to PaginatedMetricsDecorator, as JobsChScraper implements the PaginatedScraper Interface
	decoratedScraper := scraper.NewPaginatedMetricsDecorator(baseScraper)

	return map[string]scraper.Scraper{
		scraper.JobsChScraperName: decoratedScraper,
	}
}
