package app

import (
	"job-scraper/internal/config"
	"job-scraper/internal/models"
	"job-scraper/internal/scraper"
	"job-scraper/internal/scraper/jobsch"
)

func initScrapers(cfg *config.Config) map[string]scraper.Scraper {
	jobsChConfig := jobsch.Config{
		BaseURL:  cfg.Scrapers[scraper.JobsChScraperName]["base_url"],
		MaxPages: 5,
		PageSize: 20,
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
