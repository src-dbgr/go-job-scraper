package app

import (
	"job-scraper/internal/config"
	"job-scraper/internal/scraper"
	"job-scraper/internal/scraper/jobsch"
)

func initScrapers(cfg *config.Config) map[string]scraper.Scraper {
	return map[string]scraper.Scraper{
		"jobsch": jobsch.NewJobschScraper(cfg.Scrapers["jobsch"]["base_url"], cfg.Scrapers["jobsch"]["api_key"]),
	}
}
