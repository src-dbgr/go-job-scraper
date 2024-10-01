package scraper

import (
	"fmt"
	"job-scraper/internal/scraper/jobsch"
)

func NewScraper(name string, config map[string]string) (Scraper, error) {
	switch name {
	case "jobsch":
		return jobsch.NewJobschScraper(config["base_url"], config["api_key"]), nil
	// to be extended
	default:
		return nil, fmt.Errorf("unknown scraper: %s", name)
	}
}
