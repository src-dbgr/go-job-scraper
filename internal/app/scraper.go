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
			// TODO Implement job parsing logic here
			// This should convert the raw JSON to a models.Job struct
			return &models.Job{}, nil
		},
	}

	return map[string]scraper.Scraper{
		scraper.JobsChScraperName: jobsch.NewJobsChScraper(jobsChConfig),
	}
}
